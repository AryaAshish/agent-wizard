package config

import (
	"path/filepath"
	"testing"
)

func TestDefaultPathRespectsHomeWhenSet(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	got, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, FileName)
	if got != want {
		t.Fatalf("DefaultPath() = %q, want %q (so tests isolate config on Windows)", got, want)
	}
}
