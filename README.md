# agent-wizard

A **local-first CLI** for managing reusable **agent skills** across projects. One tool to discover, pin, sync, and share skills for Cursor, Claude Code, Codex, and any future agent — without a mandatory company-wide skills server.

---

## Install

### macOS (Homebrew — coming soon)

```bash
# From source (requires Go 1.22+)
go install github.com/AryaAshish/agent-wizard@latest
```

### macOS / Linux (install script)

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
```

This installs the binary into `$GOBIN` (default `~/go/bin`). Make sure it is on your `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

### Linux (from source)

```bash
git clone https://github.com/AryaAshish/agent-wizard.git
cd agent-wizard
go build -o agent-wizard .
sudo mv agent-wizard /usr/local/bin/
```

### Windows (from source)

```powershell
git clone https://github.com/AryaAshish/agent-wizard.git
cd agent-wizard
go build -o agent-wizard.exe .
# Move agent-wizard.exe to a directory on your PATH
```

### Verify installation

```bash
agent-wizard help
```

---

## Quickstart: new project setup

### 1. Initialize the manifest

Open your project folder in your terminal (e.g. in Cursor's integrated terminal):

```bash
cd /path/to/my-app
agent-wizard init
```

This creates `agentskills.yaml` with sensible defaults (`manifest-only` mode, target `.agents/skills`).

### 2. Register a skill source

Point the CLI at a skill library — a local folder, a Git repo, or an archive URL.

**Local path** (fastest for development):

```bash
agent-wizard sources add --name my-lib --kind local --path ~/skills-library
```

**Git remote** (team library):

```bash
agent-wizard sources add --name team --kind git \
  --path . \
  --gitUrl https://github.com/your-org/agent-skills.git
```

**Archive URL** (immutable release snapshot):

```bash
agent-wizard sources add --name release --kind archive \
  --path . \
  --archiveUrl https://example.com/skills-v1.0.0.zip
```

List configured sources:

```bash
agent-wizard sources list
```

### 3. Wire sources into your project manifest

Edit `agentskills.yaml` and add source names under `sources:`:

```yaml
schemaVersion: 1
targetDir: .agents/skills
installMode: manifest-only
sources:
  - my-lib
skills: []
packs: []
```

### 4. Add skills or packs

```bash
# Add a single skill
agent-wizard add pr-review

# Or add an entire pack (a bundle of related skills)
agent-wizard pack add android-starter
```

### 5. Pin versions with a lockfile

```bash
agent-wizard lock
```

This writes `agentskills.lock` with resolved refs and content digests.

### 6. Sync skills into your project

```bash
# Preview what will happen
agent-wizard sync --dry-run

# Actually copy skills into .agents/skills/
agent-wizard sync
```

### 7. Check status

```bash
# Human-readable
agent-wizard status

# Machine-readable JSON (for CI)
agent-wizard status --json

# Detect drift against lockfile (exit code 3 on drift)
agent-wizard status --check-drifts
```

---

## Using with Cursor IDE

```bash
# In Cursor's integrated terminal:
cd /path/to/my-project
agent-wizard init
agent-wizard sources add --name lib --kind local --path /path/to/skills
agent-wizard add pr-review
agent-wizard sync

# Skills are now in .agents/skills/ — Cursor reads them automatically
```

## Using with Claude Code

```bash
# Configure target dir for Claude's expected path
# Edit agentskills.yaml:
#   targetDir: .claude/skills

agent-wizard init
# ... same flow as above ...
agent-wizard sync
```

## Using with any agent tool

Edit `agentskills.yaml` and set `targetDir` to wherever your agent reads skills from, or use **profiles** for multi-agent output:

```yaml
profiles:
  - name: cursor
    targets:
      - kind: agents
        path: .agents/skills
  - name: claude
    targets:
      - kind: agents
        path: .claude/skills
```

---

## All commands

| Command | Description |
|---------|-------------|
| `agent-wizard help` | Show help |
| `agent-wizard init` | Create `agentskills.yaml` in current directory |
| `agent-wizard sources list\|add\|remove` | Manage source registry (`local`, `git`, `archive`) |
| `agent-wizard list --source DIR` | List skills available in a source directory |
| `agent-wizard list --source-name NAME` | List skills from a configured source |
| `agent-wizard list --installed` | List skills selected by current manifest |
| `agent-wizard add SKILL` | Add a skill to the manifest |
| `agent-wizard remove SKILL` | Remove a skill from the manifest |
| `agent-wizard pack add PACK` | Add a pack (bundle of skills) to the manifest |
| `agent-wizard lock` | Write `agentskills.lock` with pinned versions/digests |
| `agent-wizard sync` | Copy selected skills into target directories |
| `agent-wizard sync --dry-run` | Preview sync without writing files |
| `agent-wizard sync --prune` | Remove skill dirs not in manifest (destructive) |
| `agent-wizard sync --strict-lock` | Fail if skills don't match lockfile pins |
| `agent-wizard status` | Show manifest and source status |
| `agent-wizard status --json` | Emit status as JSON |
| `agent-wizard status --check-drifts` | Compare lockfile vs live sources (exit 3 on drift) |
| `agent-wizard migrate` | Backup manifest and normalize to current schema |
| `agent-wizard cache status\|prune` | Inspect or clear local cache |
| `agent-wizard ci-check` | Validate env policy gates for CI |
| `agent-wizard catalog validate FILE` | Validate a curated index YAML file |
| `agent-wizard import --from DIR --into DIR` | Import existing skill trees into a library |
| `agent-wizard browse [DIR]` | Interactive skill picker (stdin) |
| `agent-wizard watch` | Poll-based auto-sync loop |
| `agent-wizard icp MODE` | Validate ICP mode (`solo\|team\|enterprise`) |

---

## Source kinds

| Kind | Config fields | Notes |
|------|---------------|-------|
| `local` | `path` | Fastest for local development and air-gapped setups |
| `git` | `gitUrl`, `gitRef`, `subdir` | Clones/updates to XDG cache; pins by commit SHA |
| `archive` | `archiveUrl` | Downloads zip with safe extraction guards |

---

## Team workflow

1. **Author** publishes a skill (folder with `SKILL.md`) to the team's shared library (Git repo or shared path).
2. **Each project** declares which skills it uses in `agentskills.yaml`.
3. **Teammates** clone the repo and run `agent-wizard sync` — everyone gets the same skills.
4. **Updates** flow through the library; teams run `agent-wizard sync` or CI detects drift automatically.

---

## CI integration

```bash
# In your CI pipeline:
agent-wizard sync --strict-lock    # fail if lockfile doesn't match
agent-wizard status --check-drifts # exit 3 if drift detected

# Optional env-based policy gates:
export AGENT_WIZARD_ALLOWED_SOURCES="team-lib,community"
export AGENT_WIZARD_MIN_SCHEMA_VERSION=1
agent-wizard ci-check
```

---

## Skill format

Any directory containing a `SKILL.md` file is treated as a skill. The directory name becomes the skill ID.

```
my-library/
  pr-review/
    SKILL.md        # required
    examples/       # optional supporting files
  plan-review/
    SKILL.md
```

---

## Documentation

- [Manifest schema](docs/spec/manifest-schema.md)
- [Lockfile schema](docs/spec/lockfile-schema.md)
- [Resolver precedence](docs/spec/resolver-precedence.md)
- [Curated index schema](docs/spec/curated-index-schema.md)
- [CLI contract](docs/cli-contract.md)
- [Threat model](docs/security/threat-model.md)
- [Trust model](docs/security/trust-model.md)
- [Support matrix](docs/support-matrix.md)
- [Privacy](docs/privacy.md)
- [Compatibility](docs/compat.md)
- [Release checklist](docs/release/release-checklist.md)
- [Rollback runbook](docs/release/rollback-runbook.md)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## License

MIT
