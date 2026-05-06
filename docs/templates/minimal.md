# Template: minimal bootstrap

1. `go install github.com/aryaashish/agent-wizard@latest` (or install from source)
2. In your app repo: `agent-wizard init`
3. Add a library source: `agent-wizard sources add --name lib --kind local --path /path/to/library`
4. `agent-wizard pack add android-starter` (optional) or `agent-wizard add pr-review`
5. `agent-wizard lock`
6. `agent-wizard sync`
