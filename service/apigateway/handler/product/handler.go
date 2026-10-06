package producthandler

import (
	product_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/product"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/upload_image"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/product"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsProduct struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	Upload     upload_image.ImageUploads
	ApiHandler errors.ApiHandler
}

func RegisterProductHandler(deps *DepsProduct) {
	mapper := apimapper.NewProductResponseMapper()
	cache := product_cache.NewProductMencache(deps.CacheStore)

	queryClient := pb_product.NewProductQueryServiceClient(deps.Client)

	NewProductQueryHandleApi(&productQueryHandleDeps{
		client:     queryClient,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.QueryMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewProductCommandHandleApi(&productCommandHandleDeps{
		client:       pb_product.NewProductCommandServiceClient(deps.Client),
		router:       deps.E,
		logger:       deps.Logger,
		mapper:       mapper.CommandMapper(),
		cache:        cache,
		upload_image: deps.Upload,
		apiHandler:   deps.ApiHandler,
	})

}
