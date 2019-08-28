package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var cmd Cmd

func main() {
	cmd = parseCmd()

	log.Println("binding address: " + cmd.bind)
	log.Println("remote address: " + cmd.remote)
	log.Println("health tick: " + cmd.healthTick.String())
	log.Println("health expiry: " + cmd.healthExpiry.String())
	log.Println("preferred ip: " + cmd.preferredIp)

	u, err := url.Parse(cmd.remote)
	if err != nil {
		panic(err)
	}

	h := httputil.NewSingleHostReverseProxy(u)
	h.Transport = NewHighlanderProxy(
		cmd.healthTick,
		cmd.healthExpiry)

	http.ListenAndServe(cmd.bind, h)
}
