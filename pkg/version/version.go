package version

import "fmt"

var (
	// Version is the current version of the application.
	// It is set during the build process using ldflags.
	Version = "1.0.0"

	// Commit is the git commit hash at which the application was built.
	Commit = "unknown"

	// BuildDate is the date at which the application was built.
	BuildDate = "unknown"
)

// Info returns a formatted string with version information.
func Info() string {
	return fmt.Sprintf("v%s (%s) build at %s", Version, Commit, BuildDate)
}
