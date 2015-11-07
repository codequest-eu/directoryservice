package directoryservice

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
)

// Filter is a function which takes file metadata and decides whether we care
// about it.
type Filter func(os.FileInfo) bool

// IsGoFile is a filter checking if a file is a Go file.
func IsGoFile(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}
	return filepath.Ext(fi.Name()) == ".go"
}

// IsNotTest is a filter checking if a file is not a Go test file.
func IsNotTest(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

// DirectoryService takes responsibility of creating and removing local
// temporary directories.
type DirectoryService interface {
	// BasePath returns the base path for the DirectoryService.
	BasePath() string

	// FullPath returns a full path given a relative one.
	FullPath(string) string

	// Recurse recursively lists a directory and returns a list of paths to
	// files which pass all of the Filters.
	Recurse(string, ...Filter) ([]string, error)

	// FullPath returns a relative path given a full one.
	RelativePath(string) (string, error)

	// Cleanup removes the whole base directory. It invalidates the current
	// DirectoryService.
	Cleanup() error
}

type directoryServiceImpl string

// NewDirectoryService returns a new instance of DirectoryService. If an empty
// path is passed a temporary directory is created.
func NewDirectoryService(path string) (DirectoryService, error) {
	if path != "" {
		return directoryServiceImpl(path), nil
	}
	tempPath, err := ioutil.TempDir("", "DirectoryService")
	if err != nil {
		return nil, errors.New(err)
	}
	return NewDirectoryService(tempPath)
}

func (d directoryServiceImpl) BasePath() string {
	return string(d)
}

func (d directoryServiceImpl) FullPath(relative string) string {
	return filepath.Join(string(d), relative)
}

func (d directoryServiceImpl) Recurse(directory string, filters ...Filter) ([]string, error) {
	collector := newFileCollector()
	walkErr := filepath.Walk(
		d.FullPath(directory),
		collector.recurseWithFilters(filters...),
	)
	if walkErr != nil {
		return nil, errors.New(walkErr)
	}
	return collector.files, nil
}

func (d directoryServiceImpl) RelativePath(absolute string) (string, error) {
	ret, err := filepath.Rel(string(d), absolute)
	if err != nil {
		return "", errors.New(err)
	}
	return ret, nil
}

func (d directoryServiceImpl) Cleanup() error {
	if err := os.RemoveAll(string(d)); err != nil {
		return errors.New(err)
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
