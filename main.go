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

	navigationProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.NavigationHost, config.NavigationPort))
	http.HandleFunc("/navigation/ws", proxyRedirect(navigationProxy, "/ws"))

	usersProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.UsersHost, config.UsersPort))
	http.HandleFunc("/users", proxyHandle(usersProxy))
	http.HandleFunc("/users/", proxyHandle(usersProxy))
	http.HandleFunc("/login", proxyHandle(usersProxy))
	http.HandleFunc("/logout", proxyHandle(usersProxy))
	http.HandleFunc("/refresh", proxyHandle(usersProxy))
	http.HandleFunc("/register", proxyHandle(usersProxy))

	incidentsProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.IncidentsHost, config.IncidentsPort))
	http.HandleFunc("/incidents", proxyHandle(incidentsProxy))
	http.HandleFunc("/incidents/", proxyHandle(incidentsProxy)) // pour /types/{id} etc.
	http.HandleFunc("/incidents/interactions", proxyHandle(incidentsProxy))
	http.HandleFunc("/incidents/me/history", proxyHandle(incidentsProxy))
	http.HandleFunc("/incidents/types", proxyHandle(incidentsProxy))
	http.HandleFunc("/incidents/types/", proxyHandle(incidentsProxy))

	gisProxy := httputil.NewSingleHostReverseProxy(mustParseURL(config.GisHost, config.GisPort))
	http.HandleFunc("/address", proxyHandle(gisProxy))
	http.HandleFunc("/geocode", proxyHandle(gisProxy))
	http.HandleFunc("/route", proxyHandle(gisProxy))

	if config.Port == "" {
		log.Fatal("No gateway port provided")
	}

	log.Printf("Gateway running on :%s\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
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
