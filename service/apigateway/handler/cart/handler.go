package carthandler

import (
	cart_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/cart"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/cart"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/cart"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsCart struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterCartHandler(deps *DepsCart) {
	mapper := apimapper.NewCartResponseMapper()
	cache := cart_cache.NewCartMencache(deps.CacheStore)

	NewCartQueryHandleApi(&cartQueryHandleDeps{
		client: pb_cart.NewCartQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewCartCommandHandleApi(&cartCommandHandleDeps{
		client: pb_cart.NewCartCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})
}
