package client_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/med8bra/immersive-go-course/projects/output-and-error-handling/client"
)

type testCase struct {
	desc             string
	serverResponse   string
	serverHttpStatus int
	expectedWeather  string
	expectedError    error
}

var test_cases = []testCase{
	{desc: "Ping", serverResponse: "pong", serverHttpStatus: 200, expectedWeather: "pong", expectedError: nil},
	{desc: "Ping Weather", serverResponse: "pong weather", serverHttpStatus: 200, expectedWeather: "pong weather", expectedError: nil},
	{desc: "Server error", serverResponse: "", serverHttpStatus: 500, expectedWeather: "", expectedError: client.ErrServerNotReachable},
}

func TestClient(t *testing.T) {
	var active_test_case *testCase
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.RequestURI != "/" {
			t.Errorf("Expected GET / request, got %s %s", r.Method, r.RequestURI)
		}
		if active_test_case.serverHttpStatus == 0 {
			t.Errorf("Expected server to respond with status code")
		}
		w.WriteHeader(active_test_case.serverHttpStatus)
		w.Write([]byte(active_test_case.serverResponse))
	}))

	defer testServer.Close()

	weatherClient := client.NewWeatherClient(testServer.Client(), testServer.URL)

	for _, test_case := range test_cases {
		t.Run(test_case.desc, func(t *testing.T) {
			active_test_case = &test_case

			weatherClient.GetWeather()

			res, err := weatherClient.GetWeather()
			if res != test_case.expectedWeather || err != test_case.expectedError {
				t.Errorf("Expected [res: %s, err: %s], got [res: %s,err: %s]", test_case.expectedWeather, test_case.expectedError, res, err)
			}
		})
	}
}

func TestRateLimiting(t *testing.T) {
	retryAfter := time.Now().Add(time.Second)
	fmt.Printf("[test] retryAfter: %v\n", retryAfter)
	serverCalls := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalls++

		if r.Method != "GET" || r.RequestURI != "/" {
			t.Errorf("Expected GET / request, got %s %s", r.Method, r.RequestURI)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		now := time.Now()
		if serverCalls > 1 && now.Before(retryAfter) {
			t.Errorf("Expected request after %s", retryAfter.Sub(now))
		}
		w.Header().Set("Retry-After", retryAfter.UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusTooManyRequests)

	}))

	defer testServer.Close()

	weatherClient := client.NewWeatherClient(testServer.Client(), testServer.URL)

	for i := 0; i < 5; i++ {
		if _, err := weatherClient.GetWeather(); !errors.Is(err, client.ErrServerBusy) {
			t.Errorf("Expected [err: %s], got [err: %s]", client.ErrServerBusy, err)
		}
	}

	time.Sleep(retryAfter.Sub(time.Now()))

	if _, err := weatherClient.GetWeather(); !errors.Is(err, client.ErrServerBusy) {
		t.Errorf("Expected [err: %s], got [err: %s]", client.ErrServerBusy, err)
	}

	if serverCalls != 2 {
		t.Errorf("Expected 2 server calls, got %d", serverCalls)
	}

}
