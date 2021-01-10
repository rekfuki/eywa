package server

import (
	"net/http"
	"time"

	"github.com/fvbock/endless"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/miketonks/swag"
	sv "github.com/miketonks/swag-validator"
	"github.com/miketonks/swag/swagger"
	log "github.com/sirupsen/logrus"

	"eywa/go-libs/auth"
	"eywa/go-libs/pagination"
	"eywa/registry/builder"
	"eywa/registry/clients/docker"
	"eywa/registry/clients/mongo"
)

// ContextParams holds the objects required to initialise the server.
type ContextParams struct {
	Mongo   *mongo.Client
	Builder *builder.Client
	Docker  *docker.Client
	// Gateway *gateway.Client
}

func contextObjects(contextParams *ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mongo", contextParams.Mongo)
			c.Set("builder", contextParams.Builder)
			c.Set("docker", contextParams.Docker)
			// c.Set("gateway", contextParams.Gateway)
			return next(c)
		}
	}
}

func checkAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := auth.FromHeaders(c.Request().Header)
			if !auth.Check() {
				return c.JSON(http.StatusForbidden, "Forbidden")
			}
			c.Set("auth", auth)
			return next(c)
		}
	}
}

// Run starts the api server.
func Run(params *ContextParams) {
	r := createRouter(params)

	endless.DefaultHammerTime = 10 * time.Second
	endless.DefaultReadTimeOut = 295 * time.Second
	if err := endless.ListenAndServe(":9080", r); err != nil {
		log.Infof("Server stopped: %s", err)
	}
}

func createRouter(params *ContextParams) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(contextObjects(params))

	enableCors := true
	registryAPI := createRegistryAPI()
	e.GET("/eywa/api/registry/doc", echo.WrapHandler(registryAPI.Handler(enableCors)))

	api := e.Group("", checkAuth(), sv.SwaggerValidatorEcho(registryAPI), pagination.Validate())
	registryAPI.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c echo.Context) error)
		path = swag.ColonPath(path)
		api.Add(endpoint.Method, path, h)
	})

	return e
}

func createRegistryAPI() *swagger.API {
	return swag.New(
		swag.Title("Registry"),
		swag.Description(`Registry API`),
		swag.Version("2.0"),
		swag.BasePath("/eywa/api"),
		swag.Endpoints(aggregateEndpoints(
			imagesAPI(),
		)...,
		),
	)
}

func aggregateEndpoints(endpoints ...[]*swagger.Endpoint) (res []*swagger.Endpoint) {
	for _, v := range endpoints {
		res = append(res, v...)
	}

	return
}
