package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aryaashish/agent-wizard/internal/archive"
	"github.com/aryaashish/agent-wizard/internal/cache"
	"github.com/aryaashish/agent-wizard/internal/community"
	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/skills"
	"github.com/aryaashish/agent-wizard/internal/vcs"
)

type materializedSource struct {
	Name        string
	Root        string
	ResolvedRef string
}

func materializeSource(cfg config.Source) (materializedSource, error) {
	name := cfg.Name
	switch cfg.Kind {
	case "local":
		if cfg.Path == "" {
			return materializedSource{}, fmt.Errorf("local source %q missing path", name)
		}
		abs, err := filepath.Abs(cfg.Path)
		if err != nil {
			return materializedSource{}, err
		}
		return materializedSource{Name: name, Root: abs, ResolvedRef: "local"}, nil
	case "git":
		if cfg.GitURL == "" {
			return materializedSource{}, fmt.Errorf("git source %q missing gitUrl", name)
		}
		dest, err := cache.GitCheckoutDir(cfg.GitURL)
		if err != nil {
			return materializedSource{}, err
		}
		root, head, err := vcs.EnsureCheckout(dest, cfg.GitURL, cfg.GitRef, cfg.Subdir)
		if err != nil {
			return materializedSource{}, err
		}
		return materializedSource{Name: name, Root: root, ResolvedRef: head}, nil
	case "archive":
		if cfg.ArchiveURL == "" {
			return materializedSource{}, fmt.Errorf("archive source %q missing archiveUrl", name)
		}
		dest, err := cache.ArchiveExtractDir(cfg.ArchiveURL)
		if err != nil {
			return materializedSource{}, err
		}
		if err := os.RemoveAll(dest); err != nil {
			return materializedSource{}, err
		}
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return materializedSource{}, err
		}
		if err := archive.ExtractRemoteZip(cfg.ArchiveURL, dest, 0); err != nil {
			return materializedSource{}, err
		}
		return materializedSource{Name: name, Root: dest, ResolvedRef: "archive:" + cfg.ArchiveURL}, nil
	case community.SourceKind:
		root, err := community.Extract(false)
		if err != nil {
			return materializedSource{}, err
		}
		return materializedSource{Name: name, Root: root, ResolvedRef: "embedded"}, nil
	default:
		return materializedSource{}, fmt.Errorf("unsupported source kind %q for %q", cfg.Kind, name)
	}
}

// MaterializeSource exposes materialization rules for ancillary CLI workflows (listing without sync).
func MaterializeSource(cfg config.Source) (root string, resolvedRef string, err error) {
	ms, err := materializeSource(cfg)
	if err != nil {
		return "", "", err
	}
	return ms.Root, ms.ResolvedRef, nil
}

func discoverFromMaterialized(ms materializedSource) ([]skills.Skill, error) {
	src := skills.NewLocalPathSource(ms.Root)
	found, err := src.Discover()
	if err != nil {
		return nil, err
	}
	for i := range found {
		found[i].SourceName = ms.Name
		found[i].ResolvedRef = ms.ResolvedRef
	}
	return found, nil
}
