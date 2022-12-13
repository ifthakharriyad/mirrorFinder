package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ifthakharriyad/mirrorFinder/mirrors"
)

type response struct {
	FastestURL string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

func main() {
	http.HandleFunc("/fastest-mirror", func(w http.ResponseWriter, r *http.Request) {
		response := findFastest(mirrors.MirrorList)
		resJSON, err := json.Marshal(response)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resJSON)
	})
	port := ":8000"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Starting server on port localhost%s", port)
	log.Fatal(server.ListenAndServe())
}

func findFastest(urls []string) response {
	urlChan := make(chan string)
	latenChan := make(chan time.Duration, len(urls))
	for _, url := range urls {
		go func(url string) {
			start := time.Now()
			_, err := http.Get(url + "iso/latest/b2sums.txt")
			laten := time.Since(start) / time.Second
			if err == nil {
				urlChan <- url
				latenChan <- laten
			}
		}(url)
	}
	return response{
		<-urlChan,
		<-latenChan,
	}
}
