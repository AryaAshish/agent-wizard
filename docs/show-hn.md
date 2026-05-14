# Show HN draft (agent-wizard)

## Title options

1. Stop re-pasting the same agent instructions into every new repo  
2. agent-wizard: repeatable agent skills in your tree—second repo, same commands  
3. From scattered prompts to checked-in playbooks your whole team can sync  

## Post body (paste & tweak)

Every new repo shouldn’t mean hunting Slack or Notion for the same PR checklist or launch gates. **agent-wizard** installs reusable skill folders (`SKILL.md` + assets) into your project—think reproducible playbooks next to your code. Works with Cursor, Claude Code, Codex (anything that reads skills from disk).

### Demo (~60s)

```bash
curl -fsSL https://raw.githubusercontent.com/AryaAshish/agent-wizard/main/install.sh | sh
mkdir -p /tmp/aw-demo-a && cd /tmp/aw-demo-a && git init -q
agent-wizard init
agent-wizard list --source-name community --filter pr
agent-wizard add pr-review --source community && agent-wizard sync
```

On **v0.1.3+**, a single `add pr-review --source community` is enough in a fresh repo; this block stays compatible with **v0.1.2** (explicit `init` + `sync`).

Repeat in `/tmp/aw-demo-b` with the same commands—same playbook, zero forwarding threads.

Repo: https://github.com/AryaAshish/agent-wizard  

### Compared to manual paste

| Paste into chat | agent-wizard |
|-----------------|--------------|
| Drifts per repo | Manifest + optional lockfile |
| Second repo = scramble | Same skill ids & commands |

### Limitations

- Starter library ships embedded—grow via git/archive sources or PRs (`skill-request` template).
- Security posture: see [`docs/security/threat-model.md`](docs/security/threat-model.md)—skills run as untrusted markdown executed by **your** agent policies.

### First comment suggestion

Link threat model + invite skill requests via GitHub Issues template.
