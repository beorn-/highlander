package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HighlanderProxy struct {
	allowed  string
	isLeader func() bool
	weight   uint64
	lastCall time.Time
}

func NewHighlanderProxy(checkInterval, expirationInterval time.Duration, isLeader func() bool) *HighlanderProxy {
	f := &HighlanderProxy{isLeader: isLeader}

	go func() {
		t := time.NewTicker(checkInterval)
		for {
			select {
			case <-t.C:
				if f.lastCall.Before(time.Now().Add(-expirationInterval)) {
					log.Printf("lost source : '%s' (no data for %s)\n", f.allowed, expirationInterval.String())
					f.allowed = ""
					f.weight = 0
				}
			}
		}
	}()

	return f
}

func (f *HighlanderProxy) RoundTrip(r *http.Request) (*http.Response, error) {
	caller := r.RemoteAddr
	var weight uint64
	var err error

	keys, ok := r.URL.Query()["highlander_weight"]
	if !ok {
		weight = 0
	} else {
		weight, err = strconv.ParseUint(keys[0], 10, 64)
		if err != nil {
			weight = 0
		}
	}

	if weight > f.weight {
		log.Printf("new source : '%s' (bigger Highlander Weight) (%d -> %d)\n", caller, f.weight, weight)
		f.allowed = caller
		f.weight = weight
	}

	if f.allowed == "" {
		log.Printf("new source : '%s' (no current promoted source) (%d)\n", caller, weight)
		f.allowed = caller
		f.weight = weight
	}

	if caller != f.allowed || !f.isLeader() {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("ok")),
			Request:    r,
		}, nil
	}

	f.lastCall = time.Now()

	r.Header.Add("X-Highlander-Instance", cmd.instanceID)

	return http.DefaultTransport.RoundTrip(r)
}
