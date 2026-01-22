package items

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

const (
	// Cache key prefixes for item properties
	itemPropertyCacheKeyPrefix    = "item_property:"
	itemPropertiesListCacheKeyFmt = "item_properties:list:%s"

	// Default cache TTL for item properties
	defaultPropertyCacheTTL = 5 * time.Minute
)

type itemPropertyService struct {
	itemPropertyRepo domain.ItemPropertyRepository
	cacheRepo        domain.CacheRepository
}

func NewItemPropertyService(itemPropertyRepo domain.ItemPropertyRepository, cacheRepo domain.CacheRepository) domain.ItemPropertyService {
	return &itemPropertyService{
		itemPropertyRepo: itemPropertyRepo,
		cacheRepo:        cacheRepo,
	}
}

// GetItemPropertiesByItemID retrieves all properties for an item with lazy caching strategy.
func (s *itemPropertyService) GetItemPropertiesByItemID(ctx context.Context, itemID string) ([]*domain.ItemProperty, error) {
	cacheKey := fmt.Sprintf(itemPropertiesListCacheKeyFmt, itemID)

	// Try to get from cache first
	cached, err := s.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var properties []*domain.ItemProperty
		if err := json.Unmarshal([]byte(cached), &properties); err == nil {
			log.Printf("Cache hit for item properties list (item: %s)", itemID)
			return properties, nil
		}
	}

	// Cache miss - fetch from database
	log.Printf("Cache miss for item properties list (item: %s), fetching from database", itemID)
	properties, err := s.itemPropertyRepo.GetAllByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(properties); err == nil {
		if err := s.cacheRepo.Set(ctx, cacheKey, string(data), defaultPropertyCacheTTL); err != nil {
			log.Printf("Failed to cache item properties list (item: %s): %v", itemID, err)
		}
	}

	return properties, nil
}

// GetItemPropertyByID retrieves a single item property with lazy caching strategy.
func (s *itemPropertyService) GetItemPropertyByID(ctx context.Context, itemID string, id string) (*domain.ItemProperty, error) {
	cacheKey := fmt.Sprintf("%s%s:%s", itemPropertyCacheKeyPrefix, itemID, id)

	// Try to get from cache first
	cached, err := s.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var property domain.ItemProperty
		if err := json.Unmarshal([]byte(cached), &property); err == nil {
			log.Printf("Cache hit for item property %s (item: %s)", id, itemID)
			return &property, nil
		}
	}

	// Cache miss - fetch from database
	log.Printf("Cache miss for item property %s (item: %s), fetching from database", id, itemID)
	property, err := s.itemPropertyRepo.GetByID(ctx, itemID, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(property); err == nil {
		if err := s.cacheRepo.Set(ctx, cacheKey, string(data), defaultPropertyCacheTTL); err != nil {
			log.Printf("Failed to cache item property %s (item: %s): %v", id, itemID, err)
		}
	}

	return property, nil
}

// CreateItemProperty creates a new item property and invalidates the properties list cache.
func (s *itemPropertyService) CreateItemProperty(ctx context.Context, itemProperty *domain.ItemProperty) error {
	if err := s.itemPropertyRepo.Create(ctx, itemProperty); err != nil {
		return err
	}

	// Invalidate the item properties list cache since a new property was added
	listCacheKey := fmt.Sprintf(itemPropertiesListCacheKeyFmt, itemProperty.ItemID)
	if err := s.cacheRepo.Delete(ctx, listCacheKey); err != nil {
		log.Printf("Failed to invalidate item properties list cache (item: %s): %v", itemProperty.ItemID, err)
	}

	return nil
}

// UpdateItemProperty updates an item property and invalidates both the single property cache and the list cache.
func (s *itemPropertyService) UpdateItemProperty(ctx context.Context, itemProperty *domain.ItemProperty) error {
	if err := s.itemPropertyRepo.Update(ctx, itemProperty); err != nil {
		return err
	}

	// Invalidate the single property cache
	cacheKey := fmt.Sprintf("%s%s:%s", itemPropertyCacheKeyPrefix, itemProperty.ItemID, itemProperty.ID)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate item property cache %s (item: %s): %v", itemProperty.ID, itemProperty.ItemID, err)
	}

	// Invalidate the item properties list cache since a property was updated
	listCacheKey := fmt.Sprintf(itemPropertiesListCacheKeyFmt, itemProperty.ItemID)
	if err := s.cacheRepo.Delete(ctx, listCacheKey); err != nil {
		log.Printf("Failed to invalidate item properties list cache (item: %s): %v", itemProperty.ItemID, err)
	}

	return nil
}

// DeleteItemProperty deletes an item property and invalidates both the single property cache and the list cache.
func (s *itemPropertyService) DeleteItemProperty(ctx context.Context, itemID string, id string) error {
	if err := s.itemPropertyRepo.Delete(ctx, itemID, id); err != nil {
		return err
	}

	// Invalidate the single property cache
	cacheKey := fmt.Sprintf("%s%s:%s", itemPropertyCacheKeyPrefix, itemID, id)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate item property cache %s (item: %s): %v", id, itemID, err)
	}

	// Invalidate the item properties list cache since a property was deleted
	listCacheKey := fmt.Sprintf(itemPropertiesListCacheKeyFmt, itemID)
	if err := s.cacheRepo.Delete(ctx, listCacheKey); err != nil {
		log.Printf("Failed to invalidate item properties list cache (item: %s): %v", itemID, err)
	}

	return nil
}
