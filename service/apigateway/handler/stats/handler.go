package statshandler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	statsmapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/stats"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// StatsClients bundles every stats gRPC client exposed by stats_reader. All of
// them share one connection (ServiceConnections.StatsReader).
type StatsClients struct {
	CategoryStats              pb_category.CategoryStatsServiceClient
	CategoryStatsByMerchant    pb_category.CategoryStatsByMerchantServiceClient
	CategoryStatsById          pb_category.CategoryStatsByIdServiceClient
	OrderStats                 pb_order.OrderStatsServiceClient
	OrderStatsByMerchant       pb_order.OrderStatsByMerchantServiceClient
	TransactionStats           pb_transaction.TransactionStatsServiceClient
	TransactionStatsByMerchant pb_transaction.TransactionStatsByMerchantServiceClient
}

type DepsStats struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	ApiHandler errors.ApiHandler
}

func RegisterStatsHandler(deps *DepsStats) {
	mapper := statsmapper.NewStatsResponseMapper()

	clients := &StatsClients{
		CategoryStats:              pb_category.NewCategoryStatsServiceClient(deps.Client),
		CategoryStatsByMerchant:    pb_category.NewCategoryStatsByMerchantServiceClient(deps.Client),
		CategoryStatsById:          pb_category.NewCategoryStatsByIdServiceClient(deps.Client),
		OrderStats:                 pb_order.NewOrderStatsServiceClient(deps.Client),
		OrderStatsByMerchant:       pb_order.NewOrderStatsByMerchantServiceClient(deps.Client),
		TransactionStats:           pb_transaction.NewTransactionStatsServiceClient(deps.Client),
		TransactionStatsByMerchant: pb_transaction.NewTransactionStatsByMerchantServiceClient(deps.Client),
	}

	NewCategoryStatsHandleApi(&categoryStatsHandleDeps{
		clients:    clients,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper,
		cache:      deps.CacheStore,
		apiHandler: deps.ApiHandler,
	})

	NewOrderStatsHandleApi(&orderStatsHandleDeps{
		clients:    clients,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper,
		cache:      deps.CacheStore,
		apiHandler: deps.ApiHandler,
	})

	NewTransactionStatsHandleApi(&transactionStatsHandleDeps{
		clients:    clients,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper,
		cache:      deps.CacheStore,
		apiHandler: deps.ApiHandler,
	})
}
