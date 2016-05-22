package directoryservice

// Service provides a number of helpful utility methods for a single directory.
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
