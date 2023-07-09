package nmap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestNmapIntegration(t *testing.T) {
	server := &http.Server{Addr: "127.0.0.1:8080"}
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/login\n")
	})
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("server error: %v", err)
		}
	}()

	str, err := Scan("-sC -sV -p 8080 -T5 127.0.0.1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Println("Test output")
	fmt.Printf("%s\n", str)

	const Needle = "8080/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)"
	if !strings.Contains(str, Needle) {
		t.Fatalf("expected Golang net/http server not found in str: %v", str)
	}
	server.Shutdown(context.Background())
}
