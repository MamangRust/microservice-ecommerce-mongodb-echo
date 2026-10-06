package statshandler

import (
	"fmt"
	"net/http"

	"github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	statsmapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/stats"
	"github.com/labstack/echo/v4"
)

type categoryStatsHandlerApi struct {
	clients    *StatsClients
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

type categoryStatsHandleDeps struct {
	clients    *StatsClients
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

func NewCategoryStatsHandleApi(params *categoryStatsHandleDeps) *categoryStatsHandlerApi {
	h := &categoryStatsHandlerApi{
		clients:    params.clients,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	g := params.router.Group("/api/category-stats")

	g.GET("/monthly-total-prices", params.apiHandler.Handle("category-stats-monthly-total-prices", h.FindMonthlyTotalPrices))
	g.GET("/yearly-total-prices", params.apiHandler.Handle("category-stats-yearly-total-prices", h.FindYearlyTotalPrices))
	g.GET("/month-price", params.apiHandler.Handle("category-stats-month-price", h.FindMonthPrice))
	g.GET("/year-price", params.apiHandler.Handle("category-stats-year-price", h.FindYearPrice))

	g.GET("/by-merchant/monthly-total-prices", params.apiHandler.Handle("category-stats-by-merchant-monthly-total-prices", h.FindMonthlyTotalPricesByMerchant))
	g.GET("/by-merchant/yearly-total-prices", params.apiHandler.Handle("category-stats-by-merchant-yearly-total-prices", h.FindYearlyTotalPricesByMerchant))
	g.GET("/by-merchant/month-price", params.apiHandler.Handle("category-stats-by-merchant-month-price", h.FindMonthPriceByMerchant))
	g.GET("/by-merchant/year-price", params.apiHandler.Handle("category-stats-by-merchant-year-price", h.FindYearPriceByMerchant))

	g.GET("/by-id/monthly-total-prices", params.apiHandler.Handle("category-stats-by-id-monthly-total-prices", h.FindMonthlyTotalPricesById))
	g.GET("/by-id/yearly-total-prices", params.apiHandler.Handle("category-stats-by-id-yearly-total-prices", h.FindYearlyTotalPricesById))
	g.GET("/by-id/month-price", params.apiHandler.Handle("category-stats-by-id-month-price", h.FindMonthPriceById))
	g.GET("/by-id/year-price", params.apiHandler.Handle("category-stats-by-id-year-price", h.FindYearPriceById))

	return h
}

// @Security Bearer
// @Summary Monthly total prices
// @Tags Category Stats
// @Description Monthly total prices across all categories
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/monthly-total-prices [get]
func (h *categoryStatsHandlerApi) FindMonthlyTotalPrices(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:monthly-total-prices:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStats.FindMonthlyTotalPrices(ctx, &pb_category.FindYearMonthTotalPrices{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly total prices
// @Tags Category Stats
// @Description Yearly total prices across all categories
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/yearly-total-prices [get]
func (h *categoryStatsHandlerApi) FindYearlyTotalPrices(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:yearly-total-prices:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStats.FindYearlyTotalPrices(ctx, &pb_category.FindYearTotalPrices{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly prices
// @Tags Category Stats
// @Description Monthly prices per category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/month-price [get]
func (h *categoryStatsHandlerApi) FindMonthPrice(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:month-price:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStats.FindMonthPrice(ctx, &pb_category.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly prices
// @Tags Category Stats
// @Description Yearly prices per category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/year-price [get]
func (h *categoryStatsHandlerApi) FindYearPrice(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:year-price:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStats.FindYearPrice(ctx, &pb_category.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly total prices by merchant
// @Tags Category Stats
// @Description Monthly total prices for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-merchant/monthly-total-prices [get]
func (h *categoryStatsHandlerApi) FindMonthlyTotalPricesByMerchant(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-merchant:monthly-total-prices:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsByMerchant.FindMonthlyTotalPricesByMerchant(ctx, &pb_category.FindYearMonthTotalPriceByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly total prices by merchant
// @Tags Category Stats
// @Description Yearly total prices for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-merchant/yearly-total-prices [get]
func (h *categoryStatsHandlerApi) FindYearlyTotalPricesByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-merchant:yearly-total-prices:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsByMerchant.FindYearlyTotalPricesByMerchant(ctx, &pb_category.FindYearTotalPriceByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly prices by merchant
// @Tags Category Stats
// @Description Monthly prices per category for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-merchant/month-price [get]
func (h *categoryStatsHandlerApi) FindMonthPriceByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-merchant:month-price:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &pb_category.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly prices by merchant
// @Tags Category Stats
// @Description Yearly prices per category for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-merchant/year-price [get]
func (h *categoryStatsHandlerApi) FindYearPriceByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-merchant:year-price:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsByMerchant.FindYearPriceByMerchant(ctx, &pb_category.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly total prices by category
// @Tags Category Stats
// @Description Monthly total prices for a category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-id/monthly-total-prices [get]
func (h *categoryStatsHandlerApi) FindMonthlyTotalPricesById(c echo.Context) error {
	year, month, categoryID, err := parseYearMonthCategory(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-id:monthly-total-prices:%d:%d:%d", statsCacheNamespace, year, month, categoryID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsById.FindMonthlyTotalPricesById(ctx, &pb_category.FindYearMonthTotalPriceById{
		Year:       int32(year),
		Month:      int32(month),
		CategoryId: int32(categoryID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly total prices by category
// @Tags Category Stats
// @Description Yearly total prices for a category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-id/yearly-total-prices [get]
func (h *categoryStatsHandlerApi) FindYearlyTotalPricesById(c echo.Context) error {
	year, categoryID, err := parseYearCategory(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-id:yearly-total-prices:%d:%d", statsCacheNamespace, year, categoryID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsById.FindYearlyTotalPricesById(ctx, &pb_category.FindYearTotalPriceById{
		Year:       int32(year),
		CategoryId: int32(categoryID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearlyTotalPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly prices by category
// @Tags Category Stats
// @Description Monthly prices for a category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-id/month-price [get]
func (h *categoryStatsHandlerApi) FindMonthPriceById(c echo.Context) error {
	year, categoryID, err := parseYearCategory(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-id:month-price:%d:%d", statsCacheNamespace, year, categoryID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsById.FindMonthPriceById(ctx, &pb_category.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(categoryID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryMonthPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly prices by category
// @Tags Category Stats
// @Description Yearly prices for a category
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/category-stats/by-id/year-price [get]
func (h *categoryStatsHandlerApi) FindYearPriceById(c echo.Context) error {
	year, categoryID, err := parseYearCategory(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:category-stats:by-id:year-price:%d:%d", statsCacheNamespace, year, categoryID)

	if cached, found := cache.GetFromCache[response.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.CategoryStatsById.FindYearPriceById(ctx, &pb_category.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(categoryID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCategoryYearPrice(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}
