package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dedis/livos/controller"
)

type key int

//go:embed index.html
var siteweb embed.FS
var content embed.FS

//go:embed static
var static embed.FS

const (
	requestIDKey key = 0
)

func main() {

	listenAddr := ":9000"
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:  tracing(nextRequestID)(logging(logger)(mux)),
		ErrorLog: logger,
	}

	ctrl := controller.NewController(content)

	mux.HandleFunc("/", ctrl.HandleHome)

	// serve assets
	mux.Handle("/static/", http.FileServer(http.FS(static)))¨
	mux.Handle("/siteweb/", http.FileServer(http.FS(siteweb)))¨

	// create connection
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("failed to create conn '%s': %v", listenAddr, err)
		return
	}

	// launch server
	err = server.Serve(ln)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}

	mux.HandleFunc("/quitserver", ctrl.HandleQuit)
}

// Utility function for logging

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
