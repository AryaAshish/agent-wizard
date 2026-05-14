# agent-wizard

**What:** CLI that installs reusable **agent skills** (folders with `SKILL.md`) into your repo—like package management for agent playbooks.

**Who:** Developers using **Cursor, Claude Code, Codex**, or any workflow that reads skill files from disk.

**Outcome:** Skills land under your chosen directory in **about a minute**, ready to commit.

**Install:** `curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh` — or `npm i -g @aryaashish/agent-wizard` ([releases](https://github.com/AryaAshish/agent-wizard/releases)).

**Why not only copy-paste:** The **next repo** repeats with **one command** (`add` initializes the manifest and syncs)—no hunting Slack or Notion. Pin versions with `lock` / `sync --strict-lock` when your team needs one truth.

### Why not just copy markdown?

| Copy-paste | agent-wizard |
|------------|----------------|
| Playbooks drift across Slack / Notion | Playbooks live **next to your code**; `agentskills.yaml` records what this repo installs |
| Re-copy every new service | **Second project:** same `add … --source …` line, same skill id |
| Everyone on different revisions | **`lock`** — teammates sync the **same** pinned revision |
| Rewire paths per agent | **`targetDir` / profiles** — one manifest, multiple agent layouts |

## Happy path (first run)

With no manifest yet, `add` creates `agentskills.yaml`, wires the bundled **community** source, and runs **sync** (use `--no-sync` if you only want the manifest updated).

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
hash -r 2>/dev/null || true
agent-wizard --version
mkdir -p /tmp/my-agent-demo && cd /tmp/my-agent-demo && git init -q
agent-wizard add pr-review --source community
# Optional: browse first — agent-wizard list --source-name community --filter pr
test -f .agents/skills/pr-review/SKILL.md && echo "Skill on disk — open it in your editor."
```

Interactive **`init`** is still useful for the starter picker and browsing before you choose skills. **`add`** does not open the picker. **[Bundled skill index →](docs/SKILLS.md)**

---

## Step 1 — Install

### curl installer (works everywhere we ship binaries)

Writes the binary into `$HOME/go/bin` by default (override with `INSTALL_DIR`):

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
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

### Build from source (advanced, requires Go 1.22+)

```bash
go install github.com/aryaashish/agent-wizard@latest
```

Ensure `$HOME/go/bin` (or your `GOBIN`) is on `PATH`. Pin to a tagged release if you prefer reproducible installs (for example `@v0.1.0` once that tag exists).

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

## Step 2 — Set up your project

Open your project in Cursor (or any editor) and run:

```bash
agent-wizard init
```

This creates an `agentskills.yaml` file, auto-attaches the bundled `community` source, and starts an interactive starter picker.

---

## Step 3 — Explore and install community skills

Browse available skills from bundled community source:

```bash
# See what's available
agent-wizard list --source-name community

# Add individual skills from that source
agent-wizard add pr-review --source community
agent-wizard add plan-review --source community

# Or add a whole pack (a bundle of related skills)
agent-wizard pack add android-starter
```

Sync the skills into your project:

```bash
agent-wizard sync
```

That's it. The skills are now in `.agents/skills/` and your AI agent can use them immediately.

---

## Step 4 — Create and share your own skills (private to your team)

### Create a skill

A skill is just a folder with a `SKILL.md` file:

```
my-team-skills/
  deploy-checklist/
    SKILL.md
  code-review-guidelines/
    SKILL.md
  security-audit/
    SKILL.md
```

Push this to a **private** Git repo that your team has access to.

### Share with your team

Every team member runs these two commands once:

```bash
# Register your team's private skill library (shareable)
agent-wizard sources add --name my-team --kind git \
  --git-url https://github.com/your-org/my-team-skills.git

# See all team skills
agent-wizard list --source-name my-team
```

Then in any project:

```bash
agent-wizard add deploy-checklist --source my-team
# shorthand:
# agent-wizard add deploy-checklist -my-team
agent-wizard sync
```

Your skills stay in your private repo. No one outside your org can see them. Every teammate with repo access can install them in one command.

### Packs — bundle multiple skills together

Create a file called `.agent-wizard-pack.yaml` in your library:

```yaml
name: onboarding-kit
skills:
  - code-review-guidelines
  - deploy-checklist
  - security-audit
```

Now any teammate can install the whole bundle:

```bash
agent-wizard pack add onboarding-kit
agent-wizard sync
```

---

## Step 5 — Lock versions and keep your team in sync

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
| `init` | Create `agentskills.yaml` in your project |
| `help <command>` | Show detailed help for a command |
| `add SKILL --source NAME` | Add a skill from a specific source |
| `add SKILL -NAME` | Shorthand source selector (for example `-android`) |
| `remove SKILL` | Remove a skill |
| `pack add PACK` | Add a skill bundle |
| `list --source-name NAME [--filter SUB]` | Browse skills in a source (optional id filter) |
| `list --installed` | See what's installed |
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
