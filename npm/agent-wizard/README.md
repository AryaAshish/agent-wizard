# agent-wizard (npm wrapper)

This package installs a thin Node.js wrapper that downloads the platform-specific `agent-wizard` binary from GitHub Releases and executes it.

## Usage

```bash
npx @aryaashish/agent-wizard --version
```

or

```bash
npm i -g @aryaashish/agent-wizard
agent-wizard --version
```

## Environment overrides

- `AGENT_WIZARD_VERSION`: pin release version (default: `latest`)
- `AGENT_WIZARD_REPO`: override GitHub repo (default: `AryaAshish/agent-wizard`)
- `AGENT_WIZARD_CACHE_DIR`: override local binary cache location
