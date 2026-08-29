package main

import (
	"testing"
	"time"
)

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

func TestNewServerConfiguresHTTPTimeouts(t *testing.T) {
	originalListenAddress := listenAddress
	listenAddress = strptr("127.0.0.1:0")
	t.Cleanup(func() {
		listenAddress = originalListenAddress
	})

	server := newServer()
	if got, want := server.Addr, "127.0.0.1:0"; got != want {
		t.Errorf("server address = %q, want %q", got, want)
	}
	if got, want := server.ReadHeaderTimeout, 5*time.Second; got != want {
		t.Errorf("server ReadHeaderTimeout = %s, want %s", got, want)
	}
	if got, want := server.ReadTimeout, 10*time.Second; got != want {
		t.Errorf("server ReadTimeout = %s, want %s", got, want)
	}
	if got, want := server.IdleTimeout, time.Minute; got != want {
		t.Errorf("server IdleTimeout = %s, want %s", got, want)
	}
}
