package engine

import (
	"testing"

	"github.com/aryaashish/agent-wizard/internal/community"
	"github.com/aryaashish/agent-wizard/internal/config"
)

func TestIsLegacyCommunityGitSource(t *testing.T) {
	t.Parallel()

	legacy := config.Source{
		Name: community.SourceName,
		Kind: "git",
		GitURL: "https://github.com/AryaAshish/agent-skills-community.git",
	}
	if !IsLegacyCommunityGitSource(legacy) {
		t.Fatal("expected legacy match for canonical URL casing")
	}
	withSlash := legacy
	withSlash.GitURL = legacy.GitURL + "/"
	if !IsLegacyCommunityGitSource(withSlash) {
		t.Fatal("expected legacy match with trailing slash")
	}
	withSpace := legacy
	withSpace.GitURL = "  " + legacy.GitURL + "  "
	if !IsLegacyCommunityGitSource(withSpace) {
		t.Fatal("expected legacy match with surrounding space")
	}

	otherFork := legacy
	otherFork.GitURL = "https://github.com/other/agent-skills-community.git"
	if IsLegacyCommunityGitSource(otherFork) {
		t.Fatal("different owner should not be legacy")
	}
	wrongName := legacy
	wrongName.Name = "not-community"
	if IsLegacyCommunityGitSource(wrongName) {
		t.Fatal("wrong name should not be legacy")
	}
	local := config.Source{Name: community.SourceName, Kind: "local", Path: "/tmp"}
	if IsLegacyCommunityGitSource(local) {
		t.Fatal("local community should not be legacy git")
	}
}

func TestResolveSource_LegacyCommunityGit(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		SchemaVersion: 1,
		Sources: []config.Source{{
			Name:   community.SourceName,
			Kind:   "git",
			GitURL: "https://github.com/AryaAshish/agent-skills-community.git",
		}},
	}
	src, ok := ResolveSource(cfg, community.SourceName)
	if !ok {
		t.Fatal("expected source")
	}
	if src.Kind != community.SourceKind || src.Name != community.SourceName {
		t.Fatalf("got %+v, want embedded community kind", src)
	}
	if src.GitURL != "" {
		t.Fatalf("expected empty gitUrl, got %q", src.GitURL)
	}
}

func TestResolveSource_NonLegacyCommunityGitUnchanged(t *testing.T) {
	t.Parallel()

	want := config.Source{
		Name:   community.SourceName,
		Kind:   "git",
		GitURL: "https://github.com/AryaAshish/my-fork-skills.git",
	}
	cfg := config.Config{SchemaVersion: 1, Sources: []config.Source{want}}
	src, ok := ResolveSource(cfg, community.SourceName)
	if !ok {
		t.Fatal("expected source")
	}
	if src.Kind != want.Kind || src.GitURL != want.GitURL {
		t.Fatalf("got %+v, want %+v unchanged", src, want)
	}
}
