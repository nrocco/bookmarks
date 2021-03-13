package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewGood(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	store, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if store.fs != tmpDir {
		t.Fatalf("Store.fs=%s != tmpDir=%s", store.fs, tmpDir)
	}
}

func TestNewPathDoesNotExist(t *testing.T) {
	_, err := New("/i/do/not/exist")
	if err == nil {
		t.Fatal("This should have failed")
	}

	if "open /i/do/not/exist/data.db: no such file or directory" != err.Error() {
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

	if _, err := New(tmpDir); err == nil {
		t.Fatal("This should have failed, but it did not")
	}
}
