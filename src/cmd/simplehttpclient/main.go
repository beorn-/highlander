package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		panic("Syntax error: " + os.Args[0] + " <weight>")
	}
	for {
		r, err := http.NewRequest("GET", "http://localhost:9091", nil)
		if err != nil {
			panic(err)
		}

		r.Header.Set("X-Highlander-Weight", os.Args[1])

		res, err := http.DefaultClient.Do(r)

		log.Println(res, err)
		time.Sleep(1 * time.Second)
	}
}
