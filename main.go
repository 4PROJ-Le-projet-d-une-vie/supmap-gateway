package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const navURL = "http://localhost:8081"

func main() {
	target, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	navigationProxy := httputil.NewSingleHostReverseProxy(target)

	navigationProxy.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetURL(target)
		pr.Out.URL.Path = "/ws"
	}

	http.HandleFunc("/navigation/ws", func(w http.ResponseWriter, r *http.Request) {
		navigationProxy.ServeHTTP(w, r)
	})

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
