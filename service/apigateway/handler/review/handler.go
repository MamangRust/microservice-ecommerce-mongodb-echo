package reviewhandler

import (
	review_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/review"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/review"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsReview struct {
	Client        *grpc.ClientConn
	E             *echo.Echo
	Logger        logger.LoggerInterface
	Cache         *cache.CacheStore
	Observability observability.TraceLoggerObservability
}

func RegisterReviewHandler(deps *DepsReview) {
	mapper := apimapper.NewReviewResponseMapper()
	cache := review_cache.NewReviewMencache(deps.Cache)

	NewReviewQueryHandleApi(&reviewQueryHandleDeps{
		client:        pb_review.NewReviewQueryServiceClient(deps.Client),
		router:        deps.E,
		logger:        deps.Logger,
		mapper:        mapper.QueryMapper(),
		cache:         cache.QueryCache(),
		observability: deps.Observability,
	})

	NewReviewCommandHandleApi(&reviewCommandHandleDeps{
		client:        pb_review.NewReviewCommandServiceClient(deps.Client),
		router:        deps.E,
		logger:        deps.Logger,
		mapper:        mapper.CommandMapper(),
		queryMapper:   mapper.QueryMapper(),
		cache:         cache.CommandCache(),
		observability: deps.Observability,
	})
}
