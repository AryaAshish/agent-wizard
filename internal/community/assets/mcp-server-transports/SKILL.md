# MCP server transports (stdio vs HTTP)

Pick and verify MCP transport: stdio vs HTTP—without mixing security models.

## When to use

- Shipping MCP consumed by Cursor/IDE or internal tools.
- Debugging empty tool lists or handshake failures.

## When not to use

- Non-MCP REST only.

## Inputs

- `YOUR_HOSTING`: desktop | server
- `YOUR_MCP_BINARY` or `YOUR_HTTP_BASE` (localhost URL)
- `YOUR_PORT` for HTTP health check

## Outputs

```
TRANSPORT_CHOICE: stdio|http

RATIONALE:
- one sentence

HANDSHAKE_CHECKLIST:
- [ ] process lifecycle (stdio) or server bind (http)
- [ ] tools enumerate non-empty
- [ ] stderr vs stdout rules respected (stdio)
- [ ] auth/TLS boundary for non-localhost

FINDINGS:
- bullet or "- none -"

BLOCKERS:
- bullet or "- none -"
```

## Steps

1. Stdio binary presence.

```bash
command -v YOUR_MCP_BINARY 2>/dev/null || test -x "./YOUR_MCP_BINARY" && echo "ok" || echo "missing YOUR_MCP_BINARY"
```

2. HTTP health (only if HTTP path).

```bash
curl -fsS "http://127.0.0.1:YOUR_PORT/health" 2>/dev/null || curl -fsS "YOUR_HTTP_BASE/health" 2>/dev/null || echo "health route missing or server down"
```

3. Logs must not corrupt stdio framing.

```bash
grep -RIn 'console\.log|print\(' -- YOUR_MCP_SRC_DIR 2>/dev/null | head -n 30 || true
```

## Stop and ask

Stop if neither `YOUR_MCP_BINARY` nor `YOUR_HTTP_BASE` is provided.

## Reject if

- `TRANSPORT_CHOICE: http` with bind `0.0.0.0` and no `BLOCKERS` row for missing auth.

## Safety

- Secrets via env—not CLI flags in shell history where avoidable.

## References

- Cross-check client transport expectations when upgrading MCP SDK.
