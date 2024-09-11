package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/carlmjohnson/be"
)

// TestGetHealth tests the /health endpoint.
// Server is started by [TestMain] so that the test can make requests to it.
func TestGetHealth(t *testing.T) {
	// response is repeated but this describes intention of test better.
	// For example you can add fiels only needed for testing.
	type response struct {
		Version  string    `json:"version"`
		Revision string    `json:"vcs.revision"`
		Time     time.Time `json:"vcs.time"`
		// Modified bool      `json:"vcs.modified"`
	}

	// actual http request to the server.
	res, err := http.Get(endpoint() + "/health")
	be.NilErr(t, err)
	be.Equal(t, http.StatusOK, res.StatusCode)
	be.Equal(t, "application/json", res.Header.Get("Content-Type"))
	be.NilErr(t, json.NewDecoder(res.Body).Decode(&response{}))
	defer res.Body.Close()
}

// TestGetOpenapi tests the /openapi.yaml endpoint.
// You can add more test as needed without starting the server again.
func TestGetOpenapi(t *testing.T) {
	res, err := http.Get(endpoint() + "/openapi.yaml")
	be.NilErr(t, err)
	be.Equal(t, http.StatusOK, res.StatusCode)
	be.Equal(t, "text/plain", res.Header.Get("Content-Type"))

	sb := strings.Builder{}
	_, err = io.Copy(&sb, res.Body)
	be.NilErr(t, err)
	res.Body.Close()
	be.In(t, "openapi: 3.0.0", sb.String())
}

// TestMain starts the server and runs all the tests.
// in case of short test, (arg -short given by command line) the server is not started.
// By doing this, you can run **actual** integration tests without starting the server.
func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()
	short, _ := strconv.ParseBool(flag.CommandLine.Lookup("test.short").Value.String())
	if !short {
		args := []string{"testapp", "--port", port()}
		go func() {
			err := run(ctx, os.Stdout, args)
			if err != nil {
				log.Fatalf("failed to run with args %v got err %s\n", args, err)
			}
		}()
		waitForHealthy(ctx, 2*time.Second, endpoint()+"/health")
	}

	os.Exit(m.Run())
}

// endpoint returns the server endpoint started by [TestMain].
func endpoint() string {
	return "http://localhost:" + port()
}

// port returns the port on which the server is started by [TestMain].
// on the first call, it starts a listener on a random port and returns the port.
func port() string {
	_portOnce.Do(func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer listener.Close()
		addr := listener.Addr().(*net.TCPAddr)
		_port = strconv.Itoa(addr.Port)
	})
	return _port
}

// this variables is not intended to be used directly. use [port] function instead.
var (
	_portOnce sync.Once
	_port     string
)

// waitForHealthy waits for the server to be healthy.
// this function is used by [TestMain] to wait for the server to be healthy before running tests.
func waitForHealthy(ctx context.Context, timeout time.Duration, endpoint string) {
	startTime := time.Now()
	for {
		res, err := http.Get(endpoint)
		if err == nil && res.StatusCode == http.StatusOK {
			log.Printf("endpoint %s is ready after %v\n", endpoint, time.Since(startTime))
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
			if timeout <= time.Since(startTime) {
				log.Fatalf("timeout %v reached while waitForHealthy", timeout)
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}
