package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}
	return gormDB, mock, cleanup
}

func TestMovieRepo_GetAll_Success(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMovieRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"title", "description", "year", "genre", "image_url", "stream_url", "duration",
	}).AddRow(
		1, time.Now(), time.Now(), nil,
		"Movie 1", "Desc", 2020, "Action", "image.jpg", "stream.com", 120,
	).AddRow(
		2, time.Now(), time.Now(), nil,
		"Movie 2", "Desc", 2021, "Drama", "image2.jpg", "stream2.com", 90,
	)

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."deleted_at" IS NULL(?: LIMIT \$1(?: OFFSET \$2)?)?`).
		WillReturnRows(rows)

	ctx := context.Background()
	movies, err := repo.GetAll(ctx, 10, 0)

	require.NoError(t, err)
	require.Len(t, movies, 2)
	require.Equal(t, "Movie 1", movies[0].Title)
	require.Equal(t, "Movie 2", movies[1].Title)
}

func TestMovieRepo_GetAll_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMovieRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"title", "description", "year", "genre", "image_url", "stream_url", "duration",
	})

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."deleted_at" IS NULL(?: LIMIT \$1(?: OFFSET \$2)?)?`).
		WillReturnRows(rows)

	ctx := context.Background()
	movies, err := repo.GetAll(ctx, 10, 0)

	require.NoError(t, err)
	require.Len(t, movies, 0)
}

func TestMovieRepo_GetAll_DBError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMovieRepo(db)

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."deleted_at" IS NULL(?: LIMIT \$1(?: OFFSET \$2)?)?`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	movies, err := repo.GetAll(ctx, 10, 0)

	require.Error(t, err)
	require.Nil(t, movies)
}

func TestMovieRepo_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	movie := &domain.Movie{
		Title:       "The Matrix",
		Description: "Sci-fi action film",
		Year:        1999,
		Genre:       "Action",
		ImageURL:    "matrix.jpg",
		StreamURL:   "stream.matrix.com",
		Duration:    136,
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "movies"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()

	err = repo.Create(context.Background(), movie)
	require.NoError(t, err)
	require.Equal(t, uint(1), movie.ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestMovieRepo_Create_DBError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMovieRepo(db)

	domainMovie := &domain.Movie{
		Title:       "Inception",
		StreamURL:   "stream.inception.com",
		Description: "A mind-bending thriller",
		Year:        2010,
		Genre:       "Thriller",
		ImageURL:    "inception.jpg",
		Duration:    148,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "movies"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			domainMovie.Title, domainMovie.Description, domainMovie.Year, domainMovie.Genre,
			domainMovie.ImageURL, domainMovie.StreamURL, domainMovie.Duration).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	ctx := context.Background()
	err := repo.Create(ctx, domainMovie)

	require.Error(t, err)
	require.Equal(t, uint(0), domainMovie.ID)
}

func TestMovieRepo_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	expectedID := uint(1)
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"title", "description", "year", "genre", "image_url", "stream_url", "duration",
	}).AddRow(
		expectedID, time.Now(), time.Now(), nil,
		"The Matrix", "Sci-fi", 1999, "Action", "matrix.jpg", "stream.matrix.com", 136,
	)

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 .*`).
		WithArgs(expectedID, 1).
		WillReturnRows(rows)

	movie, err := repo.GetByID(context.Background(), expectedID)
	require.NoError(t, err)
	require.NotNil(t, movie)
	require.Equal(t, expectedID, movie.ID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 .*`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	movie, err := repo.GetByID(context.Background(), 999)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	require.Nil(t, movie)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_GetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 .*`).
		WithArgs(1, 1).
		WillReturnError(errors.New("connection error"))

	movie, err := repo.GetByID(context.Background(), 1)
	require.Error(t, err)
	require.Nil(t, movie)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	movie := &domain.Movie{
		ID:          1,
		Title:       "Updated Title",
		Description: "Updated description",
		Year:        2000,
		Genre:       "Action",
		ImageURL:    "updated.jpg",
		StreamURL:   "updated.com",
		Duration:    120,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET .* WHERE id = \$[\d]+ AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			sqlmock.AnyArg(),
			movie.Title,
			movie.Description,
			movie.Year,
			movie.Genre,
			movie.ImageURL,
			movie.StreamURL,
			movie.Duration,
			movie.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Update(context.Background(), movie)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Update_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	movie := &domain.Movie{
		ID:          2,
		Title:       "Error Movie",
		Description: "Fails in DB",
		Year:        2024,
		Genre:       "Drama",
		ImageURL:    "error.jpg",
		StreamURL:   "error.com",
		Duration:    90,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET .* WHERE id = \$[\d]+ AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			sqlmock.AnyArg(),
			movie.Title,
			movie.Description,
			movie.Year,
			movie.Genre,
			movie.ImageURL,
			movie.StreamURL,
			movie.Duration,
			movie.ID,
		).
		WillReturnError(errors.New("db failure"))
	mock.ExpectRollback()

	err = repo.Update(context.Background(), movie)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db failure")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)

	movie := &domain.Movie{
		ID:          999,
		Title:       "Ghost Movie",
		Description: "Does not exist",
		Year:        1990,
		Genre:       "Horror",
		ImageURL:    "ghost.jpg",
		StreamURL:   "ghost.com",
		Duration:    100,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET .* WHERE id = \$[\d]+ AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			sqlmock.AnyArg(),
			movie.Title,
			movie.Description,
			movie.Year,
			movie.Genre,
			movie.ImageURL,
			movie.StreamURL,
			movie.Duration,
			movie.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.Update(context.Background(), movie)
	require.ErrorIs(t, err, repository.ErrMovieNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "deleted_at"=\$1 WHERE "movies"."id" = \$2 AND "movies"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), movieID). // AnyArg para deleted_at (timestamp)
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete(context.Background(), movieID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(999)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "deleted_at"=\$1 WHERE "movies"."id" = \$2 AND "movies"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), movieID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.Delete(context.Background(), movieID)
	require.ErrorIs(t, err, repository.ErrMovieNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Delete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)
	expectedErr := errors.New("database error")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "deleted_at"=\$1 WHERE "movies"."id" = \$2 AND "movies"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), movieID).
		WillReturnError(expectedErr)
	mock.ExpectRollback()

	err = repo.Delete(context.Background(), movieID)
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_Delete_ContextCanceled(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)

	mock.ExpectBegin()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repo.Delete(ctx, movieID)
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled))
}

