package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
		remoteUrl := fmt.Sprintf("http://localhost:9091/?highlander_weight=%s", os.Args[1])
		r, err := http.NewRequest("GET", remoteUrl, nil)
		if err != nil {
			panic(err)
		}

		r.Header.Set("User-Agent", "SimpleHttpClient/0.1")

		res, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Println("Could not reach server", err)
		} else {
			log.Println(res)
			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
}
