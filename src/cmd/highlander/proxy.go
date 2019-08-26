package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HighlanderProxy struct {
	allowed  string
	lastCall time.Time
}

func NewHighlanderProxy(checkInterval, expirationInterval time.Duration) *HighlanderProxy {
	f := &HighlanderProxy{}

	go func() {
		t := time.NewTicker(checkInterval)
		for {
			select {
			case <-t.C:
				if f.lastCall.Before(time.Now().Add(-expirationInterval)) {
					f.allowed = ""
				}
			}
		}
	}()

	return f
}

func (f *HighlanderProxy) RoundTrip(r *http.Request) (*http.Response, error) {
	caller := r.Header.Get("X-Caller")
	if f.allowed == "" {
		f.allowed = caller
	}

	if caller != f.allowed {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("ok")),
			Request:    r,
		}, nil
	}

	f.lastCall = time.Now()

	return http.DefaultTransport.RoundTrip(r)
}
