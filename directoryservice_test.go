package services

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	customPath = "/bacon"
	fileName   = "fry.txt"
)

var fullPath = fmt.Sprintf("%s/%s", customPath, fileName)

func TestNewDirectoryServiceDefault(t *testing.T) {
	ds, err := NewDirectoryService("")
	if err != nil {
		t.Errorf("NewDirectoryService = %v, expected nil", err)
	}
	if ds.BasePath() == "" {
		t.Error("ds.BasePath = \"\" expected actual path")
	}
	if err := ds.Cleanup(); err != nil {
		t.Errorf("ds.Close = %v, expected nil", err)
	}
}

func TestNewDirectoryServiceCustom(t *testing.T) {
	ds, err := NewDirectoryService(customPath)
	if err != nil {
		t.Errorf("NewDirectoryService = %v, expected nil", err)
	}
	basePath := ds.BasePath()
	if basePath != customPath {
		t.Errorf("ds.BasePath = %q expected %q", basePath, customPath)
	}
}

func TestFullPath(t *testing.T) {
	ds, _ := NewDirectoryService(customPath)
	actual := ds.FullPath(fileName)
	if actual != fullPath {
		t.Errorf("ds.FullPath = %q, expected %q", actual, fullPath)
	}
}

func TestRecurse(t *testing.T) {
	ds, _ := NewDirectoryService("")
	defer ds.Cleanup()
	if err := ioutil.WriteFile(ds.FullPath("baconium.go"), []byte("bacon"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile = %v, expected nil", err)
	}
	if err := ioutil.WriteFile(ds.FullPath("cabbagium.rb"), []byte("cabbage"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile = %v, expected nil", err)
	}
	files, err := ds.Recurse(".", IsGoFile)
	if err != nil {
		t.Fatalf("ds.Recurse err = %v, expected nil", err)
	}
	listLen := len(files)
	if listLen != 1 {
		t.Fatalf("ds.Recurse returned %d files, expected 1", listLen)
	}
	if !strings.HasSuffix(files[0], "baconium.go") {
		t.Errorf("file name = %q, expected */baconium.go", files[0])
	}
}

func TestRelativePath(t *testing.T) {
	ds, _ := NewDirectoryService(customPath)
	actual, err := ds.RelativePath(fullPath)
	if err != nil {
		t.Fatalf("ds.RelativePath err = %v, expected nil", err)
	}
	if actual != fileName {
		t.Errorf("ds.RelativePath = %q, expected %q", actual, fileName)
	}
}
