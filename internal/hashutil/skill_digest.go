package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// SkillMarkdownDigest returns sha256 hex of SKILL.md for a skill directory.
func SkillMarkdownDigest(skillDir string) (string, error) {
	p := filepath.Join(skillDir, "SKILL.md")
	b, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("read SKILL.md: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
