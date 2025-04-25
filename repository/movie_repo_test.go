package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository"
)

func setupTestMovieDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Movie{})
	require.NoError(t, err)
	return db
}

func createSampleMovie(t *testing.T, repo *repository.MovieRepo, ctx context.Context) *domain.Movie {
	movie := &domain.Movie{
		Title:     "The Vintage Test",
		StreamURL: "http://stream.url/movie.mp4",
		Year:      1985,
		Genre:     "Drama",
		Duration:  120,
	}
	err := repo.Create(ctx, movie)
	require.NoError(t, err)
	return movie
}

func TestCreateMovie(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	t.Run("Create valid movie", func(t *testing.T) {
		movie := &domain.Movie{
			Title:     "Vintage Classic",
			StreamURL: "http://stream.url/movie.mp4",
			Year:      1975,
			Genre:     "Drama",
			Duration:  95,
		}
		err := repo.Create(ctx, movie)
		require.NoError(t, err)
		assert.NotZero(t, movie.ID)
	})

	t.Run("Create invalid movie (no title)", func(t *testing.T) {
		movie := &domain.Movie{
			StreamURL: "http://stream.url/movie.mp4",
		}
		err := repo.Create(ctx, movie)
		assert.Error(t, err)
	})
}

func TestGetMovieByID(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	movie := createSampleMovie(t, repo.(*repository.MovieRepo), ctx)

	t.Run("Get existing movie", func(t *testing.T) {
		found, err := repo.GetByID(ctx, movie.ID)
		require.NoError(t, err)
		assert.Equal(t, movie.Title, found.Title)
	})

	t.Run("Get non-existent movie", func(t *testing.T) {
		found, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUpdateMovie(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	movie := createSampleMovie(t, repo.(*repository.MovieRepo), ctx)

	t.Run("Update full movie", func(t *testing.T) {
		movie.Title = "Updated Title"
		err := repo.Update(ctx, movie)
		require.NoError(t, err)

		updated, _ := repo.GetByID(ctx, movie.ID)
		assert.Equal(t, "Updated Title", updated.Title)
	})

	t.Run("Update non-existent movie", func(t *testing.T) {
		err := repo.Update(ctx, &domain.Movie{Model: gorm.Model{ID: 9999}, Title: "Does not exist"})
		assert.Error(t, err)
	})
}

func TestUpdateMovieFields(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	movie := createSampleMovie(t, repo.(*repository.MovieRepo), ctx)

	t.Run("Update single field", func(t *testing.T) {
		err := repo.UpdateFields(ctx, movie.ID, map[string]interface{}{"Genre": "Comedy"})
		require.NoError(t, err)

		updated, _ := repo.GetByID(ctx, movie.ID)
		assert.Equal(t, "Comedy", updated.Genre)
	})

	t.Run("Update non-existent movie", func(t *testing.T) {
		err := repo.UpdateFields(ctx, 9999, map[string]interface{}{"Genre": "Sci-fi"})
		assert.ErrorIs(t, err, repository.ErrMovieNotFound)
	})

	t.Run("Update with DB failure", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.UpdateFields(ctx, movie.ID, map[string]interface{}{"Genre": "Drama"})
		assert.Error(t, err)
	})
}

func TestDeleteMovie(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	movie := createSampleMovie(t, repo.(*repository.MovieRepo), ctx)

	t.Run("Delete existing movie", func(t *testing.T) {
		err := repo.Delete(ctx, movie.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, movie.ID)
		assert.Error(t, err)
	})

	t.Run("Delete non-existent movie", func(t *testing.T) {
		err := repo.Delete(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestGetAllMovies(t *testing.T) {
	db := setupTestMovieDB(t)
	repo := repository.NewMovieRepo(db)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		m := &domain.Movie{
			Title:     "Movie " + string(rune(i+'0')),
			StreamURL: "http://stream.url/movie" + string(rune(i+'0')) + ".mp4",
			Year:      2000 + i,
		}
		err := repo.Create(ctx, m)
		require.NoError(t, err)
	}

	t.Run("Get first 2 movies", func(t *testing.T) {
		movies, err := repo.GetAll(ctx, 2, 0)
		require.NoError(t, err)
		assert.Len(t, movies, 2)
	})

	t.Run("Get next 2 movies", func(t *testing.T) {
		movies, err := repo.GetAll(ctx, 2, 2)
		require.NoError(t, err)
		assert.Len(t, movies, 2)
	})

	t.Run("Get beyond available", func(t *testing.T) {
		movies, err := repo.GetAll(ctx, 2, 10)
		require.NoError(t, err)
		assert.Len(t, movies, 0)
	})

	t.Run("GetAll with DB failure", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		_, err := repo.GetAll(ctx, 2, 0)
		assert.Error(t, err)
	})
}
