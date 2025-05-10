package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	//usersProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.UsersHost, config.UsersHost))
	//incidentsProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.IncidentsHost, config.IncidentsPort))
	//gisProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.GisHost, config.GisPort))
	navigationProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.NavigationHost, config.NavigationPort))

	http.HandleFunc("/navigation/ws", proxyRedirect(navigationProxy, "/ws"))

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// proxyRedirect prend un proxy et une route de réécriture optionnelle
func proxyRedirect(proxy *httputil.ReverseProxy, rewriteTo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rewriteTo != "" { // Si une réécriture est demandée, on modifie le chemin
			r.URL.Path = rewriteTo
		}

		proxy.ServeHTTP(w, r) // Passer la requête au proxy
	}
}

func proxyHandle(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return proxyRedirect(proxy, "")
}

func mustParseURL(host, port string) *url.URL {
	rawURL := fmt.Sprintf("http://%s:%s", strings.TrimSpace(host), strings.TrimSpace(port))
	parsed, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid URL %s: %v", rawURL, err)
	}
	return parsed
}
