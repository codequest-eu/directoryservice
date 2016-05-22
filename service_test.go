package directoryservice

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	service Service
	cleanup bool
}

func (s *serviceTestSuite) TearDownTest() {
	if s.cleanup {
		s.Nil(s.service.Cleanup())
	}
}

func (s *serviceTestSuite) withTestDirectory() {
	service, err := DirectoryService("testdata")
	s.NoError(err)
	s.service = service
	s.cleanup = false
}

func (s *serviceTestSuite) withTemporaryDirectory() {
	service, err := TemporaryService()
	s.NoError(err)
	s.service = service
	s.cleanup = true
}

func (s *serviceTestSuite) TestTempdirLifecycle() {
	s.withTemporaryDirectory()
	s.NotEmpty(s.service.BasePath())
}

func (s *serviceTestSuite) TestRecurse() {
	s.withTemporaryDirectory()
	for _, file := range []string{"baconium.go", "cabbagium.rb"} {
		path := s.service.FullPath(file)
		s.NoError(ioutil.WriteFile(path, []byte("bacon"), 0777))
	}
	files, err := s.service.Recurse(".", func(fi os.FileInfo) bool {
		return filepath.Ext(fi.Name()) == ".go"
	})
	s.NoError(err)
	s.Len(files, 1)
	s.Contains(files[0], "baconium.go")
}

func (s *serviceTestSuite) TestExistingDirectory() {
	s.withTestDirectory()
	s.Equal(s.service.BasePath(), "testdata")
	s.Equal(s.service.FullPath("bacon.txt"), "testdata/bacon.txt")
}

func (s *serviceTestSuite) TestRelativePathSuccess() {
	s.withTestDirectory()
	result, err := s.service.RelativePath("testdata/bacon.txt")
	s.Nil(err)
	s.Equal("bacon.txt", result)
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}
