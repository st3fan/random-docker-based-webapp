package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, "Hello, world!")
	})
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("random-docker-based-webapp listening on %s", listener.Addr())
	log.Fatal(http.Serve(listener, accessLog(mux, log.New(os.Stdout, "", 0))))
}

type loggedResponse struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *loggedResponse) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.status == 0 && (code >= 200 || code == http.StatusSwitchingProtocols) {
		w.status = code
	}
}

func (w *loggedResponse) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += int64(n)
	return n, err
}

func accessLog(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		response := &loggedResponse{ResponseWriter: w}
		next.ServeHTTP(response, r)

		status := response.status
		if status == 0 {
			status = http.StatusOK
		}
		size := "-"
		if r.Method != http.MethodHead && response.bytes > 0 {
			size = strconv.FormatInt(response.bytes, 10)
		}
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = "-"
		}
		requestLine := strconv.QuoteToASCII(r.Method + " " + r.RequestURI + " " + r.Proto)
		logger.Printf("%s - - [%s] %s %d %s", host, start.Format("02/Jan/2006:15:04:05 -0700"), requestLine, status, size)
	})
}
