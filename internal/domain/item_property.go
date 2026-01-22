package domain

import "context"

type ItemProperty struct {
	ID     string `jsonapi:"primary,item_properties" json:"id" gorm:"primaryKey;type:char(36)" validate:"omitempty,uuid4"`
	ItemID string `jsonapi:"attr,item_id" json:"item_id" gorm:"index;type:char(36)" validate:"omitempty,uuid4"`
	Name   string `jsonapi:"attr,name" json:"name" gorm:"index" validate:"required,min=1,max=255"`
	Value  string `jsonapi:"attr,value" json:"value" validate:"required,max=1000"`
}

type ItemPropertyRepository interface {
	GetAllByItemID(ctx context.Context, itemID string) ([]*ItemProperty, error)
	GetByID(ctx context.Context, itemID string, id string) (*ItemProperty, error)
	Create(ctx context.Context, itemProperty *ItemProperty) error
	Update(ctx context.Context, itemProperty *ItemProperty) error
	Delete(ctx context.Context, itemID string, id string) error
}

type ItemPropertyService interface {
	GetItemPropertiesByItemID(ctx context.Context, itemID string) ([]*ItemProperty, error)
	GetItemPropertyByID(ctx context.Context, itemID string, id string) (*ItemProperty, error)
	CreateItemProperty(ctx context.Context, itemProperty *ItemProperty) error
	UpdateItemProperty(ctx context.Context, itemProperty *ItemProperty) error
	DeleteItemProperty(ctx context.Context, itemID string, id string) error
}
