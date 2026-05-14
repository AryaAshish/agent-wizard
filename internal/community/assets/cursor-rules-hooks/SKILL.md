# Cursor rules, hooks, and skills layout

Organize `.cursor/rules`, hooks, and synced skills so agents load stable guidance.

## When to use

- Introducing Cursor assets in a repo with `agentskills.yaml` / `agent-wizard`.

## When not to use

- No repo path; experimentation only without canonical folder choice.

## Inputs

- `YOUR_REPO_ROOT`
- Canonical skills dir per manifest: `YOUR_TARGET_DIR` (from `agentskills.yaml` if present)

## Outputs

```
TREE:
- path | purpose | consumption order index

MANIFEST:
- agentskills.yaml present: yes|no
- targetDir: YOUR_TARGET_DIR or unknown

HOOKS:
- .cursor/hooks.json status: present|absent
- scripts: list or "- none -"

RISKS:
- bullet or "- none -"

NEXT_COMMANDS:
- bullet list (max 5)
```

## Steps

1. Rules directory.

```bash
cd YOUR_REPO_ROOT
mkdir -p .cursor/rules
find .cursor/rules -maxdepth 1 -type f 2>/dev/null | head -n 20
```

2. Manifest and sync tool.

```bash
grep -nE 'targetDir|profiles|skills' agentskills.yaml 2>/dev/null | head -n 40 || echo "no agentskills.yaml"
agent-wizard status --json 2>/dev/null | head -c 2000 || echo "agent-wizard not runnable here"
```

3. Hooks file.

```bash
test -f .cursor/hooks.json && wc -c .cursor/hooks.json || echo "no hooks.json"
```

## Stop and ask

Stop if `YOUR_REPO_ROOT` is not the repo root (no `.git` and no manifest).

## Reject if

- Recommends remote hook `curl | sh` without hash pin and owner review in `RISKS`.

## Safety

- Hooks run arbitrary shell—diff review mandatory; keep scripts idempotent and fast.

## References

- Use manifest `profiles` for multi-client paths per project README.
