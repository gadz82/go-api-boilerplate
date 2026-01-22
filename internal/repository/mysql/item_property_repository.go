package mysql

import (
	"context"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"gorm.io/gorm"
)

type itemPropertyRepository struct {
	db *gorm.DB
}

func NewItemPropertyRepository(db *gorm.DB) domain.ItemPropertyRepository {
	return &itemPropertyRepository{db: db}
}

func (r *itemPropertyRepository) GetAllByItemID(ctx context.Context, itemID string) ([]*domain.ItemProperty, error) {
	var itemProperties []*domain.ItemProperty
	if err := r.db.WithContext(ctx).Where("item_id = ?", itemID).Find(&itemProperties).Error; err != nil {
		return nil, err
	}
	return itemProperties, nil
}

func (r *itemPropertyRepository) GetByID(ctx context.Context, itemID string, id string) (*domain.ItemProperty, error) {
	var itemProperty domain.ItemProperty
	if err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&itemProperty, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &itemProperty, nil
}

func (r *itemPropertyRepository) Create(ctx context.Context, itemProperty *domain.ItemProperty) error {
	return r.db.WithContext(ctx).Create(itemProperty).Error
}

func (r *itemPropertyRepository) Update(ctx context.Context, itemProperty *domain.ItemProperty) error {
	return r.db.WithContext(ctx).Save(itemProperty).Error
}

func (r *itemPropertyRepository) Delete(ctx context.Context, itemID string, id string) error {
	return r.db.WithContext(ctx).Where("item_id = ?", itemID).Delete(&domain.ItemProperty{}, "id = ?", id).Error
}
