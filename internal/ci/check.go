package ci

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aryaashish/agent-wizard/internal/manifest"
)

// Check validates optional environment policy hooks for CI.
func Check(m manifest.Manifest) error {
	if raw := strings.TrimSpace(os.Getenv("AGENT_WIZARD_ALLOWED_SOURCES")); raw != "" {
		allowed := map[string]struct{}{}
		for _, part := range strings.Split(raw, ",") {
			name := strings.TrimSpace(part)
			if name == "" {
				continue
			}
			allowed[name] = struct{}{}
		}
		for _, s := range m.Sources {
			if _, ok := allowed[s]; !ok {
				return fmt.Errorf("source %q not in AGENT_WIZARD_ALLOWED_SOURCES allowlist", s)
			}
		}
	}

	if raw := strings.TrimSpace(os.Getenv("AGENT_WIZARD_MIN_SCHEMA_VERSION")); raw != "" {
		minV, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid AGENT_WIZARD_MIN_SCHEMA_VERSION: %w", err)
		}
		if m.SchemaVersion < minV {
			return fmt.Errorf("manifest schemaVersion %d below required %d", m.SchemaVersion, minV)
		}
	}
	return nil
}
