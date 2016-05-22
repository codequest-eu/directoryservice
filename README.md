# directoryservice [![codebeat badge](https://codebeat.co/badges/fe5a26da-9b90-47be-b8b3-6ce870adb3ff)](https://codebeat.co/projects/directoryservice-master) [![GoDoc](https://godoc.org/github.com/codequest-eu/directoryservice?status.svg)](https://godoc.org/github.com/codequest-eu/directoryservice) [![GoCover](http://gocover.io/_badge/github.com/codequest-eu/directoryservice)](http://gocover.io/github.com/codequest-eu/directoryservice)

Package directoryservice provides a useful abstraction for working with temporary directories and files they contain without the need to explicitly deal with the underlying OS's impementation. This can be useful for things like temporarily cloning Git repositories.

[![wercker status](https://app.wercker.com/status/aea4ead7c0e7add1f0448f17f18d03f2/m/master "wercker status")](https://app.wercker.com/project/bykey/aea4ead7c0e7add1f0448f17f18d03f2)

## Usage

    import "github.com/codequest-eu/directoryservice"


#### type Service

```go
type Service interface {
	// BasePath returns the base path for the Service.
	BasePath() string

	// FullPath returns a full path given a relative one.
	FullPath(string) string

	// Recurse recursively lists a directory and returns a list of paths to
	// files which pass all of the Filters.
	Recurse(string, ...Filter) ([]string, error)

	// FullPath returns a relative path given a full one.
	RelativePath(string) (string, error)

	// Cleanup removes the whole base directory. It invalidates the current
	// Service.
	Cleanup() error
}
```

Service provides a number of helpful utility methods for a single directory.

#### func TemporaryService

```go
func TemporaryService() (Service, error)
```

TemporaryService returns a new instance of Service which creates a temporary directory for its' purposes.

#### func DirectoryService

```go
func DirectoryService(path string) (Service, error)
```

DirectoryService takes an existing directory and creates an instance of Service around it. The consturctor verifies that the provided path leads to a valid directory.

#### type Filter

```go
type Filter func(os.FileInfo) bool
```

Filter is a function which takes file metadata and decides whether we care about it.
