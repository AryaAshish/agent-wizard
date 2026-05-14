# Template: minimal bootstrap

1. `go install github.com/aryaashish/agent-wizard@latest` (or install from source; see README for **v0.1.3+** vs older tags).
2. In your app repo: **`agent-wizard add pr-review --source community`** (v0.1.3+ creates manifest + sync), **or** on older CLIs: `agent-wizard init` then `add` then `sync`.
3. Add a library source: `agent-wizard sources add --name lib --kind local --path /path/to/library`
4. `agent-wizard pack add android-starter` (optional) or `agent-wizard add <other-skill> --source lib`
5. `agent-wizard lock`
6. `agent-wizard sync` (idempotent if `add` already synced)
