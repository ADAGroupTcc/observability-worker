package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"worker/config"
)

func main() {
	envs, err := config.LoadEnvVars()
	if err != nil {
		log.Fatalf("config.LoadEnvVars: %v", err)
	}
	cl := http.Client{}
	go func() {
		http.HandleFunc("/health", HealthHandler)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", envs.ApiPort), nil); err != nil {
			log.Fatal(fmt.Sprintf("Failed to listen port %s", envs.ApiPort), err)
			os.Exit(1)
		}
	}()
	var wg sync.WaitGroup
	for {
		log.Default().Printf("Start checking APIs healthiness")
		for _, baseUrl := range envs.ApiBaseUrls {
			wg.Add(1)
			go func(baseUrl string) {
				defer wg.Done()
				url, _ := url.Parse(fmt.Sprintf("%s/health", baseUrl))
				req := &http.Request{Method: "GET", URL: url, Header: http.Header{"Api_Token": []string{envs.ApiToken}}}
				res, err := cl.Do(req)
				if err != nil {
					log.Fatal(err)
					log.Fatalf("API %s is unhealthy", baseUrl)
				}
				if res.StatusCode == 200 {
					log.Printf("API %s is healthy", baseUrl)
				} else {
					log.Fatalf("API %s is unhealthy", baseUrl)
				}
			}(baseUrl)
		}
		wg.Wait()
		log.Default().Printf("Finish checking APIs healthiness")
		time.Sleep(time.Duration(envs.PollingInterval) * time.Second)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
