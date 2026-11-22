package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestAccessLog(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hi")
	})
	for _, tc := range []struct {
		method, path, result string
	}{
		{"GET", "/", "200 2"},
		{"HEAD", "/", "200 -"},
		{"GET", "/missing", "404 19"},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			var output bytes.Buffer
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.RemoteAddr = "[2001:db8::1]:4231"
			accessLog(mux, log.New(&output, "", 0)).ServeHTTP(httptest.NewRecorder(), req)
			pattern := `^2001:db8::1 - - \[\d{2}/[A-Z][a-z]{2}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}\] "` +
				regexp.QuoteMeta(tc.method+" "+tc.path+" HTTP/1.1") + `" ` + tc.result + "\n$"
			if !regexp.MustCompile(pattern).MatchString(output.String()) {
				t.Errorf("unexpected access log %q", output.String())
			}
		})
	}
}

func TestAccessLogEscapesRequestLine(t *testing.T) {
	var output bytes.Buffer
	req := httptest.NewRequest("GET", "/", nil)
	req.RequestURI = "/\"\n"
	accessLog(http.NewServeMux(), log.New(&output, "", 0)).ServeHTTP(httptest.NewRecorder(), req)
	if !strings.Contains(output.String(), `"GET /\"\n HTTP/1.1"`) || strings.Count(output.String(), "\n") != 1 {
		t.Fatal(output.String())
	}
}
