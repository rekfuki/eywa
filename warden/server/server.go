package server

import (
	"io"
	"text/template"
	"time"

	"github.com/fvbock/endless"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

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
}

func contextObjects(contextParams *ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", contextParams.DB)
			c.Set("session_signing_key", contextParams.SessionSigningKey)
			c.Set("session_timeout_duration", contextParams.SessionTimeoutDuration)
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

	// Expose metrics for prometheus
	// e.GET("/metrics", echo.WrapHandler(params.Metrics.PrometheusHandler()))

	t := &Template{
		templates: template.Must(template.New("").ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.GET("/login", controllers.ShowLogin)
	e.GET("/auth", controllers.OAuth)
	e.GET("/auth/callback", controllers.CompleteOAuth)

	return e
}
