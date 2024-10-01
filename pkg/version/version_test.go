package version

import (
	"testing"
)

func TestGet(t *testing.T) {
	// Set up test data
	expectedInfo := Info{
		Commit:  Commit,
		Version: Version,
		Branch:  Branch,
		BuiltAt: BuiltAt,
		Builder: Builder,
	}

	// Call Get() and check the result
	Info, err := Get()
	if err != nil {
		t.Errorf("Error getting Info object: %v", err)
	}
	if Info != expectedInfo {
		t.Errorf("Loaded Info object does not match expected. Got %v, expected %v", Info, expectedInfo)
	}
}

// func TestCommand(t *testing.T) {
//	// Set up test data
//	expectedInfo := Info{
//		Commit:  "commit_here",
//		Version: "version_here",
//		BuiltAt: "date_here",
//	}
//
//	expectedRunE := *cobra.Command {
//		expectedInfo
//	}
//
//	expectedCommand := *cobra.Command {
//		Use:   "version",
//		Short: "Print the version.",
//		Long:  `Print the version and build information.`,
//		RunE:
//	}
// }
