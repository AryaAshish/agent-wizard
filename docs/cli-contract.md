# CLI contract (v1)

## Grammar

- Prefer subcommands: `init`, `list`, `create-skill`, `add`, `sync`, `status`, `sources`, `lock`, `migrate`, `cache`, `ci-check`, `browse`, `watch`, `import`, `pack`, `catalog`, `icp`, `wizard` (alias `guide`).
- Flags complement subcommands; keep stable long flags over time.
- **`list`:** prints skill id and a short summary (first paragraph under the title in `SKILL.md`, aligned columns; sorted by id). Empty results print hints.
- **`create-skill`:** scaffolds `./<skill-id>/SKILL.md` locally; no upload.
- **`add` (v0.1.3+):** if `agentskills.yaml` is missing, creates it (headless init with bundled **community** source), appends the skill, then runs **`sync`** unless **`--no-sync`**. Older releases required **`init`** before the first **`add`**.
- **Project directory:** resolve by walking parents from `cwd` for `agentskills.yaml`; commands that need an existing manifest (`sync`, `status`, …) use that root. When creating a **new** manifest on first **`add`**, the directory is the nearest ancestor containing **`.git`**, otherwise **`cwd`**.
- **`wizard` / `guide`:** optional interactive menu when stdin/stdout are a TTY; scripts keep flag-based commands unchanged.

## Machine output

- `status --json` emits a single JSON object suitable for CI parsing.
- Exit codes:
  - `0` success
  - `1` general failure
  - `3` lockfile drift detected (`status --check-drifts`)
