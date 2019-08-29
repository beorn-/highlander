package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr + " " + r.Method + " " + r.URL.String() + " " + r.Proto + " " + r.UserAgent() + " " + r.Header.Get("X-Highlander-Weight"))
	})
	http.ListenAndServe(":9090", nil)
}
