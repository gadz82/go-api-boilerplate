package domain

import (
	"context"
	"time"
)

type Item struct {
	ID             string          `jsonapi:"primary,items" json:"id" gorm:"primaryKey;type:char(36)" validate:"omitempty,uuid4"`
	Title          string          `jsonapi:"attr,title" json:"title" gorm:"index" validate:"required,min=1,max=255"`
	Description    string          `jsonapi:"attr,description" json:"description" validate:"max=1000"`
	CreatedAt      *time.Time      `jsonapi:"attr,created_at,iso8601" json:"created_at,omitempty" gorm:"type:timestamp;default:null"`
	UpdatedAt      time.Time       `jsonapi:"attr,updated_at,iso8601" json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ItemProperties []*ItemProperty `jsonapi:"relation,item_properties" json:"item_properties,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
}

type ItemRepository interface {
	GetAll(ctx context.Context) ([]*Item, error)
	GetByID(ctx context.Context, id string) (*Item, error)
	Create(ctx context.Context, item *Item) error
	Update(ctx context.Context, item *Item) error
	Delete(ctx context.Context, id string) error
}

type ItemService interface {
	GetAllItems(ctx context.Context) ([]*Item, error)
	GetItemByID(ctx context.Context, id string) (*Item, error)
	CreateItem(ctx context.Context, item *Item) error
	UpdateItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, id string) error
}
