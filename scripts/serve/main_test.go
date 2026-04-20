package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDir_Valid(t *testing.T) {
	tmp := t.TempDir()

	got, err := resolveDir(tmp)
	if err != nil {
		t.Fatalf("resolveDir(%q) returned error: %v", tmp, err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
}

func TestResolveDir_Missing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	if _, err := resolveDir(missing); err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

func TestResolveDir_NotADirectory(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "regular.txt")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := resolveDir(file); err == nil {
		t.Fatal("expected error when path is a file, got nil")
	}
}

func TestHandler_WasmContentType(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "game.wasm"), []byte("\x00asm"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	rr := serve(t, tmp, "/game.wasm")

	if got := rr.Header().Get("Content-Type"); got != "application/wasm" {
		t.Errorf("Content-Type = %q, want application/wasm", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
}

func TestHandler_NonWasmKeepsDefaultContentType(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "index.html"), []byte("<html></html>"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	rr := serve(t, tmp, "/index.html")

	if ct := rr.Header().Get("Content-Type"); ct == "application/wasm" {
		t.Errorf("Content-Type should not be application/wasm for .html, got %q", ct)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
}

func TestHandler_CacheControlAlwaysSet(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "plain.txt"), []byte("hi"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	rr := serve(t, tmp, "/plain.txt")

	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
}

func TestHandler_PathTraversalBlocked(t *testing.T) {
	root := t.TempDir()
	secretsDir := filepath.Dir(root)
	secret := filepath.Join(secretsDir, "secret.txt")
	if err := os.WriteFile(secret, []byte("top-secret"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(secret) })

	rr := serve(t, root, "/../secret.txt")

	if strings.Contains(rr.Body.String(), "top-secret") {
		t.Fatalf("path traversal leaked file contents: %q", rr.Body.String())
	}
}

func TestHandler_MissingFileReturns404(t *testing.T) {
	rr := serve(t, t.TempDir(), "/nope.wasm")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func serve(t *testing.T, dir, target string) *httptest.ResponseRecorder {
	t.Helper()
	h := newHandler(http.Dir(dir))
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
