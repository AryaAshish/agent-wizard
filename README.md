# agent-wizard

CLI installs versioned agent skills into your repo—one command, no Slack drift.

- **Pain:** playbooks rot in chat and get re-pasted wrong per repo.
- **Do:** `add <skill> --source <library>` materializes `SKILL.md` under your tree.
- **Win:** same skill id everywhere; diff and review updates like code.

v0.1.3+: first `add` writes project config and runs **sync** (pass `--no-sync` to skip the copy step).

Commands walk **upward** from your shell directory for `agentskills.yaml`; if it’s missing, the first **`add`** creates it at the **nearest Git repo root** (when `.git` exists above you). When your cwd isn’t that directory, the CLI prints **`project: /absolute/path`** so you know where files landed.

**TTY:** run **`agent-wizard`** with no arguments, or **`agent-wizard wizard`** / **`guide`**, for a short guided menu (install a community skill or register a team Git URL).

## Use a community skill in this repo

1. `curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh`
2. `cd /path/to/your/repo`
3. `agent-wizard add pr-review --source community` → **`.agents/skills/pr-review/SKILL.md`** — e.g. `head -n1 .agents/skills/pr-review/SKILL.md` prints `# Pull request review`.
4. You get a real file in git; agents and editors load that path instead of fragile pasted prompts.
5. `agent-wizard list --source-name community` · [Bundled skills](docs/SKILLS.md)
6. npm: `npm i -g @aryaashish/agent-wizard` ([releases](https://github.com/AryaAshish/agent-wizard/releases)), then repeat step 3.
7. v0.1.2 or older: run `init`, then `add`, then `sync`, or upgrade.

## Use your team’s skill library

1. Put each skill as `skill-id/SKILL.md` in one Git repo; push (private GitHub is fine).
2. `agent-wizard sources add --name my-team --kind git --git-url https://github.com/your-org/my-team-skills.git`
3. In each app repo: `agent-wizard add deploy-checklist --source my-team` → **`.agents/skills/deploy-checklist/SKILL.md`**
4. Use the same workflow across all repos.
5. Definitions stay in one skills repo; every service pulls by id; ship updates through normal PRs.

`init` is optional (interactive picker / defaults)—not required for steps above on v0.1.3+.

---

## Detailed install options

### Why not just copy markdown?

| Copy-paste | agent-wizard |
|------------|----------------|
| Playbooks drift across Slack / Notion | Playbooks live **next to your code**; the project file lists what’s installed |
| Re-copy every new service | **Next repo:** same `add … --source …` line, same skill id |
| Everyone on different revisions | **`lock`** — teammates sync the **same** pinned revision |
| Rewire paths per agent | **Profiles** — one config, multiple install paths |

### curl installer (works everywhere we ship binaries)

Writes the binary into `$HOME/go/bin` by default (override with `INSTALL_DIR`):

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
```

Install a **specific GitHub release** (must exist with matching tarball + `checksums.txt`):

```bash
VERSION=v0.1.3 curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
```

### npm / npx

The published package is **`@aryaashish/agent-wizard`**. It downloads the matching release from GitHub Releases, verifies checksums, caches under `~/.cache/agent-wizard/npm`, and runs the real binary.

```bash
# one-off usage
npx @aryaashish/agent-wizard --version

# or global install
npm i -g @aryaashish/agent-wizard
agent-wizard --version
```

Until a release has been published to npm, install from the repo path or use the curl installer above.

### Homebrew (optional)

Homebrew support is **maintainer-provided**: CI can push a formula to a separate tap repository when `HOMEBREW_TAP_REPO` and `HOMEBREW_TAP_TOKEN` are configured. `brew tap …` must point at whatever **public** GitHub repo actually holds the formula (for example `your-user/homebrew-agent-wizard`), not a placeholder name:

```bash
brew tap <github-user>/<tap-repo>
brew install agent-wizard
```

If you do not have a tap configured yet, use **curl** or **npm** instead.

### Build from source (advanced, requires Go 1.25+)

```bash
go install github.com/aryaashish/agent-wizard@latest
```

Ensure `$HOME/go/bin` (or your `GOBIN`) is on `PATH`. Pin to a tagged release if you prefer reproducible installs (for example `@v0.1.3` once that tag exists).

**`go install` and `--version`:** Plain `go install` does not pass release ldflags, so **`agent-wizard --version`** may print `dev (commit=none date=unknown)`. That is expected; use a **release binary** (`install.sh`, GitHub Release asset, or npm’s downloaded binary) when you need semver + embedded build metadata.

**Windows** (from source):

```powershell
git clone https://github.com/AryaAshish/agent-wizard.git
cd agent-wizard
go build -o agent-wizard.exe .
```

Verify installation:

```bash
agent-wizard --version
```

If `agent-wizard --version` shows an older version, your shell might be using a stale binary path. Run `which agent-wizard`, open a new terminal, or run `hash -r` (bash/zsh) so `PATH` picks up `$HOME/go/bin`.

### Troubleshooting: `agent-skills-community` / `Repository not found`

Older setups stored the `community` source in **`~/.agent-wizard-config.yaml`** as a **git** remote pointing at `https://github.com/AryaAshish/agent-skills-community.git`, which no longer exists. Starter skills ship **inside the CLI** (embedded library), not from that repo.

**Fix:** upgrade `agent-wizard`, then run **`agent-wizard init`** (in any project) or **`agent-wizard community fetch`** so the global config is rewritten to the embedded `community` source. Or edit `~/.agent-wizard-config.yaml` and remove the `community` git entry (or run `agent-wizard sources remove community` and run `init` again so the default is re-added).

---

## Building your team library (detail)

Each skill is a folder with `SKILL.md`. Register once and `add` from each app repo as in **Use your team’s skill library** above.

**Packs:** add `.agent-wizard-pack.yaml` at the library root:

```yaml
name: onboarding-kit
skills:
  - code-review-guidelines
  - deploy-checklist
  - security-audit
```

Install the bundle in a project: `agent-wizard pack add onboarding-kit` (then `agent-wizard sync` if you passed `--no-sync` on prior adds).

---

## Lock versions and keep your team in sync

Pin the exact versions everyone should use:

```bash
agent-wizard lock
```

This creates `agentskills.lock` — commit it to your repo. Now when a teammate clones the project:

```bash
agent-wizard sync --strict-lock
```

Everyone gets the exact same skill versions. No surprises.

Check if anything has drifted:

```bash
agent-wizard status --check-drifts
```

---

## Use with different agents

**Cursor** — works out of the box, skills go to `.agents/skills/`.

**Claude Code** — change the target directory in `agentskills.yaml`:

```yaml
targetDir: .claude/skills
```

**Multiple agents at once** — use profiles:

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

## CI — enforce skills in your pipeline

Add to your CI script:

```bash
agent-wizard sync --strict-lock    # fail if lockfile doesn't match
agent-wizard status --check-drifts # exit code 3 if drift detected
agent-wizard ci-check              # validate policy gates
```

Set policy via environment variables:

```bash
export AGENT_WIZARD_ALLOWED_SOURCES="my-team,community"
export AGENT_WIZARD_MIN_SCHEMA_VERSION=1
```

---

## Source types

You can point `agent-wizard` at three kinds of skill sources:

| Type | Command | Best for |
|------|---------|----------|
| **Local folder** | `sources add --name dev --kind local --path ~/my-skills` | Developing skills locally |
| **Git repo** | `sources add --name team --kind git --git-url https://... [--git-ref main]` | Team/community repos |
| **Zip archive** | `sources add --name release --kind archive --archive-url https://...` | Pinned release snapshots |

Optional for local sources: add `--quiet` to suppress the collaboration warning in scripts.

`local` paths are machine-specific. They are great for your own development machine, but not team-shareable unless everyone mounts the same shared filesystem.

---

## All commands

| Command | What it does |
|---------|--------------|
| `wizard` / `guide` | Guided menu on an interactive terminal (install skill / add Git source) |
| `init` | Create `agentskills.yaml` in your project |
| `help <command>` | Show detailed help for a command |
| `add SKILL --source NAME` | Add a skill from a specific source |
| `add SKILL -NAME` | Shorthand source selector (for example `-android`) |
| `remove SKILL` | Remove a skill |
| `pack add PACK` | Add a skill bundle |
| `list --source-name NAME [--filter SUB]` | Browse skills (id + summary; sorted, aligned); optional id filter |
| `list --installed` | See what's installed (same id + summary columns) |
| `create-skill ID` | Create `<ID>/SKILL.md` from template (resolved next to manifest / git root like `add`) |
| `sync` | Copy skills into your project |
| `sync --dry-run` | Preview without writing |
| `sync --prune` | Remove skills not in manifest |
| `sync --strict-lock` | Fail if lockfile mismatch |
| `lock` | Pin current versions |
| `status` | Show project status |
| `status --json` | Status as JSON |
| `status --check-drifts` | Detect lockfile drift |
| `sources list` | Show configured sources |
| `sources add` | Register a new source |
| `sources remove` | Remove a source |
| `community fetch` | Refresh bundled community assets cache |
| `migrate` | Upgrade manifest schema |
| `cache status` | Show cache info |
| `cache prune` | Clear cached downloads |
| `ci-check` | Run CI policy checks |
| `browse` | Interactive skill picker |
| `watch` | Auto-sync on changes |
| `import --from DIR --into DIR` | Import existing skills |

---

## Documentation

- [Shipping test plan](docs/test-plan-ship.md)
- [Bundled skills index](docs/SKILLS.md)
- [Roadmap](ROADMAP.md)
- [Launch metrics log](docs/metrics.md)
- [Show HN draft](docs/show-hn.md)
- [Manifest schema](docs/spec/manifest-schema.md)
- [Lockfile schema](docs/spec/lockfile-schema.md)
- [CLI contract](docs/cli-contract.md)
- [Threat model](docs/security/threat-model.md)
- [Release checklist](docs/release/release-checklist.md)
- [Release notes template](docs/release/release-notes-template.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md).

## License

MIT
