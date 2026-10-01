package wordlist

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "words.txt")

	content := `
# comment

admin
login

api
# another comment
test
`

	err := os.WriteFile(
		path,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"admin",
		"login",
		"api",
		"test",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"Read() = %#v, want %#v",
			got,
			want,
		)
	}
}

func TestReadMissingFile(t *testing.T) {
	_, err := Read("/does/not/exist")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
