package statshandler

import (
	"strconv"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
)

// statsCacheTTL is how long stats responses are cached at the gateway.
const statsCacheTTL = 5 * time.Minute

// statsCacheNamespace prefixes every stats cache key so they can be invalidated
// independently of the per-domain gateway caches.
const statsCacheNamespace = "apigw"

func parseYear(c echo.Context) (int, error) {
	year, err := strconv.Atoi(c.QueryParam("year"))
	if err != nil || year <= 0 {
		return 0, errors.NewBadRequestError("year is required and must be a positive integer")
	}
	return year, nil
}

func parseYearMonth(c echo.Context) (int, int, error) {
	year, err := parseYear(c)
	if err != nil {
		return 0, 0, err
	}

	month, err := strconv.Atoi(c.QueryParam("month"))
	if err != nil || month < 1 || month > 12 {
		return 0, 0, errors.NewBadRequestError("month is required and must be between 1 and 12")
	}

	return year, month, nil
}

func parseYearMerchant(c echo.Context) (int, int, error) {
	year, err := parseYear(c)
	if err != nil {
		return 0, 0, err
	}

	merchantID, err := parseID(c, "merchant_id")
	if err != nil {
		return 0, 0, err
	}

	return year, merchantID, nil
}

func parseYearMonthMerchant(c echo.Context) (int, int, int, error) {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return 0, 0, 0, err
	}

	merchantID, err := parseID(c, "merchant_id")
	if err != nil {
		return 0, 0, 0, err
	}

	return year, month, merchantID, nil
}

func parseYearCategory(c echo.Context) (int, int, error) {
	year, err := parseYear(c)
	if err != nil {
		return 0, 0, err
	}

	categoryID, err := parseID(c, "category_id")
	if err != nil {
		return 0, 0, err
	}

	return year, categoryID, nil
}

func parseYearMonthCategory(c echo.Context) (int, int, int, error) {
	year, month, err := parseYearMonth(c)
	if err != nil {
		return 0, 0, 0, err
	}

	categoryID, err := parseID(c, "category_id")
	if err != nil {
		return 0, 0, 0, err
	}

	return year, month, categoryID, nil
}

func parseID(c echo.Context, name string) (int, error) {
	id, err := strconv.Atoi(c.QueryParam(name))
	if err != nil || id <= 0 {
		return 0, errors.NewBadRequestError(name + " is required and must be a positive integer")
	}
	return id, nil
}
