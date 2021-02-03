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

	"eywa/execution-tracker/db"
	"eywa/go-libs/auth"
	"eywa/go-libs/pagination"
)

// Config is a populated by env variables and Vault
type Config struct {
	DefaultTimeRange string `default:"1h"`
	MaxTimeRange     string `default:"10080m"`
	Postgres         db.Config
	Component        string
	Container        string
}

// ContextParams holds the objects required to initialise the server
type ContextParams struct {
	DB               *db.Client
	DefaultTimeRange time.Duration
	MaxTimeRange     time.Duration
	MaxCount         int
}

func contextObjects(contextParams *ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", contextParams.DB)
			c.Set("defaultTimeRange", contextParams.DefaultTimeRange)
			c.Set("maxTimeRange", contextParams.MaxTimeRange)
			c.Set("maxCount", contextParams.MaxCount)
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

func createRouter(params *ContextParams) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(contextObjects(params))

	lumberjackAPI := createLogsSwaggerAPI()
	api := e.Group("", checkAuth(), sv.SwaggerValidatorEcho(lumberjackAPI), pagination.Validate())
	lumberjackAPI.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c echo.Context) error)
		path = swag.ColonPath(path)
		api.Add(endpoint.Method, path, h)
	})

	return e
}

// Run starts the gin Router and listens forever, recovering from panics
func Run(conf Config, dbc *db.Client) {
	defaultTimeRange, err := time.ParseDuration(conf.DefaultTimeRange)
	if err != nil {
		log.Fatalf("Invalid default time range in config: %s", err)
	}

	maxTimeRange, err := time.ParseDuration(conf.MaxTimeRange)
	if err != nil {
		log.Fatalf("Invalid maximum time range in config: %s", err)
	}

	params := &ContextParams{
		DB:               dbc,
		DefaultTimeRange: defaultTimeRange,
		MaxTimeRange:     maxTimeRange,
	}

	r := createRouter(params)

	endless.DefaultHammerTime = 10 * time.Second
	endless.DefaultReadTimeOut = 295 * time.Second
	if err := endless.ListenAndServe(":10080", r); err != nil {
		log.Infof("Server stopped: %s", err)
	}
}

func createLogsSwaggerAPI() *swagger.API {

	api := swag.New(
		swag.Title("Execution-tracker API"),
		swag.Description("Provides access to execution logs"),
		swag.Version("2.0"),
		swag.BasePath("/eywa/api"),
		swag.Endpoints(aggregateEndpoints(
			logsAPI(),
		)...,
		),
	)
	return api
}

func aggregateEndpoints(endpoints ...[]*swagger.Endpoint) (res []*swagger.Endpoint) {
	for _, v := range endpoints {
		res = append(res, v...)
	}

	return
}
