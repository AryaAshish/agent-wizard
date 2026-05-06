package archive

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractZipRejectTraversal(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, pair := range []struct {
		Name string
		Data string
	}{
		{`../evil.txt`, "boom"},
	} {
		f, err := zw.Create(pair.Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write([]byte(pair.Data))
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	br := bytes.NewReader(buf.Bytes())
	zr, err := zip.NewReader(br, int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	dest := t.TempDir()
	if err := ExtractZip(zr, dest); err == nil {
		t.Fatal("expected error for traversal path")
	}
}

func TestExtractZipAllowsSafeNested(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	dir := filepath.ToSlash(filepath.Join("a", "b"))
	if _, err := zw.Create(dir + "/"); err != nil {
		t.Fatal(err)
	}
	f, err := zw.Create(dir + "/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("# ok")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	br := bytes.NewReader(buf.Bytes())
	zr, err := zip.NewReader(br, int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	if err := ExtractZip(zr, dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, filepath.FromSlash("a/b/SKILL.md"))); err != nil {
		t.Fatal(err)
	}
}
