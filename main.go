package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coder/websocket"
)

const (
	navURL = "ws://127.0.0.1:8081/navigation"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/navigation", WebSocketProxyHandler)
	log.Println("Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// WebSocketProxyHandler acts as a Websocket reverse proxy between the client and backend.
func WebSocketProxyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Gérer l'authentification ici

	log.Println("proxy handler called")

	u, err := url.Parse(navURL)
	if err != nil {
		http.Error(w, "Backend URL invalid", http.StatusInternalServerError)
		return
	}
	u.RawQuery = r.URL.RawQuery // Forward of the query params

	upgradeCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Update the headers to forward
	header := http.Header{}
	for k, vals := range r.Header {
		// Forward useful headers, but not all for security
		if strings.HasPrefix(strings.ToLower(k), "sec-websocket") ||
			strings.ToLower(k) == "authorization" {
			header[k] = vals
		}
	}

	// Connect to the backend Websocket server
	backendConn, _, err := websocket.Dial(upgradeCtx, u.String(), &websocket.DialOptions{
		HTTPHeader: header,
	})
	if err != nil {
		http.Error(w, "Could not connect to backend", http.StatusBadGateway)
		log.Printf("WebSocket backend dial error: %v", err)
		return
	}
	defer backendConn.Close(websocket.StatusInternalError, "proxy error")

	// Accept incoming client connection
	clientConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// OriginPatterns avec CORS ??
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("WebSocket accept error: %v", err)
		return
	}
	defer clientConn.Close(websocket.StatusInternalError, "proxy error")

	proxyCtx, proxyCancel := context.WithCancel(context.Background())
	defer proxyCancel()

	// Proxy bidirectionnel : client <=> backend
	errc := make(chan error, 2)
	go proxyWebSocket(proxyCtx, clientConn, backendConn, errc)
	go proxyWebSocket(proxyCtx, backendConn, clientConn, errc)

	// Wait for one of the two connection to close or an error
	err = <-errc
	status := websocket.StatusNormalClosure
	reason := "proxy closed"
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		status = websocket.StatusInternalError
		reason = err.Error()
	}
	if err := clientConn.Close(status, reason); err != nil {
		log.Printf("WebSocket client close error: %v", err)
	}
	if err := backendConn.Close(status, reason); err != nil {
		log.Printf("WebSocket backend close error: %v", err)
	}
}

// proxyWebSocket copies messages between the two WS connections.
func proxyWebSocket(ctx context.Context, src, dst *websocket.Conn, errc chan error) {
	for {
		typ, data, err := src.Read(ctx)
		if err != nil {
			errc <- err
			return
		}
		err = dst.Write(ctx, typ, data)
		if err != nil {
			errc <- err
			return
		}
	}
}
