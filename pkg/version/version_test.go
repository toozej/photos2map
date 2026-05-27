package version

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGet_Defaults(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "local"
	Commit = ""
	Branch = ""
	BuiltAt = ""
	Builder = ""

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "local" {
		t.Errorf("expected Version='local', got '%s'", info.Version)
	}
	if info.Commit != "" {
		t.Errorf("expected Commit='', got '%s'", info.Commit)
	}
	if info.Branch != "" {
		t.Errorf("expected Branch='', got '%s'", info.Branch)
	}
	if info.BuiltAt != "" {
		t.Errorf("expected BuiltAt='', got '%s'", info.BuiltAt)
	}
	if info.Builder != "" {
		t.Errorf("expected Builder='', got '%s'", info.Builder)
	}
}

func TestGet_LocalBuild(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "v1.2.3-2-gabcdef"
	Commit = "abcdef1"
	Branch = "feature/test"
	BuiltAt = "2026-05-04T21:00:00Z"
	Builder = "developer@workstation"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "v1.2.3-2-gabcdef" {
		t.Errorf("expected Version='v1.2.3-2-gabcdef', got '%s'", info.Version)
	}
	if info.Commit != "abcdef1" {
		t.Errorf("expected Commit='abcdef1', got '%s'", info.Commit)
	}
	if info.Branch != "feature/test" {
		t.Errorf("expected Branch='feature/test', got '%s'", info.Branch)
	}
	if info.BuiltAt != "2026-05-04T21:00:00Z" {
		t.Errorf("expected BuiltAt='2026-05-04T21:00:00Z', got '%s'", info.BuiltAt)
	}
	if info.Builder != "developer@workstation" {
		t.Errorf("expected Builder='developer@workstation', got '%s'", info.Builder)
	}
}

func TestGet_GoReleaserBuild(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "1.2.3"
	Commit = "a1b2c3d4e5f6"
	Branch = "main"
	BuiltAt = "2026-05-04T12:00:00Z"
	Builder = "goreleaser"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "1.2.3" {
		t.Errorf("expected Version='1.2.3', got '%s'", info.Version)
	}
	if info.Commit != "a1b2c3d4e5f6" {
		t.Errorf("expected Commit='a1b2c3d4e5f6', got '%s'", info.Commit)
	}
	if info.Branch != "main" {
		t.Errorf("expected Branch='main', got '%s'", info.Branch)
	}
	if info.BuiltAt != "2026-05-04T12:00:00Z" {
		t.Errorf("expected BuiltAt='2026-05-04T12:00:00Z', got '%s'", info.BuiltAt)
	}
	if info.Builder != "goreleaser" {
		t.Errorf("expected Builder='goreleaser', got '%s'", info.Builder)
	}
}

func TestGet_DockerBuild(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "v2.0.0"
	Commit = "deadbeef"
	Branch = "release/v2.0.0"
	BuiltAt = "2026-06-01T08:30:00Z"
	Builder = "docker"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "v2.0.0" {
		t.Errorf("expected Version='v2.0.0', got '%s'", info.Version)
	}
	if info.Commit != "deadbeef" {
		t.Errorf("expected Commit='deadbeef', got '%s'", info.Commit)
	}
	if info.Branch != "release/v2.0.0" {
		t.Errorf("expected Branch='release/v2.0.0', got '%s'", info.Branch)
	}
	if info.Builder != "docker" {
		t.Errorf("expected Builder='docker', got '%s'", info.Builder)
	}
}

func TestGet_PartiallyPopulated(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "v0.1.0"
	Commit = "abc123"
	Branch = ""
	BuiltAt = ""
	Builder = ""

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "v0.1.0" {
		t.Errorf("expected Version='v0.1.0', got '%s'", info.Version)
	}
	if info.Commit != "abc123" {
		t.Errorf("expected Commit='abc123', got '%s'", info.Commit)
	}
	if info.Branch != "" {
		t.Errorf("expected Branch='', got '%s'", info.Branch)
	}
	if info.BuiltAt != "" {
		t.Errorf("expected BuiltAt='', got '%s'", info.BuiltAt)
	}
	if info.Builder != "" {
		t.Errorf("expected Builder='', got '%s'", info.Builder)
	}
}

