package reviewdetail_cache

import (
	"context"
	"fmt"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
)

type reviewDetailCommandCache struct {
	store *cache.CacheStore
}

func NewReviewDetailCommandCache(store *cache.CacheStore) *reviewDetailCommandCache {
	return &reviewDetailCommandCache{store: store}
}

func (s *reviewDetailCommandCache) DeleteReviewDetailCache(ctx context.Context, review_id int) {
	cache.DeleteFromCache(ctx, s.store, fmt.Sprintf(reviewDetailByIdCacheKey, review_id))
}
