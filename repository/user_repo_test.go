package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository"
)

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

	cleanup := func() { sqlDB.Close() }
	return db, mock, cleanup
}

func TestUserRepo_Create_Success(t *testing.T) {
	db, mock, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)

	user := &domain.User{
		Email:    "alice@example.com",
		Password: "securepwd",
		IsAdmin:  false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password","is_admin") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			user.Email,
			user.Password,
			user.IsAdmin,
		).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(7))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, uint(7), user.ID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Create_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)

	user := &domain.User{
		Email:    "bob@example.com",
		Password: "pwd123",
		IsAdmin:  true,
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password","is_admin") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			user.Email,
			user.Password,
			user.IsAdmin,
		).
		WillReturnError(errors.New("db insert failure"))

	mock.ExpectRollback()

	err := repo.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db insert failure")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Create_ContextCanceled(t *testing.T) {
	db, _, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	user := &domain.User{
		Email:    "charlie@example.com",
		Password: "pass",
		IsAdmin:  false,
	}

	err := repo.Create(ctx, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestUserRepo_FindByEmail_Success(t *testing.T) {
	db, mock, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)
	ctx := context.Background()

	email := "bob@example.com"

	userRows := sqlmock.NewRows([]string{"id", "email", "password", "is_admin", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, email, "securepass", false, time.Now(), time.Now(), nil)

	profileRows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 1, "Perfil 1", time.Now(), time.Now(), nil)

	historyRows := sqlmock.NewRows([]string{"id", "profile_id", "movie_id", "watched_at"}).
		AddRow(1, 1, 123, time.Now())

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnRows(userRows)

	mock.ExpectQuery(`SELECT \* FROM "profiles" WHERE "profiles"\."user_id" = \$1 AND "profiles"\."deleted_at" IS NULL`).
		WithArgs(1).
		WillReturnRows(profileRows)

	mock.ExpectQuery(`SELECT \* FROM "watch_histories" WHERE "watch_histories"\."profile_id" = \$1`).
		WithArgs(1).
		WillReturnRows(historyRows)

	user, err := repo.FindByEmail(ctx, email)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Len(t, user.Profiles, 1)
	assert.Equal(t, "Perfil 1", user.Profiles[0].Name)
	assert.Len(t, user.Profiles[0].WatchHistory, 1)
	assert.Equal(t, uint(123), user.Profiles[0].WatchHistory[0].MovieID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)
	ctx := context.Background()

	email := "nonexistent@example.com"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail(ctx, email)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_FindByEmail_DBError(t *testing.T) {
	db, mock, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepo(db)
	ctx := context.Background()

	email := "error@example.com"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnError(fmt.Errorf("error de base de datos"))

	user, err := repo.FindByEmail(ctx, email)

	require.Error(t, err)
	require.Nil(t, user)
	require.NoError(t, mock.ExpectationsWereMet())
}
