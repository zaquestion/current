package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err:", err)
			return
		}

		query, err := url.ParseQuery(string(body))
		if err != nil {
			log.Println("err:", err)
			return
		}

		uri := os.Getenv("BB_PROXY_URI")
		urlStr := uri + "?" + query.Encode()
		fmt.Println(urlStr)

		_, err = http.Get(urlStr)
		if err != nil {
			log.Println(err)
			return
		}
	})

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
