package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdate(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/owner/repo/releases/latest":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"tag_name":"v1.0.0","assets":[{"name":"kvasx","browser_download_url":"%s/asset"}]}`, srv.URL)
		case "/asset":
			fmt.Fprint(w, "binary")
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	apiURL = srv.URL

	dir := t.TempDir()
	p, err := Update("owner", "repo", dir)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(data) != "binary" {
		t.Fatalf("unexpected file contents: %s", string(data))
	}
	if filepath.Base(p) != "kvasx" {
		t.Fatalf("unexpected file name: %s", filepath.Base(p))
	}
}