func TestMovieRepo_UpdateFields_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)
	updates := map[string]interface{}{
		"title":       "New Title",
		"description": "New Description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "description"=\$1,"title"=\$2,"updated_at"=\$3 WHERE id = \$4 AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			updates["description"],
			updates["title"],
			sqlmock.AnyArg(),
			movieID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.UpdateFields(context.Background(), movieID, updates)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_UpdateFields_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(999)
	updates := map[string]interface{}{
		"title": "Non-existent Movie",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "title"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			updates["title"],
			sqlmock.AnyArg(),
			movieID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.UpdateFields(context.Background(), movieID, updates)
	require.ErrorIs(t, err, repository.ErrMovieNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_UpdateFields_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)
	updates := map[string]interface{}{
		"title": "Error Movie",
	}
	expectedErr := errors.New("database error")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "title"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			updates["title"],
			sqlmock.AnyArg(),
			movieID,
		).
		WillReturnError(expectedErr)
	mock.ExpectRollback()

	err = repo.UpdateFields(context.Background(), movieID, updates)
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_UpdateFields_InvalidFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)
	updates := map[string]interface{}{
		"invalid_field": "value",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "movies" SET "invalid_field"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "movies"."deleted_at" IS NULL`).
		WithArgs(
			updates["invalid_field"],
			sqlmock.AnyArg(),
			movieID,
		).
		WillReturnError(gorm.ErrInvalidField)
	mock.ExpectRollback()

	err = repo.UpdateFields(context.Background(), movieID, updates)
	require.Error(t, err)
	require.ErrorIs(t, err, gorm.ErrInvalidField)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMovieRepo_UpdateFields_ContextCanceled(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewMovieRepo(gormDB)
	movieID := uint(1)
	updates := map[string]interface{}{
		"title": "Canceled Update",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repo.UpdateFields(ctx, movieID, updates)
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled))
}
