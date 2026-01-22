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
	// Cache key prefixes
	itemCacheKeyPrefix = "item:"
	itemsListCacheKey  = "items:list"

	// Default cache TTL
	defaultCacheTTL = 5 * time.Minute
)

type itemService struct {
	itemRepo  domain.ItemRepository
	cacheRepo domain.CacheRepository
}

func NewItemService(itemRepo domain.ItemRepository, cacheRepo domain.CacheRepository) domain.ItemService {
	return &itemService{
		itemRepo:  itemRepo,
		cacheRepo: cacheRepo,
	}
}

// GetAllItems retrieves all items with lazy caching strategy.
// It first checks the cache, and if not found, fetches from the database and caches the result.
func (s *itemService) GetAllItems(ctx context.Context) ([]*domain.Item, error) {
	// Try to get from cache first
	cached, err := s.cacheRepo.Get(ctx, itemsListCacheKey)
	if err == nil && cached != "" {
		var items []*domain.Item
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			log.Printf("Cache hit for items list")
			return items, nil
		}
	}

	// Cache miss - fetch from database
	log.Printf("Cache miss for items list, fetching from database")
	items, err := s.itemRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(items); err == nil {
		if err := s.cacheRepo.Set(ctx, itemsListCacheKey, string(data), defaultCacheTTL); err != nil {
			log.Printf("Failed to cache items list: %v", err)
		}
	}

	return items, nil
}

// GetItemByID retrieves an item by ID with lazy caching strategy.
// It first checks the cache, and if not found, fetches from the database and caches the result.
func (s *itemService) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	cacheKey := fmt.Sprintf("%s%s", itemCacheKeyPrefix, id)

	// Try to get from cache first
	cached, err := s.cacheRepo.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var item domain.Item
		if err := json.Unmarshal([]byte(cached), &item); err == nil {
			log.Printf("Cache hit for item %s", id)
			return &item, nil
		}
	}

	// Cache miss - fetch from database
	log.Printf("Cache miss for item %s, fetching from database", id)
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(item); err == nil {
		if err := s.cacheRepo.Set(ctx, cacheKey, string(data), defaultCacheTTL); err != nil {
			log.Printf("Failed to cache item %s: %v", id, err)
		}
	}

	return item, nil
}

// CreateItem creates a new item and invalidates the items list cache.
func (s *itemService) CreateItem(ctx context.Context, item *domain.Item) error {
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return err
	}

	// Invalidate the items list cache since a new item was added
	if err := s.cacheRepo.Delete(ctx, itemsListCacheKey); err != nil {
		log.Printf("Failed to invalidate items list cache: %v", err)
	}

	return nil
}

// UpdateItem updates an item and invalidates both the single item cache and the items list cache.
func (s *itemService) UpdateItem(ctx context.Context, item *domain.Item) error {
	if err := s.itemRepo.Update(ctx, item); err != nil {
		return err
	}

	// Invalidate the single item cache
	cacheKey := fmt.Sprintf("%s%s", itemCacheKeyPrefix, item.ID)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate item cache %s: %v", item.ID, err)
	}

	// Invalidate the items list cache since an item was updated
	if err := s.cacheRepo.Delete(ctx, itemsListCacheKey); err != nil {
		log.Printf("Failed to invalidate items list cache: %v", err)
	}

	return nil
}

// DeleteItem deletes an item and invalidates both the single item cache and the items list cache.
func (s *itemService) DeleteItem(ctx context.Context, id string) error {
	if err := s.itemRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate the single item cache
	cacheKey := fmt.Sprintf("%s%s", itemCacheKeyPrefix, id)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate item cache %s: %v", id, err)
	}

	// Invalidate the items list cache since an item was deleted
	if err := s.cacheRepo.Delete(ctx, itemsListCacheKey); err != nil {
		log.Printf("Failed to invalidate items list cache: %v", err)
	}

	return nil
}
