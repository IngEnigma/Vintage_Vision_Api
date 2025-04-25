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

func setupProfileTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Profile{})
	require.NoError(t, err)

	return db
}

func TestProfileRepo_Create(t *testing.T) {
	db := setupProfileTestDB(t)
	repo := repository.NewProfileRepo(db)

	ctx := context.Background()

	t.Run("should create profile successfully", func(t *testing.T) {
		profile := &domain.Profile{Name: "Pepe", AvatarURL: "url", UserID: 1}
		err := repo.Create(ctx, profile)

		require.NoError(t, err)
		assert.NotZero(t, profile.ID)
	})

	t.Run("should fail validation (empty name)", func(t *testing.T) {
		profile := &domain.Profile{Name: "", UserID: 1}
		err := repo.Create(ctx, profile)

		assert.Error(t, err)
	})
}

func TestProfileRepo_FindByUser(t *testing.T) {
	db := setupProfileTestDB(t)
	repo := repository.NewProfileRepo(db)
	ctx := context.Background()

	profile1 := &domain.Profile{Name: "User1", UserID: 100}
	profile2 := &domain.Profile{Name: "User2", UserID: 100}
	repo.Create(ctx, profile1)
	repo.Create(ctx, profile2)

	t.Run("should return all profiles for user", func(t *testing.T) {
		profiles, err := repo.FindByUser(ctx, 100)

		require.NoError(t, err)
		assert.Len(t, profiles, 2)
	})

	t.Run("should return empty slice if no profiles", func(t *testing.T) {
		profiles, err := repo.FindByUser(ctx, 999)
		require.NoError(t, err)
		assert.Len(t, profiles, 0)
	})
}

func TestProfileRepo_FindByIDAndUser(t *testing.T) {
	db := setupProfileTestDB(t)
	repo := repository.NewProfileRepo(db)
	ctx := context.Background()

	profile := &domain.Profile{Name: "Solo", UserID: 22}
	repo.Create(ctx, profile)

	t.Run("should find profile by ID and user", func(t *testing.T) {
		result, err := repo.FindByIDAndUser(ctx, profile.ID, 22)

		require.NoError(t, err)
		assert.Equal(t, profile.Name, result.Name)
	})

	t.Run("should return error if not found", func(t *testing.T) {
		result, err := repo.FindByIDAndUser(ctx, 999, 22)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProfileRepo_Update(t *testing.T) {
	db := setupProfileTestDB(t)
	repo := repository.NewProfileRepo(db)
	ctx := context.Background()

	profile := &domain.Profile{Name: "Before", UserID: 88}
	repo.Create(ctx, profile)

	t.Run("should update profile successfully", func(t *testing.T) {
		profile.Name = "After"
		err := repo.Update(ctx, profile)

		require.NoError(t, err)

		updated, _ := repo.FindByIDAndUser(ctx, profile.ID, 88)
		assert.Equal(t, "After", updated.Name)
	})
}

func TestProfileRepo_DeleteByID(t *testing.T) {
	db := setupProfileTestDB(t)
	repo := repository.NewProfileRepo(db)
	ctx := context.Background()

	profile := &domain.Profile{Name: "To Delete", UserID: 55}
	repo.Create(ctx, profile)

	t.Run("should delete profile successfully", func(t *testing.T) {
		err := repo.DeleteByID(ctx, profile.ID, 55)

		require.NoError(t, err)

		result, err := repo.FindByIDAndUser(ctx, profile.ID, 55)
		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("should return error if profile not found", func(t *testing.T) {
		err := repo.DeleteByID(ctx, 999, 55)
		assert.Error(t, err)
	})
}
