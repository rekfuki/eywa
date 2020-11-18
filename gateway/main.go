package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/cache"
	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
)

func runFunction(c echo.Context) error {
	cache := c.Get("cache").(*cache.Client)
	k8s := c.Get("k8s").(*k8s.Client)
	fnName := c.Param("name")

	log.Infof("Recevied request to run function: %s", fnName)

	fn := cache.LookupFunction(fnName)
	if fn == nil {
		log.Infof("Function not found in the cache, assuming not running...")

		newFn := &types.Function{
			Name:      "test-k8s-deployment",
			Image:     "registry.eywa.rekfuki.dev/go-watchdog:0.0.1",
			Namespace: "faas",
		}

		k8s.CreateDeployment(newFn)
	}

	return c.JSON(http.StatusOK, "ok")
}

// ContextParams holds the objects required to initialise the server.
type ContextParams struct {
	Cache *cache.Client
	K8s   *k8s.Client
}

// ContextObjects attaches backend clients to the API context
func ContextObjects(contextParams ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("cache", contextParams.Cache)
			c.Set("k8s", contextParams.K8s)
			return next(c)
		}
	}
}

func main() {
	k8s, err := k8s.Setup(&k8s.Config{})
	if err != nil {
		log.Fatalf("Failed to setup k8s client: %s", err)
	}

	cache := cache.Setup(&cache.Config{})

	params := ContextParams{
		Cache: cache,
		K8s:   k8s,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ContextObjects(params))

	e.POST("/system/function/:name", runFunction)

	e.Logger.Fatal(e.Start(":8080"))
}
