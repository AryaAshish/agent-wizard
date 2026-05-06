# agent-wizard

Manage AI agent skills across all your projects. Install community skills, create your own, share them privately with your team — all from one CLI.

Works with **Cursor, Claude Code, Codex**, and any agent that reads skill files.

---

## Step 1 — Install

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
```

**Or build from source** (requires Go 1.22+):

```bash
go install github.com/AryaAshish/agent-wizard@latest
```

**Windows** (from source):

```powershell
git clone https://github.com/AryaAshish/agent-wizard.git
cd agent-wizard
go build -o agent-wizard.exe .
```

Make sure the binary is on your `PATH`, then verify:

```bash
agent-wizard help
```

---

## Step 2 — Set up your project

Open your project in Cursor (or any editor) and run:

```bash
agent-wizard init
```

This creates an `agentskills.yaml` file in your project root. That file tracks which skills your project uses.

---

## Step 3 — Explore and install community skills

Browse available skills from any public skill library:

```bash
# Point at a community library
agent-wizard sources add --name community --kind git \
  --gitUrl https://github.com/AryaAshish/agent-skills-community.git

# See what's available
agent-wizard list --source-name community

# Add individual skills
agent-wizard add pr-review
agent-wizard add plan-review

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
# Register your team's private skill library
agent-wizard sources add --name my-team --kind git \
  --gitUrl https://github.com/your-org/my-team-skills.git

# See all team skills
agent-wizard list --source-name my-team
```

Then in any project:

```bash
agent-wizard add my-team/deploy-checklist
agent-wizard sync
```

Your skills stay in your private repo. No one outside your org can see them. But every teammate with repo access can install them in one command.

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
| **Git repo** | `sources add --name team --kind git --gitUrl https://...` | Team & community libraries |
| **Zip archive** | `sources add --name release --kind archive --archiveUrl https://...` | Pinned release snapshots |

---

## All commands

| Command | What it does |
|---------|--------------|
| `init` | Create `agentskills.yaml` in your project |
| `add SKILL` | Add a skill to your project |
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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md).

## License

MIT
