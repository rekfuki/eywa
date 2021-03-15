package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// FunctionRunner runs a function
type FunctionRunner interface {
	Run(f FunctionRequest) error
}

// FunctionRequest stores request for function execution
type FunctionRequest struct {
	Process     string
	ProcessArgs []string
	Environment []string

	InputReader   io.ReadCloser
	OutputWriter  io.Writer
	ContentLength *int64
}

// HTTPFunctionRunner creates and maintains one process responsible for handling all calls
type HTTPFunctionRunner struct {
	ExecTimeout    time.Duration // ExecTimeout the maximum duration or an upstream function call
	ReadTimeout    time.Duration // ReadTimeout for HTTP server
	WriteTimeout   time.Duration // WriteTimeout for HTTP Server
	Process        string        // Process to run as fprocess
	ProcessArgs    []string      // ProcessArgs to pass to command
	Command        *exec.Cmd
	Stdout         *[]string
	Stderr         *[]string
	Client         *http.Client
	UpstreamURL    *url.URL
	BufferHTTPBody bool
	WriteDebug     bool
}

type FunctionResponse struct {
	Body    []byte      `json:"body,omitempty"`
	Headers http.Header `json:"headers,omitempty"`
	Status  int         `json:"status"`
	Stdout  []string    `json:"stdout,omitempty"`
	Stderr  []string    `json:"stderr,omitempty"`
}

// Start forks the process used for processing incoming requests
func (f *HTTPFunctionRunner) Start() error {
	cmd := exec.Command(f.Process, f.ProcessArgs...)

	f.Stdout = &[]string{}
	f.Stderr = &[]string{}

	f.Command = cmd
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Logs lines from stderr and stdout to the output buffer
	bindLoggingPipe("stderr", stdoutPipe, f.Stdout)
	bindLoggingPipe("stdout", stderrPipe, f.Stderr)

	f.Client = makeProxyClient(f.ExecTimeout)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Fatalf("Failed to process SIGTERM: %s", err)
		}

	}()

	err = cmd.Start()
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Fatalf("Forked function has terminated: %s", err.Error())
		}
	}()

	return err
}

// Run a function with a long-running process with a HTTP protocol for communication
func (f *HTTPFunctionRunner) Run(req FunctionRequest, contentLength int64, r *http.Request, w http.ResponseWriter) error {
	startedTime := time.Now()

	// Wipe any previous output
	*f.Stderr = []string{}
	*f.Stdout = []string{}

	upstreamURL := f.UpstreamURL.String()

	if len(r.RequestURI) > 0 {
		upstreamURL += r.RequestURI
	}

	var body io.Reader
	if f.BufferHTTPBody {
		reqBody, _ := ioutil.ReadAll(r.Body)
		body = bytes.NewReader(reqBody)
	} else {
		body = r.Body
	}

	request, _ := http.NewRequest(r.Method, upstreamURL, body)
	request.Host = r.Host

	copyHeaders(request.Header, &r.Header)

	var reqCtx context.Context
	var cancel context.CancelFunc

	if f.ExecTimeout.Nanoseconds() > 0 {
		reqCtx, cancel = context.WithTimeout(r.Context(), f.ExecTimeout)
	} else {
		reqCtx = r.Context()
		cancel = func() {

		}
	}

	defer cancel()

	res, err := f.Client.Do(request.WithContext(reqCtx))
	if err != nil {
		log.Printf("Upstream HTTP request error: %s\n", err.Error())

		// Error unrelated to context / deadline
		if reqCtx.Err() == nil {
			w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%f", time.Since(startedTime).Seconds()))

			w.WriteHeader(http.StatusInternalServerError)

			return nil
		}

		<-reqCtx.Done()
		if reqCtx.Err() != nil {
			// Error due to timeout / deadline
			log.Printf("Upstream HTTP killed due to exec_timeout: %s\n", f.ExecTimeout)
			w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%f", time.Since(startedTime).Seconds()))

			w.WriteHeader(http.StatusGatewayTimeout)
			return nil
		}

		w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%f", time.Since(startedTime).Seconds()))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%f", time.Since(startedTime).Seconds()))
	w.Header().Set("Content-Type", "application/json")

	resp := FunctionResponse{
		Status:  res.StatusCode,
		Headers: http.Header{},
	}

	copyHeaders(resp.Headers, &res.Header)

	if f.WriteDebug {
		resp.Stdout = *f.Stdout
		resp.Stderr = *f.Stderr
	}

	if res.Body != nil {
		defer res.Body.Close()

		bodyBytes, bodyErr := ioutil.ReadAll(res.Body)
		if bodyErr != nil {
			log.Printf("Failed to read body: %s", bodyErr)
		}
		resp.Body = bodyBytes
	}

	if !resp.isEmpty() {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to write response")
		}
	}

	return nil
}

func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination)[k] = vClone
	}
}

func makeProxyClient(dialTimeout time.Duration) *http.Client {
	proxyClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			DisableKeepAlives:     false,
			IdleConnTimeout:       500 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &proxyClient
}

func (fr *FunctionResponse) isEmpty() bool {
	if len(fr.Body) > 0 || len(fr.Stderr) > 0 || len(fr.Stdout) > 0 || len(fr.Headers) > 0 {
		return false
	}
	return true
}
