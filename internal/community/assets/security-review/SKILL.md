# Security review (triage)

Triage **secrets, injection, authz, crypto, SSRF** from code/diff context with severity tags.

## When to use

- PR touches auth, parsing, subprocess, SQL, templates, file IO, or network clients.

## When not to use

- No code or diff supplied and no named entrypoints.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_SCOPE` (paths or `BASE...HEAD`)

## Outputs

```
FINDINGS:
- [CRITICAL|HIGH|MED|LOW] YOUR_PATH → title
  evidence: one line from diff or grep context
  exploitability: one sentence
  fix: one concrete mitigation

FALSE_POSITIVES_DISMISSED:
- bullet or "- none -"

OPEN_ITEMS:
- need dynamic test | need threat model | "- none -"
```

## Steps

1. Establish scope.

```bash
cd YOUR_REPO_ROOT
git diff YOUR_BASE_REF...YOUR_HEAD_REF --name-only 2>/dev/null | head -n 200
```

2. High-signal greps (tune extensions to repo).

```bash
grep -RInE 'exec\.|spawn\(|child_process|eval\(|innerHTML|dangerouslySetInnerHTML|pickle\.loads|yaml\.load\(|requests\.get\(|http\.Get\(|fmt\.Sprintf\(.*SELECT|SELECT \+|sqlite3\.connect\(' -- YOUR_SCOPE 2>/dev/null | head -n 80
find YOUR_SCOPE -maxdepth 6 -type f \( -name '*.pem' -o -name '*.key' -o -name '*id_rsa*' \) 2>/dev/null
```

3. For each finding, tie to **one** evidence line; otherwise omit.

## Stop and ask

Stop if `YOUR_SCOPE` is the whole repo without subpath limits and exceeds ~200 files changed—request narrower paths.

## Reject if

- A `CRITICAL` without plausible attacker control of the sink.

## Safety

- Redact secrets; never echo private keys into the structured output.
