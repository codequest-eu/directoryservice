// Package directoryservice provides a useful abstraction for working with
// temporary directories and files they contain without the need to explicitly
// deal with the underlying OS's impementation. This can be useful for things
// like temporarily cloning Git repositories.
package directoryservice

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Filter is a function which takes file metadata and decides whether we care
// about it.
type Filter func(os.FileInfo) bool

type serviceImpl string

// TemporaryService returns a new instance of Service which creates a temporary
// directory for its' purposes.
func TemporaryService() (Service, error) {
	tempPath, err := ioutil.TempDir("", "DirectoryService")
	return serviceImpl(tempPath), err
}

// DirectoryService takes an existing directory and creates an instance of
// Service around it. The consturctor verifies that the provided path leads to
// a valid directory.
func DirectoryService(path string) (Service, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}
	return serviceImpl(path), nil
}

func (si serviceImpl) BasePath() string {
	return string(si)
}

func (si serviceImpl) FullPath(relative string) string {
	return filepath.Join(string(si), relative)
}

func (si serviceImpl) Recurse(directory string, filters ...Filter) ([]string, error) {
	collector := newFileCollector()
	walkErr := filepath.Walk(
		si.FullPath(directory),
		collector.recurseWithFilters(filters...),
	)
	if walkErr != nil {
		return nil, walkErr
	}
	return collector.files, nil
}

func (si serviceImpl) RelativePath(absolute string) (string, error) {
	ret, err := filepath.Rel(string(si), absolute)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (si serviceImpl) Cleanup() error {
	if err := os.RemoveAll(string(si)); err != nil {
		return err
	}
	return nil
}

func passesFilters(fi os.FileInfo, filters ...Filter) bool {
	for _, filter := range filters {
		if !filter(fi) {
			return false
		}
	}
	return true
}

type fileCollector struct {
	files []string
}

func newFileCollector() *fileCollector {
	return &fileCollector{make([]string, 0)}
}

func (f *fileCollector) recurseWithFilters(filters ...Filter) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir // continue at all cost
		}
		if passesFilters(fi, filters...) {
			f.files = append(f.files, path)
		}
		return nil
	}
}
