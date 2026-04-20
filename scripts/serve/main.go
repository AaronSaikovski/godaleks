// Command serve is a small static file server for local WASM testing.
// It sets the correct Content-Type for .wasm files and disables caching
// so rebuilds are picked up on the next reload.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("dir", "web", "directory to serve")
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()

	absDir, err := resolveDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	handler := newHandler(http.Dir(absDir))

	log.Printf("serving %s on http://%s", absDir, *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}

func resolveDir(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", &os.PathError{Op: "serve", Path: abs, Err: os.ErrInvalid}
	}
	return abs, nil
}

func newHandler(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}
		w.Header().Set("Cache-Control", "no-store")
		fileServer.ServeHTTP(w, r)
	})
}
