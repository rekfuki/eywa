package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"eywa/watchdog/config"
	"eywa/watchdog/executor"
	limiter "eywa/watchdog/limiter"
	"eywa/watchdog/metrics"
)

var (
	acceptingConnections int32
)

func main() {
	var runHealthcheck bool

	flag.BoolVar(&runHealthcheck,
		"run-healthcheck",
		false,
		"Check for the a lock-file, when using an exec healthcheck. Exit 0 for present, non-zero when not found.")

	flag.Parse()

	if runHealthcheck {
		if lockFilePresent() {
			os.Exit(0)
		}

		log.Fatal("Unable to find lock file")
	}

	atomic.StoreInt32(&acceptingConnections, 0)

	wc := config.New(os.Environ())
	fmt.Printf("post config parse")

	if len(wc.FunctionProcess) == 0 && wc.OperationalMode != config.ModeStatic {
		log.Fatal(`Provide a "function_process" or "fprocess" environmental variable for your function"`)
	}

	if wc.Debug == true {
		log.Infof("Debug mode enabled")
	}

	requestHandler := buildRequestHandler(wc)

	log.Infof("OperationalMode: %s\n", config.WatchdogMode(wc.OperationalMode))

	httpMetrics := metrics.NewHttp()
	http.HandleFunc("/", metrics.InstrumentHandler(requestHandler, httpMetrics))
	http.HandleFunc("/_/health", makeHealthHandler())

	metricsServer := metrics.MetricsServer{}
	metricsServer.Register(wc.MetricsPort)

	cancel := make(chan bool)

	go metricsServer.Serve(cancel)

	shutdownTimeout := wc.HTTPWriteTimeout
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", wc.TCPPort),
		ReadTimeout:    wc.HTTPReadTimeout,
		WriteTimeout:   wc.HTTPWriteTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	log.Printf("Timeouts: read: %s, write: %s hard: %s.\n",
		wc.HTTPReadTimeout,
		wc.HTTPWriteTimeout,
		wc.ExecTimeout)
	log.Printf("Listening on port: %d\n", wc.TCPPort)

	listenUntilShutdown(shutdownTimeout, s, wc.SuppressLock)
}

func markUnhealthy() error {
	atomic.StoreInt32(&acceptingConnections, 0)

	path := filepath.Join(os.TempDir(), ".lock")
	log.Printf("Removing lock-file : %s\n", path)
	removeErr := os.Remove(path)
	return removeErr
}

func listenUntilShutdown(shutdownTimeout time.Duration, s *http.Server, suppressLock bool) {

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig

		log.Printf("SIGTERM received.. shutting down server in %s\n", shutdownTimeout.String())

		healthErr := markUnhealthy()

		if healthErr != nil {
			log.Printf("Unable to mark unhealthy during shutdown: %s\n", healthErr.Error())
		}

		<-time.Tick(shutdownTimeout)

		if err := s.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Error in Shutdown: %v", err)
		}

		log.Printf("No new connections allowed. Exiting in: %s\n", shutdownTimeout.String())

		<-time.Tick(shutdownTimeout)

		close(idleConnsClosed)
	}()

	// Run the HTTP server in a separate go-routine.
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Error ListenAndServe: %v", err)
			close(idleConnsClosed)
		}
	}()

	if !suppressLock {
		path, writeErr := createLockFile()

		if writeErr != nil {
			log.Panicf("Cannot write %s. To disable lock-file set env suppress_lock=true.\n Error: %s.\n", path, writeErr.Error())
		}
	} else {
		log.Println("Warning: \"suppress_lock\" is enabled. No automated health-checks will be in place for your function.")

		atomic.StoreInt32(&acceptingConnections, 1)
	}

	<-idleConnsClosed
}

func buildRequestHandler(wc config.WatchdogConfig) http.Handler {
	var requestHandler http.HandlerFunc

	switch wc.OperationalMode {
	case config.ModeHTTP:
		requestHandler = makeHTTPRequestHandler(wc)
	default:
		log.Panicf("unknown watchdog mode: %d", wc.OperationalMode)
	}

	if wc.MaxInflight > 0 {
		return limiter.NewConcurrencyLimiter(requestHandler, wc.MaxInflight)
	}

	return requestHandler
}

// createLockFile returns a path to a lock file and/or an error
// if the file could not be created.
func createLockFile() (string, error) {
	path := filepath.Join(os.TempDir(), ".lock")
	log.Printf("Writing lock-file to: %s\n", path)

	mkdirErr := os.MkdirAll(os.TempDir(), os.ModePerm)
	if mkdirErr != nil {
		return path, mkdirErr
	}

	writeErr := ioutil.WriteFile(path, []byte{}, 0660)
	if writeErr != nil {
		return path, writeErr
	}

	atomic.StoreInt32(&acceptingConnections, 1)
	return path, nil
}

func makeHTTPRequestHandler(wc config.WatchdogConfig) func(http.ResponseWriter, *http.Request) {
	commandName, arguments := wc.Process()
	functionInvoker := executor.HTTPFunctionRunner{
		ExecTimeout:    wc.ExecTimeout,
		Process:        commandName,
		ProcessArgs:    arguments,
		BufferHTTPBody: wc.BufferHTTPBody,
		WriteDebug:     wc.Debug,
	}

	if len(wc.UpstreamURL) == 0 {
		log.Fatal(`For "mode=http" you must specify a valid URL for "http_upstream_url"`)
	}

	urlValue, upstreamURLErr := url.Parse(wc.UpstreamURL)
	if upstreamURLErr != nil {
		log.Fatal(upstreamURLErr)
	}
	functionInvoker.UpstreamURL = urlValue

	fmt.Printf("Forking - %s %s\n", commandName, arguments)
	err := functionInvoker.Start()
	if err != nil {
		log.Fatalf("Failed to start function invoker: %s", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := executor.FunctionRequest{
			Process:      commandName,
			ProcessArgs:  arguments,
			InputReader:  r.Body,
			OutputWriter: w,
		}

		if r.Body != nil {
			defer r.Body.Close()
		}

		err := functionInvoker.Run(req, r.ContentLength, r, w)
		if err != nil {
			log.Errorf("Failed to run function invocation: %s", err)
			w.WriteHeader(500)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				log.Printf("Failed to wrtie error response: %s", err)
			}
		}
	}
}

func lockFilePresent() bool {
	path := filepath.Join(os.TempDir(), ".lock")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func makeHealthHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if atomic.LoadInt32(&acceptingConnections) == 0 || !lockFilePresent() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				log.Printf("Failed to write to health handler response: %s", err)
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
