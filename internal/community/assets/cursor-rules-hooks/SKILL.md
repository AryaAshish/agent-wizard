# Cursor rules, hooks, and skills layout

Organize Cursor-facing assets so agents pick up stable guidance: `.cursor/rules`, hooks JSON, and synced skill folders (`agent-wizard` targets).

## When to use

- Introducing Cursor on a repo alongside other agents (Claude Code paths differ—manifest profiles cover both).

## When not to use

- Solo experimentation—still helpful but prioritize one canonical folder layout before proliferating duplicates.

## Inputs

- Whether `.agents/skills` or `.cursor/skills` is canonical for Cursor on this repo (`agentskills.yaml` `targetDir`).
- Hooks automation appetite (validate manifest before push vs advisory).

## Outputs

- Recommended directory tree + one sentence per artifact explaining consumption order.

## Steps

1. Place durable repo-wide guidance under `.cursor/rules/` using small scoped markdown files—not one mega-file nobody reads.

```bash
mkdir -p .cursor/rules
ls -la .cursor/rules 2>/dev/null || echo "Create rules files here."
```

2. Sync agent skills via `agent-wizard` into the folder Cursor indexes—avoid editing generated copies manually.

```bash
agent-wizard status --json 2>/dev/null || echo "Run init in repo root first."
grep -n targetDir agentskills.yaml 2>/dev/null || true
```

3. Hooks live under `.cursor/hooks.json` or team-standard location—keep scripts idempotent and fast (< few seconds).

```bash
test -f .cursor/hooks.json && cat .cursor/hooks.json || echo "No hooks.json yet."
```

4. Document precedence: manifest-declared skills vs rule files vs chat instructions—prefer manifest for shared truth.

## Safety

- Hooks running arbitrary shell—review diff carefully; never fetch remote scripts without pinning hashes.

## References

- Use manifest `profiles` to mirror paths for Claude/Cursor simultaneously—see README profiles snippet.
