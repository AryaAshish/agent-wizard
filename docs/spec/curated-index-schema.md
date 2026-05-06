# Curated Index Schema (Draft)

## File

`curated-index.yaml`

## Top-level fields

- `schemaVersion`
- `entries`

## Entry fields

- `id`: unique pack or skill id
- `kind`: `skill` or `pack`
- `title`
- `description`
- `source`: reference URL/path
- `tags`: list of taxonomy tags

## Tag taxonomy

- `stack:*` (e.g. `stack:android`, `stack:backend`, `stack:fullstack`)
- `area:*` (e.g. `area:frontend`, `area:security`)
- `phase:*` (e.g. `phase:planning`, `phase:release`)
- `workflow:*` (e.g. `workflow:pull-request-review`, `workflow:launch-readiness`)

## Governance requirements

- every entry has a maintainer contact
- every entry declares license metadata
- optional trust badge: `community`, `org-curated`, `verified`
