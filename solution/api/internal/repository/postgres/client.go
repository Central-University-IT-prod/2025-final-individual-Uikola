package postgres

import (
	"api/internal/entity"
	"api/internal/errorz"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ClientRepository handles database operations for the Client entity.
type ClientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new instance of ClientRepository.
func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

// CreateMany inserts multiple client records into the database.
// If a client with the same primary key (ClientID) already exists, it updates the existing record.
func (r *ClientRepository) CreateMany(ctx context.Context, clients []entity.Client) ([]entity.Client, error) {
	err := r.db.WithContext(ctx).Save(&clients).Error
	return clients, err
}

// GetByID retrieves a client by its ClientID from the database.
func (r *ClientRepository) GetByID(ctx context.Context, clientID string) (entity.Client, error) {
	var client entity.Client
	err := r.db.WithContext(ctx).Where("client_id = ?", clientID).First(&client).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return entity.Client{}, errorz.ErrClientNotFound
	case err != nil:
		return entity.Client{}, err
	}

	return client, nil
}
