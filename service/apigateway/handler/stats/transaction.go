package statshandler

import (
	"fmt"
	"net/http"

	"github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	statsmapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/stats"
	"github.com/labstack/echo/v4"
)

type transactionStatsHandlerApi struct {
	clients    *StatsClients
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

type transactionStatsHandleDeps struct {
	clients    *StatsClients
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     statsmapper.StatsResponseMapper
	cache      *cache.CacheStore
	apiHandler errors.ApiHandler
}

func NewTransactionStatsHandleApi(params *transactionStatsHandleDeps) *transactionStatsHandlerApi {
	h := &transactionStatsHandlerApi{
		clients:    params.clients,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	g := params.router.Group("/api/transaction-stats")

	g.GET("/month-status-success", params.apiHandler.Handle("transaction-stats-month-status-success", h.GetMonthlyAmountSuccess))
	g.GET("/year-status-success", params.apiHandler.Handle("transaction-stats-year-status-success", h.GetYearlyAmountSuccess))
	g.GET("/month-status-failed", params.apiHandler.Handle("transaction-stats-month-status-failed", h.GetMonthlyAmountFailed))
	g.GET("/year-status-failed", params.apiHandler.Handle("transaction-stats-year-status-failed", h.GetYearlyAmountFailed))
	g.GET("/month-method-success", params.apiHandler.Handle("transaction-stats-month-method-success", h.GetMonthlyMethodSuccess))
	g.GET("/year-method-success", params.apiHandler.Handle("transaction-stats-year-method-success", h.GetYearlyMethodSuccess))
	g.GET("/month-method-failed", params.apiHandler.Handle("transaction-stats-month-method-failed", h.GetMonthlyMethodFailed))
	g.GET("/year-method-failed", params.apiHandler.Handle("transaction-stats-year-method-failed", h.GetYearlyMethodFailed))

	g.GET("/by-merchant/month-status-success", params.apiHandler.Handle("transaction-stats-by-merchant-month-status-success", h.GetMonthlyAmountSuccessByMerchant))
	g.GET("/by-merchant/year-status-success", params.apiHandler.Handle("transaction-stats-by-merchant-year-status-success", h.GetYearlyAmountSuccessByMerchant))
	g.GET("/by-merchant/month-status-failed", params.apiHandler.Handle("transaction-stats-by-merchant-month-status-failed", h.GetMonthlyAmountFailedByMerchant))
	g.GET("/by-merchant/year-status-failed", params.apiHandler.Handle("transaction-stats-by-merchant-year-status-failed", h.GetYearlyAmountFailedByMerchant))
	g.GET("/by-merchant/month-method-success", params.apiHandler.Handle("transaction-stats-by-merchant-month-method-success", h.GetMonthlyMethodByMerchantSuccess))
	g.GET("/by-merchant/year-method-success", params.apiHandler.Handle("transaction-stats-by-merchant-year-method-success", h.GetYearlyMethodByMerchantSuccess))
	g.GET("/by-merchant/month-method-failed", params.apiHandler.Handle("transaction-stats-by-merchant-month-method-failed", h.GetMonthlyMethodByMerchantFailed))
	g.GET("/by-merchant/year-method-failed", params.apiHandler.Handle("transaction-stats-by-merchant-year-method-failed", h.GetYearlyMethodByMerchantFailed))

	return h
}

// @Security Bearer
// @Summary Monthly successful transactions
// @Tags Transaction Stats
// @Description Monthly amount of successful transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/month-status-success [get]
func (h *transactionStatsHandlerApi) GetMonthlyAmountSuccess(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:month-status-success:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetMonthlyAmountSuccess(ctx, &pb_transaction.MonthAmountTransactionRequest{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthSuccess(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly successful transactions
// @Tags Transaction Stats
// @Description Yearly amount of successful transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/year-status-success [get]
func (h *transactionStatsHandlerApi) GetYearlyAmountSuccess(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:year-status-success:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetYearlyAmountSuccess(ctx, &pb_transaction.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearSuccess(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly failed transactions
// @Tags Transaction Stats
// @Description Monthly amount of failed transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/month-status-failed [get]
func (h *transactionStatsHandlerApi) GetMonthlyAmountFailed(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:month-status-failed:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetMonthlyAmountFailed(ctx, &pb_transaction.MonthAmountTransactionRequest{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthFailed(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly failed transactions
// @Tags Transaction Stats
// @Description Yearly amount of failed transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/year-status-failed [get]
func (h *transactionStatsHandlerApi) GetYearlyAmountFailed(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:year-status-failed:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearFailed](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetYearlyAmountFailed(ctx, &pb_transaction.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearFailed(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly successful transactions by payment method
// @Tags Transaction Stats
// @Description Monthly payment-method breakdown of successful transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/month-method-success [get]
func (h *transactionStatsHandlerApi) GetMonthlyMethodSuccess(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:month-method-success:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetMonthlyTransactionMethodSuccess(ctx, &pb_transaction.MonthMethodTransactionRequest{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly successful transactions by payment method
// @Tags Transaction Stats
// @Description Yearly payment-method breakdown of successful transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/year-method-success [get]
func (h *transactionStatsHandlerApi) GetYearlyMethodSuccess(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:year-method-success:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetYearlyTransactionMethodSuccess(ctx, &pb_transaction.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly failed transactions by payment method
// @Tags Transaction Stats
// @Description Monthly payment-method breakdown of failed transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/month-method-failed [get]
func (h *transactionStatsHandlerApi) GetMonthlyMethodFailed(c echo.Context) error {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:month-method-failed:%d:%d", statsCacheNamespace, year, month)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetMonthlyTransactionMethodFailed(ctx, &pb_transaction.MonthMethodTransactionRequest{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly failed transactions by payment method
// @Tags Transaction Stats
// @Description Yearly payment-method breakdown of failed transactions
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/year-method-failed [get]
func (h *transactionStatsHandlerApi) GetYearlyMethodFailed(c echo.Context) error {
	year, err := parseYear(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:year-method-failed:%d", statsCacheNamespace, year)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStats.GetYearlyTransactionMethodFailed(ctx, &pb_transaction.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly successful transactions by merchant
// @Tags Transaction Stats
// @Description Monthly amount of successful transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/month-status-success [get]
func (h *transactionStatsHandlerApi) GetMonthlyAmountSuccessByMerchant(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:month-status-success:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, &pb_transaction.MonthAmountTransactionMerchantRequest{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthSuccess(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly successful transactions by merchant
// @Tags Transaction Stats
// @Description Yearly amount of successful transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/year-status-success [get]
func (h *transactionStatsHandlerApi) GetYearlyAmountSuccessByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:year-status-success:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetYearlyAmountSuccessByMerchant(ctx, &pb_transaction.YearAmountTransactionMerchantRequest{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearSuccess(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly failed transactions by merchant
// @Tags Transaction Stats
// @Description Monthly amount of failed transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/month-status-failed [get]
func (h *transactionStatsHandlerApi) GetMonthlyAmountFailedByMerchant(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:month-status-failed:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetMonthlyAmountFailedByMerchant(ctx, &pb_transaction.MonthAmountTransactionMerchantRequest{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthFailed(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly failed transactions by merchant
// @Tags Transaction Stats
// @Description Yearly amount of failed transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/year-status-failed [get]
func (h *transactionStatsHandlerApi) GetYearlyAmountFailedByMerchant(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:year-status-failed:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearFailed](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetYearlyAmountFailedByMerchant(ctx, &pb_transaction.YearAmountTransactionMerchantRequest{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearFailed(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly successful transactions by merchant and payment method
// @Tags Transaction Stats
// @Description Monthly payment-method breakdown of successful transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/month-method-success [get]
func (h *transactionStatsHandlerApi) GetMonthlyMethodByMerchantSuccess(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:month-method-success:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetMonthlyTransactionMethodByMerchantSuccess(ctx, &pb_transaction.MonthMethodTransactionMerchantRequest{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly successful transactions by merchant and payment method
// @Tags Transaction Stats
// @Description Yearly payment-method breakdown of successful transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/year-method-success [get]
func (h *transactionStatsHandlerApi) GetYearlyMethodByMerchantSuccess(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:year-method-success:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetYearlyTransactionMethodByMerchantSuccess(ctx, &pb_transaction.YearMethodTransactionMerchantRequest{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Monthly failed transactions by merchant and payment method
// @Tags Transaction Stats
// @Description Monthly payment-method breakdown of failed transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/month-method-failed [get]
func (h *transactionStatsHandlerApi) GetMonthlyMethodByMerchantFailed(c echo.Context) error {
	year, month, merchantID, err := parseYearMonthMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:month-method-failed:%d:%d:%d", statsCacheNamespace, year, month, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetMonthlyTransactionMethodByMerchantFailed(ctx, &pb_transaction.MonthMethodTransactionMerchantRequest{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionMonthMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Yearly failed transactions by merchant and payment method
// @Tags Transaction Stats
// @Description Yearly payment-method breakdown of failed transactions for a merchant
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/transaction-stats/by-merchant/year-method-failed [get]
func (h *transactionStatsHandlerApi) GetYearlyMethodByMerchantFailed(c echo.Context) error {
	year, merchantID, err := parseYearMerchant(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	key := fmt.Sprintf("%s:transaction-stats:by-merchant:year-method-failed:%d:%d", statsCacheNamespace, year, merchantID)

	if cached, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, h.cache, key); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.clients.TransactionStatsByMerchant.GetYearlyTransactionMethodByMerchantFailed(ctx, &pb_transaction.YearMethodTransactionMerchantRequest{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransactionYearMethod(res)
	cache.SetToCache(ctx, h.cache, key, apiResponse, statsCacheTTL)

	return c.JSON(http.StatusOK, apiResponse)
}
