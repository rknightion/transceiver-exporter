package main

import "testing"

// strptr returns a pointer to s, matching the *string type of the flag globals.
func strptr(s string) *string { return &s }

func TestCompileRegexFlags(t *testing.T) {
	// Save and restore the flag globals so tests don't leak into each other.
	origExclude, origInclude := excludeInterfacesRegex, includeInterfacesRegex
	t.Cleanup(func() {
		excludeInterfacesRegex, includeInterfacesRegex = origExclude, origInclude
	})

	t.Run("valid regexes compile", func(t *testing.T) {
		excludeInterfacesRegex = strptr("^eth[0-9]+$")
		includeInterfacesRegex = strptr("^swp")
		if err := compileRegexFlags(); err != nil {
			t.Fatalf("compileRegexFlags() unexpected error: %v", err)
		}
		if !excludeInterfacesRegexCompiled.MatchString("eth0") {
			t.Error("compiled exclude regex should match eth0")
		}
		if !includeInterfacesRegexCompiled.MatchString("swp1") {
			t.Error("compiled include regex should match swp1")
		}
	})

	t.Run("empty regexes are allowed", func(t *testing.T) {
		excludeInterfacesRegex = strptr("")
		includeInterfacesRegex = strptr("")
		if err := compileRegexFlags(); err != nil {
			t.Fatalf("compileRegexFlags() with empty regexes: %v", err)
		}
	})

	t.Run("invalid exclude regex errors", func(t *testing.T) {
		excludeInterfacesRegex = strptr("eth[")
		includeInterfacesRegex = strptr("")
		if err := compileRegexFlags(); err == nil {
			t.Error("compileRegexFlags() should error on invalid exclude regex")
		}
	})

	t.Run("invalid include regex errors", func(t *testing.T) {
		excludeInterfacesRegex = strptr("")
		includeInterfacesRegex = strptr("swp(")
		if err := compileRegexFlags(); err == nil {
			t.Error("compileRegexFlags() should error on invalid include regex")
		}
	})
}
