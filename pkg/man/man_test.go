package man

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewManCmd(t *testing.T) {
	expectedUse := "man"
	if NewManCmd().Use != expectedUse {
		t.Errorf("Unexpected command use text: got %q, expected %q", NewManCmd().Use, expectedUse)
	}

	expectedShort := "Generates photos2map's command line manpages"
	if NewManCmd().Short != expectedShort {
		t.Errorf("Unexpected command short text: got %q, expected %q", NewManCmd().Short, expectedShort)
	}

	expectedSilenceUsage := true
	if NewManCmd().SilenceUsage != expectedSilenceUsage {
		t.Errorf("Unexpected command SilenceUsage field: got %t, expected %t", NewManCmd().SilenceUsage, expectedSilenceUsage)
	}

	expectedDisableFlagsInUseLine := true
	if NewManCmd().DisableFlagsInUseLine != expectedDisableFlagsInUseLine {
		t.Errorf("Unexpected command DisableFlagsInUseLine field: got %t, expected %t", NewManCmd().DisableFlagsInUseLine, expectedDisableFlagsInUseLine)
	}

	expectedHidden := true
	if NewManCmd().Hidden != expectedHidden {
		t.Errorf("Unexpected command Hidden field: got %t, expected %t", NewManCmd().Hidden, expectedHidden)
	}
}

func TestNewManCmd_NoArgs(t *testing.T) {
	cmd := NewManCmd()
	if err := cmd.Args(cmd, []string{}); err != nil {
		t.Errorf("expected no error with zero args, got: %v", err)
	}
}

func TestNewManCmd_RejectsArgs(t *testing.T) {
	cmd := NewManCmd()
	if err := cmd.Args(cmd, []string{"extra"}); err == nil {
		t.Error("expected error when args provided to man command, got nil")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	var out bytes.Buffer
	_, _ = out.ReadFrom(r)
	os.Stdout = oldStdout

	return out.String()
}

func TestNewManCmd_RunE(t *testing.T) {
	rootCmd := &cobra.Command{Use: "testroot", Short: "Test Root"}
	manCmd := NewManCmd()
	rootCmd.AddCommand(manCmd)

	manCmd.SetArgs([]string{})
	output := captureStdout(t, func() {
		err := manCmd.Execute()
		if err != nil {
			t.Fatalf("man command execution failed: %v", err)
		}
	})

	if output == "" {
		t.Error("expected man command to produce output, got empty string")
	}

	if !strings.Contains(output, "testroot") {
		t.Errorf("expected man page output to contain root command name 'testroot', got: %q", output[:min(len(output), 200)])
	}
}

func TestNewManCmd_ManPageGeneration(t *testing.T) {
	rootCmd := &cobra.Command{Use: "myapp", Short: "My Application"}
	manCmd := NewManCmd()
	rootCmd.AddCommand(manCmd)

	manCmd.SetArgs([]string{})
	err := manCmd.Execute()
	if err != nil {
		t.Fatalf("man command execution failed: %v", err)
	}
}

func TestNewManCmd_RunEWithSubcommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "myapp", Short: "My Application"}
	subCmd := &cobra.Command{Use: "sub", Short: "A subcommand", Run: func(cmd *cobra.Command, args []string) {}}
	rootCmd.AddCommand(subCmd)

	manCmd := NewManCmd()
	rootCmd.AddCommand(manCmd)

	manCmd.SetArgs([]string{})
	output := captureStdout(t, func() {
		err := manCmd.Execute()
		if err != nil {
			t.Fatalf("man command execution failed: %v", err)
		}
	})

	if !strings.Contains(output, "sub") {
		t.Errorf("expected man page output to mention subcommand 'sub', got: %q", output[:min(len(output), 500)])
	}
}

func TestNewManCmd_RunEWithRootFlags(t *testing.T) {
	rootCmd := &cobra.Command{Use: "flagapp", Short: "Flag Application"}
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.Flags().String("name", "", "Application name")

	manCmd := NewManCmd()
	rootCmd.AddCommand(manCmd)

	manCmd.SetArgs([]string{})
	output := captureStdout(t, func() {
		err := manCmd.Execute()
		if err != nil {
			t.Fatalf("man command execution failed: %v", err)
		}
	})

	if !strings.Contains(output, "flagapp") {
		t.Errorf("expected man page output to contain 'flagapp', got: %q", output[:min(len(output), 200)])
	}
}

func TestNewManCmd_ReturnsCobraCommand(t *testing.T) {
	cmd := NewManCmd()
	var _ *cobra.Command = cmd
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
