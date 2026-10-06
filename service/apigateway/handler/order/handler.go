package orderhandler

import (
	order_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/order"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/order"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsOrder struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterOrderHandler(deps *DepsOrder) {
	mapper := apimapper.NewOrderResponseMapper()
	cache := order_cache.OrderNewMencache(deps.CacheStore)

	queryClient := pb_order.NewOrderQueryServiceClient(deps.Client)

	NewOrderQueryHandleApi(&orderQueryHandleDeps{
		client: queryClient,
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewOrderCommandHandleApi(&orderCommandHandleDeps{
		client: pb_order.NewOrderCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})
}
