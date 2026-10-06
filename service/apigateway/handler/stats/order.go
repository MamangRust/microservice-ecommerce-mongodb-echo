package statshandler

import (
	"fmt"
	"net/http"

	"github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	statsmapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/stats"
	"github.com/labstack/echo/v4"
)

type orderStatsHandlerApi struct {
	clients    *StatsClients
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

type orderStatsHandleDeps struct {
	clients    *StatsClients
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

func NewOrderStatsHandleApi(params *orderStatsHandleDeps) *orderStatsHandlerApi {
	h := &orderStatsHandlerApi{
		clients:    params.clients,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	g := params.router.Group("/api/order-stats")

	g.GET("/monthly-total-revenue", params.apiHandler.Handle("order-stats-monthly-total-revenue", h.FindMonthlyTotalRevenue))
	g.GET("/yearly-total-revenue", params.apiHandler.Handle("order-stats-yearly-total-revenue", h.FindYearlyTotalRevenue))
	g.GET("/monthly-revenue", params.apiHandler.Handle("order-stats-monthly-revenue", h.FindMonthlyRevenue))
	g.GET("/yearly-revenue", params.apiHandler.Handle("order-stats-yearly-revenue", h.FindYearlyRevenue))

	g.GET("/by-merchant/monthly-total-revenue", params.apiHandler.Handle("order-stats-by-merchant-monthly-total-revenue", h.FindMonthlyTotalRevenueByMerchant))
	g.GET("/by-merchant/yearly-total-revenue", params.apiHandler.Handle("order-stats-by-merchant-yearly-total-revenue", h.FindYearlyTotalRevenueByMerchant))
	g.GET("/by-merchant/monthly-revenue", params.apiHandler.Handle("order-stats-by-merchant-monthly-revenue", h.FindMonthlyRevenueByMerchant))
	g.GET("/by-merchant/yearly-revenue", params.apiHandler.Handle("order-stats-by-merchant-yearly-revenue", h.FindYearlyRevenueByMerchant))

	return h
}

// @Security Bearer
// @Summary Monthly total revenue
// @Tags Order Stats
// @Description Monthly total order revenue
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/monthly-total-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyTotalRevenue(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:monthly-total-revenue:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStats.FindMonthlyTotalRevenue(ctx, &pb_order.FindYearMonthTotalRevenue{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderMonthlyTotalRevenue(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly total revenue
// @Tags Order Stats
// @Description Yearly total order revenue
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/yearly-total-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyTotalRevenue(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:yearly-total-revenue:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStats.FindYearlyTotalRevenue(ctx, &pb_order.FindYearTotalRevenue{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderYearlyTotalRevenue(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly revenue
// @Tags Order Stats
// @Description Monthly order aggregates (count, revenue, items sold)
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/monthly-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyRevenue(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:monthly-revenue:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderMonthly](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStats.FindMonthlyRevenue(ctx, &pb_order.FindYearOrder{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderMonthly(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly revenue
// @Tags Order Stats
// @Description Yearly order aggregates (count, revenue, items sold)
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/yearly-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyRevenue(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:yearly-revenue:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderYearly](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStats.FindYearlyRevenue(ctx, &pb_order.FindYearOrder{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderYearly(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly total revenue by merchant
// @Tags Order Stats
// @Description Monthly total order revenue for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/by-merchant/monthly-total-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyTotalRevenueByMerchant(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:by-merchant:monthly-total-revenue:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &pb_order.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderMonthlyTotalRevenue(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly total revenue by merchant
// @Tags Order Stats
// @Description Yearly total order revenue for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/by-merchant/yearly-total-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyTotalRevenueByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:by-merchant:yearly-total-revenue:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStatsByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &pb_order.FindYearTotalRevenueByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderYearlyTotalRevenue(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly revenue by merchant
// @Tags Order Stats
// @Description Monthly order aggregates for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/by-merchant/monthly-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyRevenueByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:by-merchant:monthly-revenue:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderMonthly](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStatsByMerchant.FindMonthlyRevenueByMerchant(ctx, &pb_order.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderMonthly(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly revenue by merchant
// @Tags Order Stats
// @Description Yearly order aggregates for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/order-stats/by-merchant/yearly-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyRevenueByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:order-stats:by-merchant:yearly-revenue:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseOrderYearly](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.OrderStatsByMerchant.FindYearlyRevenueByMerchant(ctx, &pb_order.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrderYearly(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}
