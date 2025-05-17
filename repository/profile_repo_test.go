package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupProfileTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}
	return db, mock, cleanup
}

func TestCreateProfile_Success(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "profiles" ("created_at","updated_at","deleted_at","name","avatar_url","user_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Test Profile", "http://avatar.url", uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	profile := &domain.Profile{
		Name:      "Test Profile",
		AvatarURL: "http://avatar.url",
		UserID:    1,
	}

	err := repo.Create(context.Background(), profile)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), profile.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProfile_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "profiles" ("created_at","updated_at","deleted_at","name","avatar_url","user_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Test Profile", "http://avatar.url", uint(1)).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	profile := &domain.Profile{
		Name:      "Test Profile",
		AvatarURL: "http://avatar.url",
		UserID:    1,
	}

	err := repo.Create(context.Background(), profile)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProfile_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	profile := &domain.Profile{
		Name:      "Test Profile",
		AvatarURL: "http://avatar.url",
		UserID:    1,
	}

	err := repo.Create(ctx, profile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestFindByUser_Success(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	userID := uint(1)
	now := time.Now()

	profileRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "avatar_url", "user_id"}).
		AddRow(1, now, now, nil, "Profile 1", "http://avatar1.url", userID).
		AddRow(2, now, now, nil, "Profile 2", "http://avatar2.url", userID)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "profiles" WHERE user_id = $1 AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(profileRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "watch_histories" WHERE "watch_histories"."profile_id" IN ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	profiles, err := repo.FindByUser(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, profiles, 2)
	assert.Equal(t, "Profile 1", profiles[0].Name)
	assert.Equal(t, "Profile 2", profiles[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUser_WithWatchHistory(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	userID := uint(1)
	now := time.Now()

	profileRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "avatar_url", "user_id"}).
		AddRow(1, now, now, nil, "Profile 1", "http://avatar1.url", userID)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "profiles" WHERE user_id = $1 AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(profileRows)

	watchHistoryRows := sqlmock.NewRows([]string{"id", "profile_id", "movie_id", "watched_at", "progress"}).
		AddRow(1, 1, 101, now.Add(-time.Hour), 0.75).
		AddRow(2, 1, 102, now.Add(-30*time.Minute), 0.5)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "watch_histories" WHERE "watch_histories"."profile_id" = $1`)).
		WithArgs(1).
		WillReturnRows(watchHistoryRows)

	profiles, err := repo.FindByUser(context.Background(), userID)

	assert.NoError(t, err)
	require.Len(t, profiles, 1, "Debería haber exactamente 1 perfil")

	profile := profiles[0]
	assert.Equal(t, "Profile 1", profile.Name)
	assert.Equal(t, "http://avatar1.url", profile.AvatarURL)
	assert.Equal(t, userID, profile.UserID)

	require.Len(t, profile.WatchHistory, 2, "Debería haber 2 elementos en el historial")

	assert.Equal(t, uint(101), profile.WatchHistory[0].MovieID)
	assert.Equal(t, float32(0.75), profile.WatchHistory[0].Progress)
	assert.WithinDuration(t, now.Add(-time.Hour), profile.WatchHistory[0].WatchedAt, time.Second)

	assert.Equal(t, uint(102), profile.WatchHistory[1].MovieID)
	assert.Equal(t, float32(0.5), profile.WatchHistory[1].Progress)
	assert.WithinDuration(t, now.Add(-30*time.Minute), profile.WatchHistory[1].WatchedAt, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUser_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	userID := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "profiles" WHERE user_id = $1 AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	profiles, err := repo.FindByUser(context.Background(), userID)

	assert.NoError(t, err)
	assert.Empty(t, profiles)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUser_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	userID := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "profiles" WHERE user_id = $1 AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	profiles, err := repo.FindByUser(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, profiles)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUser_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	userID := uint(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	profiles, err := repo.FindByUser(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, profiles)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestDeleteByID_Success(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)
	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "profiles" SET "deleted_at"=$1 WHERE (id = $2 AND user_id = $3) AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), profileID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), profileID, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByID_NotFound(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)
	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "profiles" SET "deleted_at"=$1 WHERE (id = $2 AND user_id = $3) AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), profileID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), profileID, userID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrProfileNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByID_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)
	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "profiles" SET "deleted_at"=$1 WHERE (id = $2 AND user_id = $3) AND "profiles"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), profileID, userID).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.DeleteByID(context.Background(), profileID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByID_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)
	userID := uint(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.DeleteByID(ctx, profileID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestUpdate_Success(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "profiles" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"name"=$4,"avatar_url"=$5,"user_id"=$6 WHERE "profiles"."deleted_at" IS NULL AND "id" = $7`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			"Updated Profile",
			"http://new.avatar.url",
			uint(1),
			profileID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	profile := &domain.Profile{
		ID:        profileID,
		Name:      "Updated Profile",
		AvatarURL: "http://new.avatar.url",
		UserID:    1,
	}

	err := repo.Update(context.Background(), profile)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "profiles" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"name"=$4,"avatar_url"=$5,"user_id"=$6 WHERE "profiles"."deleted_at" IS NULL AND "id" = $7`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			"Updated Profile",
			"http://new.avatar.url",
			uint(1),
			profileID,
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	profile := &domain.Profile{
		ID:        profileID,
		Name:      "Updated Profile",
		AvatarURL: "http://new.avatar.url",
		UserID:    1,
	}

	err := repo.Update(context.Background(), profile)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	profile := &domain.Profile{
		ID:        profileID,
		Name:      "Updated Profile",
		AvatarURL: "http://new.avatar.url",
		UserID:    1,
	}

	err := repo.Update(ctx, profile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestFindByIDAndUser_Success(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	profileID := uint(1)
	userID := uint(42)
	now := time.Now()

	mock.
		ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "profiles" `+
				`WHERE (id = $1 AND user_id = $2) `+
				`AND "profiles"."deleted_at" IS NULL `+
				`ORDER BY "profiles"."id" LIMIT $3`,
		)).
		WithArgs(profileID, userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "name", "avatar_url", "user_id",
		}).AddRow(
			profileID, now, now, nil,
			"Alice", "http://a.jpg", userID,
		))

	mock.
		ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "watch_histories" ` +
				`WHERE "watch_histories"."profile_id" = $1`,
		)).
		WithArgs(profileID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "profile_id", "movie_id", "watched_at", "progress",
		}).AddRow(uint(10), profileID, uint(99), now.Add(-time.Hour), float32(0.5)))

	got, err := repo.FindByIDAndUser(context.Background(), profileID, userID)
	require.NoError(t, err)
	assert.Equal(t, profileID, got.ID)
	assert.Equal(t, "Alice", got.Name)
	assert.Len(t, got.WatchHistory, 1)
	assert.Equal(t, uint(99), got.WatchHistory[0].MovieID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDAndUser_NotFound(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)

	mock.
		ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "profiles" `+
				`WHERE (id = $1 AND user_id = $2) `+
				`AND "profiles"."deleted_at" IS NULL `+
				`ORDER BY "profiles"."id" LIMIT $3`,
		)).
		WithArgs(uint(1), uint(2), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := repo.FindByIDAndUser(context.Background(), 1, 2)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, repository.ErrProfileNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDAndUser_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)

	mock.
		ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "profiles" `+
				`WHERE (id = $1 AND user_id = $2) `+
				`AND "profiles"."deleted_at" IS NULL `+
				`ORDER BY "profiles"."id" LIMIT $3`,
		)).
		WithArgs(uint(5), uint(6), sqlmock.AnyArg()).
		WillReturnError(errors.New("db down"))

	got, err := repo.FindByIDAndUser(context.Background(), 5, 6)
	assert.Nil(t, got)
	assert.EqualError(t, err, "db down")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDAndUser_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupProfileTestDB(t)
	defer cleanup()

	repo := repository.NewProfileRepo(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := repo.FindByIDAndUser(ctx, 1, 1)
	assert.Nil(t, got)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
