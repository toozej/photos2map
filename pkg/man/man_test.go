package man

import (
	"testing"
)

func TestNewManCmd(t *testing.T) {
	// test each field of NewManCmd Cobra command

	expectedUse := "man"
	if NewManCmd().Use != expectedUse {
		t.Errorf("Unexpected command use text: got %q, expected %q", NewManCmd().Use, expectedUse)
	}

	expectedShort := "Generates golang-starter's command line manpages"
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

	//	c := cobra.Command()
	//	expectedManPage := mcoral.NewManPage(1, c.Root())
	//	if NewManCmd().RunE(c) != expectedManPage {
	//		t.Errorf("Unexpected command Hidden field: got %t, expected %t", NewManCmd().Hidden, expectedHidden)
	//	}

}
