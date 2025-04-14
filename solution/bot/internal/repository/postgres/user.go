package postgres

import (
	"bot/internal/entity"
	"bot/internal/errorz"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Save(&user).Error
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, userID int64) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.User{}, errorz.ErrUserNotFound
	case err != nil:
		return entity.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", userID).Error
}
