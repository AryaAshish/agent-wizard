# agent-wizard

Manage AI agent skills across all your projects. Install community skills, create your own, share them privately with your team — all from one CLI.

Works with **Cursor, Claude Code, Codex**, and any agent that reads skill files.

---

## 30-second quick start

```bash
agent-wizard init
agent-wizard list --source-name community
```

`init` now auto-wires the bundled starter community library and opens an interactive picker.
If you skip the picker, install one skill manually:

```bash
agent-wizard add pr-review --source community
agent-wizard sync
```

After `init`, the CLI also prints suggested commands:

```bash
agent-wizard list --source-name community
agent-wizard add pr-review --source community
agent-wizard pack add android-starter
agent-wizard sync
```

---

## Step 1 — Install

### Homebrew (recommended on macOS/Linux)

```bash
brew tap aryaashish/tap
brew install agent-wizard
```

### npm / npx

```bash
# one-off usage
npx agent-wizard --version

# or global install
npm i -g agent-wizard
agent-wizard --version
```

### curl installer

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
```

### Build from source (advanced, requires Go 1.22+)

```bash
go install github.com/aryaashish/agent-wizard@main
```

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

If `agent-wizard --version` shows an older version, your shell might be using a stale binary path. Run `which agent-wizard` and ensure it points to your intended install location.

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
agent-wizard add plan-review -community

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
| `list --source-name NAME` | Browse skills in a source |
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
