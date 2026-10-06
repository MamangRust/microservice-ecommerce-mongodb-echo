package role_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apicache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache"
	rolehandler "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/handler/role"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	app_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	tests "github.com/MamangRust/microservice-ecommerce-test"

	apigatewaymiddlewares "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type RoleApiTestSuite struct {
	tests.BaseTestSuite
	echo   *echo.Echo
	roleID int
}

func (s *RoleApiTestSuite) SetupSuite() {
	s.PrimaryContext = database.IdentityDB
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()

	// Setup Echo and API Handler
	e := echo.New()
	apigatewaymiddlewares.RegisterErrorHandler(e)
	s.echo = e

	// Bypass auth middleware by setting user_id and seeding roles in Redis
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", 1)
			return next(c)
		}
	})

	roles := []string{"Admin_Role_10", "Admin_Admin_14"}
	// F6-2: the gateway namespaces its cache keys with apigw: — the role
	// middleware reads cacheRoleKey "apigw:user_roles:%s", so the seeded cache
	// entry must use the same prefixed key.
	cacheStore := s.GetCacheStore()
	cache.SetToCache(s.Ctx, cacheStore, "apigw:user_roles:1", &roles, 5*time.Minute)

	apiErrorHandler := app_errors.NewApiHandler(s.Obs, s.Log)
	rolehandler.RegisterRoleHandler(&rolehandler.DepsRole{
		Client:     s.Conns["role"],
		Kafka:      nil,
		E:          e,
		Logger:     s.Log,
		CacheStore: cacheStore,
		Cache:      apicache.NewRoleCache(cacheStore),
		ApiHandler: apiErrorHandler,
	})
}

func (s *RoleApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateRoleRequest{
		Name: "API Role",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role-command/create", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.Name, data["name"])
	s.roleID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/role-query", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.roleID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/role-query/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.roleID), data["id"])

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/active", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. FindByUserId
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/user/1", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Update
	updateBody := requests.UpdateRoleRequest{
		Name: "Updated API Role",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/update/%d", s.roleID), bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/trashed/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. FindByTrashed (verify trashed)
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/restore/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. Trash again for DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/trashed/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 12. DeletePermanent
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/role-command/permanent/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 13. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 14. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
