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

func setupTestUserDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepo_Create(t *testing.T) {
	db := setupTestUserDB(t)
	repo := repository.NewUserRepo(db)

	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		user := &domain.User{Email: "test@example.com", Password: "secure123"}
		err := repo.Create(ctx, user)

		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("should fail to create duplicate email", func(t *testing.T) {
		user1 := &domain.User{Email: "duplicate@example.com", Password: "abc"}
		user2 := &domain.User{Email: "duplicate@example.com", Password: "xyz"}

		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		err = repo.Create(ctx, user2)
		assert.Error(t, err)
	})
}

func TestUserRepo_FindByEmail(t *testing.T) {
	db := setupTestUserDB(t)
	repo := repository.NewUserRepo(db)
	ctx := context.Background()

	user := &domain.User{Email: "findme@example.com", Password: "123456"}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("should find existing user by email", func(t *testing.T) {
		result, err := repo.FindByEmail(ctx, "findme@example.com")

		require.NoError(t, err)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("should return error for non-existing email", func(t *testing.T) {
		result, err := repo.FindByEmail(ctx, "notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
