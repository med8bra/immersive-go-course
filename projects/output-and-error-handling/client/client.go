package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type WeatherClient interface {
	GetWeather() (string, error)
}

type weatherClient struct {
	httpClient     *http.Client
	weatherBaseUrl string
	retryAfter     time.Time
	m              sync.RWMutex
}

var ErrServerBusy = errors.New("server is too busy")
var ErrServerNotReachable = errors.New("server is not reachable")

// GetWeather implements WeatherClient.
func (wc *weatherClient) GetWeather() (string, error) {
	d := wc.getRetryAfter()

	if d != 0 {
		return "", fmt.Errorf("Please retry after %s %w", d, ErrServerBusy)
	}

	r, err := wc.httpClient.Get(wc.weatherBaseUrl)
	if err != nil {
		return "", ErrServerNotReachable
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := time.Parse(http.TimeFormat, r.Header.Get("Retry-After"))
		if err != nil {
			return "", err
		}
		wc.setRetryAfter(retryAfter)
		return "", ErrServerBusy
	}
	if r.StatusCode != http.StatusOK {
		return "", ErrServerNotReachable
	}

	body, err := io.ReadAll(r.Body)

	return string(body), nil
}

func (wc *weatherClient) setRetryAfter(retryAfter time.Time) {
	wc.m.Lock()
	defer wc.m.Unlock()
	wc.retryAfter = retryAfter
}

func (wc *weatherClient) getRetryAfter() time.Duration {
	wc.m.RLock()
	defer wc.m.RUnlock()
	now := time.Now()

	d := wc.retryAfter.Sub(now)
	if d > 0 {
		return d
	}
	return 0
}

func NewWeatherClient(client *http.Client, weatherBaseUrl string) WeatherClient {
	return &weatherClient{
		httpClient:     client,
		weatherBaseUrl: weatherBaseUrl,
		retryAfter:     time.Now(),
	}
}
