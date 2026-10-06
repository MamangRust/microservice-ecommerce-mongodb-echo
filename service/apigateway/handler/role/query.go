package rolehandler

import (
	"net/http"
	"strconv"

	mencache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/role"
	"github.com/MamangRust/microservice-ecommerce-grpc-apigateway/middlewares"
	"github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pb_user_role "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/role"
	"github.com/labstack/echo/v4"
)

type roleQueryHandlerApi struct {
	kafka      *kafka.Kafka
	role       pb_role.RoleQueryServiceClient
	userRole   pb_user_role.UserRoleServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.RoleQueryResponseMapper
	cache      role_cache.RoleQueryCache
	apiHandler errors.ApiHandler
}

type roleQueryHandleDeps struct {
	client         pb_role.RoleQueryServiceClient
	userRoleClient pb_user_role.UserRoleServiceClient
	router         *echo.Echo
	logger         logger.LoggerInterface
	mapper         apimapper.RoleQueryResponseMapper
	kafka          *kafka.Kafka
	cache_role     mencache.RoleCache
	cache          role_cache.RoleQueryCache
	apiHandler     errors.ApiHandler
}

func NewRoleQueryHandleApi(params *roleQueryHandleDeps) *roleQueryHandlerApi {
	handler := &roleQueryHandlerApi{
		role:       params.client,
		userRole:   params.userRoleClient,
		logger:     params.logger,
		mapper:     params.mapper,
		kafka:      params.kafka,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	roleMiddleware := middlewares.RoleValidatorGRPC(params.userRoleClient, params.logger, params.cache_role)
	routerRole := params.router.Group("/api/role-query")
	requireAdmin := middlewares.RequireRoles("Admin", "ROLE_ADMIN", "Admin_Role_10")

	routerRole.GET("", roleMiddleware(requireAdmin(handler.FindAll)))
	routerRole.GET("/:id", roleMiddleware(requireAdmin(handler.FindById)))
	routerRole.GET("/active", roleMiddleware(requireAdmin(handler.FindByActive)))
	routerRole.GET("/trashed", roleMiddleware(requireAdmin(handler.FindByTrashed)))
	routerRole.GET("/user/:user_id", roleMiddleware(requireAdmin(handler.FindByUserId)))

	return handler
}

// @Security Bearer
// @Summary Find all roles
// @Tags Role Query
// @Description Retrieve a list of all roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query [get]
func (h *roleQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoles(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindAllRole(ctx, &pb_role.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRole(res)
	h.cache.SetCachedRoles(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find role by ID
// @Tags Role Query
// @Description Retrieve a role by ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/{id} [get]
func (h *roleQueryHandlerApi) FindById(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedRoleById(ctx, roleID); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByIdRole(ctx, &pb_role.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseRole(res)
	h.cache.SetCachedRoleById(ctx, roleID, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active roles
// @Tags Role Query
// @Description Retrieve a list of active roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/active [get]
func (h *roleQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByActive(ctx, &pb_role.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed roles
// @Tags Role Query
// @Description Retrieve a list of trashed role records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed role data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/trashed [get]
func (h *roleQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByTrashed(ctx, &pb_role.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find roles by user ID
// @Tags Role Query
// @Description Retrieve a list of roles assigned to a specific user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponsesRole "List of user roles"
// @Failure 400 {object} errors.ErrorResponse "Invalid user ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/user/{user_id} [get]
func (h *roleQueryHandlerApi) FindByUserId(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		return errors.NewBadRequestError("user_id is required")
	}

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedRoleByUserId(ctx, userID); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.userRole.FindByUserId(ctx, &pb_user_role.FindByIdUserRoleRequest{UserId: int32(userID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsesRole(res)
	h.cache.SetCachedRoleByUserId(ctx, userID, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}
