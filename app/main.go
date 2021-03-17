package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/fvbock/endless"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

//go:embed build
var embeddedFiles embed.FS

func main() {
	fsys, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Fatalf("Failed to subpath embedded files: %s", err)
	}

	e := echo.New()
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(fsys)))
	e.GET("/*", contentHandler, rewritePath)

	endless.DefaultHammerTime = 10 * time.Second
	endless.DefaultReadTimeOut = 295 * time.Second
	if err := endless.ListenAndServe(":5000", e); err != nil {
		log.Infof("Server stopped: %s", err)
	}
}

// Required to enable single page app routing.
// Anything that is not a static file is redirected to the index.html
// so it can be routed by the react-router
func rewritePath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Request().URL.Path, "static") {
			c.Request().URL.Path = "/"
		}
		return next(c)
	}
}
