package shippingaddresshandler

import (
	shippingaddress_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/shipping_address"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/shipping_address"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsShippingAddress struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
	Cache  *cache.CacheStore
}

func RegisterShippingAddressHandler(deps *DepsShippingAddress) {
	mapper := apimapper.NewShippingAddressResponseMapper()
	cache := shippingaddress_cache.NewShippingAddressMencache(deps.Cache)

	NewShippingAddressQueryHandleApi(&shippingAddressQueryHandleDeps{
		client: pb_shipping_address.NewShippingQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache.QueryCache(),
	})

	NewShippingAddressCommandHandleApi(&shippingAddressCommandHandleDeps{
		client: pb_shipping_address.NewShippingCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache.CommandCache(),
	})
}
