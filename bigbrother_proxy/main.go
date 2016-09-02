package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path:", r.URL.Path)
		spew.Dump("query:", r.URL.Query())
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err:", err)
			return
		}
		fmt.Println("body:", string(body))

		query, err := url.ParseQuery(string(body))
		if err != nil {
			log.Println("err:", err)
			return
		}

		/*
			bigBrotherPayload := struct {
				Latitude  string `json:"Latitude"`
				Longitude string `json:"longitude"`
				Accuracy  string `json:"accuracy"`
				Altitude  string `json:"altitude"`
				Provider  string `json:"provider"`
				Bearing   string `json:"bearing"`
				Speed     string `json:"speed"`
				Time      string `json:"time"`
			}{
				Latitude:  query.Get("latitude"),
				Longitude: query.Get("longitude"),
				Accuracy:  query.Get("accuracy"),
				Altitude:  query.Get("altitude"),
				Provider:  query.Get("provider"),
				Bearing:   query.Get("bearing"),
				Speed:     query.Get("speed"),
				Time:      query.Get("time"),
			}

			j, err := json.Marshal(bigBrotherPayload)
			if err != nil {
				log.Println("err:", err)
				return
			}
			urlStr := "http://localhost:8081/location/bigbrother"
			postBody := ioutil.NopCloser(bytes.NewBuffer(j))
			fmt.Println(string(j))
		*/

		urlStr := "http://localhost:8081/location/bigbrother?" + query.Encode()
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
