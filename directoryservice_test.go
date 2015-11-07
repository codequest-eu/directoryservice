package directoryservice

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/suite"
)

const (
	customPath = "/bacon"
	fileName   = "fry.txt"
)

var fullPath = filepath.Join(customPath, fileName)

type DirectoryServiceTestSuite struct {
	suite.Suite
	service DirectoryService
	cleanup bool
}

func (s *DirectoryServiceTestSuite) TearDownTest() {
	if s.cleanup {
		s.Nil(s.service.Cleanup())
	}
}

func (s *DirectoryServiceTestSuite) withDirectory(name string) {
	service, err := NewDirectoryService(name)
	s.Require().Nil(err)
	s.service = service
}

func (s *DirectoryServiceTestSuite) withTemporaryDirectory() {
	s.withDirectory("")
	s.cleanup = true
}

func (s *DirectoryServiceTestSuite) withCustomDirectory() {
	s.withDirectory(customPath)
}

func (s *DirectoryServiceTestSuite) TestTempdirLifecycle() {
	s.withTemporaryDirectory()
	s.NotEmpty(s.service.BasePath())
}

func (s *DirectoryServiceTestSuite) TestRecurse() {
	s.withTemporaryDirectory()
	for _, file := range []string{"baconium.go", "cabbagium.rb"} {
		path := s.service.FullPath(file)
		s.Require().Nil(ioutil.WriteFile(path, []byte("bacon"), 0777))
	}
	files, err := s.service.Recurse(".", IsGoFile)
	s.Require().Nil(err)
	s.Len(files, 1)
	s.Contains(files[0], "baconium.go")
}

func (s *DirectoryServiceTestSuite) TestCustomDirectory() {
	s.withCustomDirectory()
	s.Equal(s.service.BasePath(), customPath)
	s.Equal(s.service.FullPath(fileName), fullPath)
}

func (s *DirectoryServiceTestSuite) TestRelativePathSuccess() {
	s.withCustomDirectory()
	result, err := s.service.RelativePath(fullPath)
	s.Require().Nil(err)
	s.Equal(fileName, result)
}

func (s *DirectoryServiceTestSuite) TestRelativePathError() {
	s.withCustomDirectory()
	_, err := s.service.RelativePath(fileName)
	s.Require().NotNil(err)
	s.IsType((*errors.Error)(nil), err)
}

func TestDirectoryService(t *testing.T) {
	suite.Run(t, new(DirectoryServiceTestSuite))
}
