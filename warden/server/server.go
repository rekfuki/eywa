package server

import (
	"io"
	"net/http"
	"text/template"
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
	"eywa/warden/clients/tugrik"
	"eywa/warden/controllers"
	"eywa/warden/db"
)

// Template stores templates to be rendered
type Template struct {
	templates *template.Template
}

// Render renders templates
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// ContextParams holds the objects required to initialise the server.
type ContextParams struct {
	DB                     *db.Client
	SessionSigningKey      string
	SessionTimeoutDuration time.Duration
	Tugrik                 *tugrik.Client
}

func contextObjects(contextParams *ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", contextParams.DB)
			c.Set("session_signing_key", contextParams.SessionSigningKey)
			c.Set("session_timeout_duration", contextParams.SessionTimeoutDuration)
			c.Set("tugrik", contextParams.Tugrik)
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
	if err := endless.ListenAndServe(":1080", r); err != nil {
		log.Infof("Server stopped: %s", err)
	}
}

func createRouter(params *ContextParams) *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(contextObjects(params))

	t := &Template{
		templates: template.Must(template.New("").ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.GET("/login", controllers.ShowLogin)
	e.POST("/logout", controllers.Logout)
	e.GET("/oauth", controllers.OAuth)
	e.GET("/authn", controllers.CheckAuth)
	e.GET("/users/me", controllers.GetUser)

	wardenAPI := createWardenAPI()
	e.GET("/eywa/api/warden/doc", echo.WrapHandler(wardenAPI.Handler(true)))

	api := e.Group("", checkAuth(), sv.SwaggerValidatorEcho(wardenAPI), pagination.Validate())
	wardenAPI.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c echo.Context) error)
		path = swag.ColonPath(path)
		api.Add(endpoint.Method, path, h)
	})

	return e
}

func createWardenAPI() *swagger.API {
	return swag.New(
		swag.Title("Warden"),
		swag.Description(`Warden API`),
		swag.Version("2.0"),
		swag.BasePath("/eywa/api"),
		swag.Endpoints(aggregateEndpoints(
			accessTokensAPI(),
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
