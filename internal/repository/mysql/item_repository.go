package mysql

import (
	"context"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domain.ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) GetAll(ctx context.Context) ([]*domain.Item, error) {
	var items []*domain.Item
	db := r.db.WithContext(ctx)
	if ctx.Value("include_properties") == true {
		db = db.Preload("ItemProperties")
	}
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *itemRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	var item domain.Item
	db := r.db.WithContext(ctx)
	if ctx.Value("include_properties") == true {
		db = db.Preload("ItemProperties")
	}
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) Create(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *itemRepository) Update(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *itemRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Item{}, "id = ?", id).Error
}
