package merchantbusinesshandler

import (
	merchantbusiness_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/merchant_business"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_business"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	merchantapimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_business"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantBusiness struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterMerchantBusinessHandler(deps *DepsMerchantBusiness) {
	mapper := apimapper.NewMerchantBusinessResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantbusiness_cache.NewMerchantBusinessMencache(deps.CacheStore)

	NewMerchantBusinessQueryHandleApi(&merchantBusinessQueryHandleDeps{
		client: pb_merchant_business.NewMerchantBusinessQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantBusinessCommandHandleApi(&merchantBusinessCommandHandleDeps{
		client:         pb_merchant_business.NewMerchantBusinessCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
