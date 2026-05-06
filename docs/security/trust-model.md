# Trust Model (v1 baseline)

## Trust tiers

1. **Local trusted**: local filesystem sources configured by the user.
2. **Remote untrusted**: remote source kinds represented but disabled by default.
3. **Future trusted remote**: allowlisted hosts + signed artifacts (roadmap).

## v1 behavior

- Only `local` sources are executable in resolver/sync.
- Ambiguous skills across sources fail fast.
- `sync --dry-run` is available for non-mutating inspection.

## Roadmap

- source allowlist policy
- signature and digest verification
- warning banners for untrusted or unsigned remotes
