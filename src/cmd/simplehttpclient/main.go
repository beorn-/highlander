package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	for {
		r, err := http.NewRequest("GET", "http://localhost:8091", nil)
		if err != nil {
			panic(err)
		}
		r.Header.Set("X-Caller", os.Args[1])

		res, err := http.DefaultClient.Do(r)

		log.Println(res, err)
		time.Sleep(1 * time.Second)
	}
}
