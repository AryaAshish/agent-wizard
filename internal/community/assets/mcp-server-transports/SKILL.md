# MCP server transports (stdio vs HTTP)

Pick and verify a Model Context Protocol server transport suitable for local dev and production: stdio vs Streamable HTTP—without mixing security models.

## When to use

- Shipping an MCP server consumed by Cursor/Claude Code or internal tools.
- Debugging “client can’t see tools” failures.

## When not to use

- One-off REST APIs without MCP framing—don’t force MCP semantics.

## Inputs

- Host environment: local IDE (stdio typical) vs remote multi-user (HTTP with auth boundary).
- Required tool surface area and latency envelope.

## Outputs

- Explicit transport choice + checklist of verified handshake steps.

## Steps

1. **Stdio:** simplest for desktop—one process lifecycle bound to client spawn; avoid daemonizing unintentionally.

```bash
# Example: run server binary directly — replace with your entrypoint
which YOUR_MCP_BINARY || echo "Build your MCP server first."
```

2. **HTTP/SSE variants:** prefer documented Streamable HTTP patterns from current MCP spec guidance—terminate TLS at gateway when exposed beyond localhost.

```bash
curl -fsS http://127.0.0.1:YOUR_PORT/health || echo "Expose minimal health route."
```

3. Validate tools enumerate—empty tool lists usually signal startup crash before handshake completes.

4. Logging: stderr for diagnostics on stdio servers—avoid corrupting stdout JSON-RPC framing.

## Safety

- Binding MCP HTTP servers `0.0.0.0` without auth exposes arbitrary tool execution—default localhost + proxy auth.

- Secrets via env—not CLI flags captured in shell history where avoidable.

## References

- Cross-check transport expectations against your client docs when upgrading MCP SDK versions.
