package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDo(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Header.Get("X-Test") != "hello" {
				t.Errorf(
					"X-Test = %q, want %q",
					r.Header.Get("X-Test"),
					"hello",
				)
			}

			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(
				[]byte("hello world\nthis is a test"),
			)
		}),
	)

	defer server.Close()

	client := New(
		5*time.Second,
		map[string]string{
			"X-Test": "hello",
		},
		"myfuzz-test",
	)

	result, err := client.Do(http.MethodGet, server.URL, "")

	if err != nil {
		t.Fatal(err)
	}

	if result.StatusCode != 200 {
		t.Fatalf("StatusCode = %d, want 200", result.StatusCode)
	}

	if result.Size != 26 {
		t.Fatalf("Size = %d, want 26", result.Size)
	}

	if result.Words != 6 {
		t.Fatalf("Words = %d, want 6", result.Words)
	}

	if result.Lines != 2 {
		t.Fatalf("Lines = %d, want 2", result.Lines)
	}

	if result.Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}

func TestClientUserAgent(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.UserAgent(); got != "myfuzz-test" {
				t.Errorf("User-Agent = %q, want %q", got, "myfuzz-test")
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	client := New(5*time.Second, nil, "myfuzz-test")

	result, err := client.Do(http.MethodGet, server.URL, "")

	if err != nil {
		t.Fatal(err)
	}

	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("StatusCode = %d, want %d", result.StatusCode, http.StatusNoContent)
	}
}