func TestGet_UnknownValues(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "unknown"
	Commit = "unknown"
	Branch = "unknown"
	BuiltAt = "unknown"
	Builder = "unknown"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, err := Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if info.Version != "unknown" {
		t.Errorf("expected Version='unknown', got '%s'", info.Version)
	}
	if info.Commit != "unknown" {
		t.Errorf("expected Commit='unknown', got '%s'", info.Commit)
	}
	if info.Branch != "unknown" {
		t.Errorf("expected Branch='unknown', got '%s'", info.Branch)
	}
}

func TestGet_ReturnsNoError(t *testing.T) {
	_, err := Get()
	if err != nil {
		t.Errorf("Get() should never return an error, got: %v", err)
	}
}

func TestCommand_Structure(t *testing.T) {
	cmd := Command()

	if cmd.Use != "version" {
		t.Errorf("expected Use='version', got '%s'", cmd.Use)
	}
	if cmd.Short != "Print the version." {
		t.Errorf("expected Short='Print the version.', got '%s'", cmd.Short)
	}
	if cmd.Long != "Print the version and build information." {
		t.Errorf("expected Long='Print the version and build information.', got '%s'", cmd.Long)
	}
	if cmd.RunE == nil {
		t.Error("expected RunE to be set, got nil")
	}
}

func TestCommand_JSONOutput(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "1.0.0"
	Commit = "abc123"
	Branch = "main"
	BuiltAt = "2026-01-01T00:00:00Z"
	Builder = "test"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	cmd := Command()

	var stdout strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command() execution failed: %v", err)
	}

	output := stdout.String()
	if output == "" {
		t.Fatal("Command() produced no output")
	}

	var result Info
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
		t.Fatalf("Command() output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if result.Version != "1.0.0" {
		t.Errorf("expected Version='1.0.0' in JSON, got '%s'", result.Version)
	}
	if result.Commit != "abc123" {
		t.Errorf("expected Commit='abc123' in JSON, got '%s'", result.Commit)
	}
	if result.Branch != "main" {
		t.Errorf("expected Branch='main' in JSON, got '%s'", result.Branch)
	}
	if result.BuiltAt != "2026-01-01T00:00:00Z" {
		t.Errorf("expected BuiltAt='2026-01-01T00:00:00Z' in JSON, got '%s'", result.BuiltAt)
	}
	if result.Builder != "test" {
		t.Errorf("expected Builder='test' in JSON, got '%s'", result.Builder)
	}
}

func TestCommand_JSONOutput_DefaultValues(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "local"
	Commit = ""
	Branch = ""
	BuiltAt = ""
	Builder = ""

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	cmd := Command()

	var stdout strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command() execution failed: %v", err)
	}

	output := stdout.String()

	var result Info
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
		t.Fatalf("Command() output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if result.Version != "local" {
		t.Errorf("expected Version='local' in JSON, got '%s'", result.Version)
	}
	if result.Commit != "" {
		t.Errorf("expected Commit='' in JSON, got '%s'", result.Commit)
	}
}

func TestInfo_FieldOrder(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	Version = "v1"
	Commit = "c1"
	Branch = "b1"
	BuiltAt = "t1"
	Builder = "x1"

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	info, _ := Get()

	if info.Commit != "c1" {
		t.Errorf("Info.Commit not populated correctly, got '%s'", info.Commit)
	}
	if info.Version != "v1" {
		t.Errorf("Info.Version not populated correctly, got '%s'", info.Version)
	}
	if info.Branch != "b1" {
		t.Errorf("Info.Branch not populated correctly, got '%s'", info.Branch)
	}
	if info.BuiltAt != "t1" {
		t.Errorf("Info.BuiltAt not populated correctly, got '%s'", info.BuiltAt)
	}
	if info.Builder != "x1" {
		t.Errorf("Info.Builder not populated correctly, got '%s'", info.Builder)
	}
}

func TestGet_VariableSync(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origBranch := Branch
	origBuiltAt := BuiltAt
	origBuilder := Builder

	defer func() {
		Version = origVersion
		Commit = origCommit
		Branch = origBranch
		BuiltAt = origBuiltAt
		Builder = origBuilder
	}()

	Version = "before"
	Commit = "before"
	Branch = "before"
	BuiltAt = "before"
	Builder = "before"

	info1, _ := Get()
	if info1.Version != "before" {
		t.Errorf("expected 'before', got '%s'", info1.Version)
	}

	Version = "after"
	Commit = "after"
	Branch = "after"
	BuiltAt = "after"
	Builder = "after"

	info2, _ := Get()
	if info2.Version != "after" {
		t.Errorf("expected 'after', got '%s'", info2.Version)
	}
}
