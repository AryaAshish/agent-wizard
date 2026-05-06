package engine

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/config"
)

func TestMaterializeSource_GitLocalRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	run(t, repo, "git", "init")
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, repo, "git", "add", ".")
	run(t, repo, "git", "commit", "-m", "init")

	cfg := config.Source{
		Name:   "g",
		Kind:   "git",
		GitURL: repo,
	}
	root, ref, err := MaterializeSource(cfg)
	if err != nil {
		t.Fatalf("MaterializeSource git err=%v", err)
	}
	if root == "" || ref == "" {
		t.Fatalf("root/ref empty root=%q ref=%q", root, ref)
	}
	if _, err := os.Stat(filepath.Join(root, "README.md")); err != nil {
		t.Fatalf("git materialized file missing: %v", err)
	}
}

func TestMaterializeSource_ArchiveHTTPZip(t *testing.T) {
	zipBytes := makeZip(t, map[string]string{
		"pr-review/SKILL.md": "# pr",
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(zipBytes)
	}))
	defer srv.Close()

	cfg := config.Source{Name: "a", Kind: "archive", ArchiveURL: srv.URL}
	root, ref, err := MaterializeSource(cfg)
	if err != nil {
		t.Fatalf("MaterializeSource archive err=%v", err)
	}
	if root == "" || ref == "" {
		t.Fatalf("root/ref empty")
	}
	if _, err := os.Stat(filepath.Join(root, "pr-review", "SKILL.md")); err != nil {
		t.Fatalf("archive extracted file missing: %v", err)
	}
}

func makeZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range files {
		f, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func run(t *testing.T, dir string, cmd string, args ...string) {
	t.Helper()
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=agent-wizard-test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=agent-wizard-test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", cmd, args, err, string(out))
	}
}
