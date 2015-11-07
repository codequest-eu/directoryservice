# directoryservice
--
    import "github.com/codequest-eu/directoryservice"

Package directoryservice provides a useful abstraction for working with
temporary directories and files they contain without the need to explicitly deal
with the underlying OS's impementation. This can be useful for things like
temporarily cloning Git repositories.

## Usage

#### func  IsGoFile

```go
func IsGoFile(fi os.FileInfo) bool
```
IsGoFile is a filter checking if a file is a Go file.

#### func  IsNotTest

```go
func IsNotTest(fi os.FileInfo) bool
```
IsNotTest is a filter checking if a file is not a Go test file.

#### type DirectoryService

```go
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
```

DirectoryService takes responsibility of creating and removing local temporary
directories.

#### func  NewDirectoryService

```go
func NewDirectoryService(path string) (DirectoryService, error)
```
NewDirectoryService returns a new instance of DirectoryService. If an empty path
is passed a temporary directory is created.

#### type Filter

```go
type Filter func(os.FileInfo) bool
```

Filter is a function which takes file metadata and decides whether we care about
it.
