/*
Copyright © 2026 NAME HERE jhettenh@gmail.com
*/
package cmd

import (
	"bytes"
	"testing"

	"github.com/borland502/wordgen/internal/version"
)

func TestVersionCommand(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"version"})

	output := &bytes.Buffer{}
	cmd.SetOut(output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	result := output.String()
	if result == "" {
		t.Fatal("version command produced no output")
	}

	// Verify expected format
	if !bytes.Contains(output.Bytes(), []byte("wordgen version")) {
		t.Errorf("output missing 'wordgen version': %s", result)
	}

	if !bytes.Contains(output.Bytes(), []byte("commit")) {
		t.Errorf("output missing 'commit': %s", result)
	}

	if !bytes.Contains(output.Bytes(), []byte("built")) {
		t.Errorf("output missing 'built': %s", result)
	}

	// Verify the actual version is present
	if !bytes.Contains(output.Bytes(), []byte(version.Version)) {
		t.Errorf("output missing version %q: %s", version.Version, result)
	}
}
