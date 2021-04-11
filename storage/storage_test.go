package storage

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewGood(t *testing.T) {
	if _, err := New(context.Background(), "file::memory:"); err != nil {
		t.Fatal(err)
	}
}

func TestNewPathDoesNotExist(t *testing.T) {
	_, err := New(context.Background(), "/i/do/not/exist.db")
	if err == nil {
		t.Fatal("This should have failed")
	}

	if "unable to open database file: out of memory (14)" != err.Error() {
		t.Fatal(err)
	}
}

func TestNewCannotCreateDatabase(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "data.db"), []byte("hello\ngo\n"), 0644); err != nil {
		t.Fatal("Could not create a data.db file for testing")
	}

	if _, err := New(context.Background(), tmpDir); err == nil {
		t.Fatal("This should have failed, but it did not")
	}
}
