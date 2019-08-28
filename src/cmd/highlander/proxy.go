package main

import (
	"io/ioutil"
	"log"
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
	caller := r.RemoteAddr

	if caller == cmd.preferredIp {
		f.allowed = caller
		log.Println("new source : '" + caller + "' (preferred ip)")
	}

	if f.allowed == "" {
		f.allowed = caller
		log.Println("new source : '" + caller + "' (no current promoted source")
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
