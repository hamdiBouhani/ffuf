package engine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hamdiBouhani/ffuf/internal/config"
)

func TestRun(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			switch r.URL.Path {
			case "/admin":
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("admin page"))

			case "/login":
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("login page"))

			default:
				http.NotFound(w, r)
			}
		}),
	)

	defer server.Close()

	dir := t.TempDir()

	wordlistPath := filepath.Join(
		dir,
		"words.txt",
	)

	err := os.WriteFile(
		wordlistPath,
		[]byte("admin\nlogin\nmissing\n"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Config{
		URL:      server.URL + "/FUZZ",
		Wordlist: wordlistPath,

		Threads: 3,
		Method:  "GET",

		Timeout: 5 * time.Second,

		MatchStatus: map[int]bool{
			200: true,
		},
	}

	results, err := Run(
		context.Background(),
		cfg,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"got %d results, want 2",
			len(results),
		)
	}

	found := map[string]bool{}

	for _, result := range results {
		found[result.Input] = true

		if result.StatusCode != 200 {
			t.Errorf(
				"%s returned status %d",
				result.Input,
				result.StatusCode,
			)
		}
	}

	if !found["admin"] {
		t.Error("admin result not found")
	}

	if !found["login"] {
		t.Error("login result not found")
	}

	if found["missing"] {
		t.Error("missing should have been filtered")
	}
}
