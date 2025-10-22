# kcskit 0.1.0 — Kaspersky Container Security CLI

kcskit is a small command-line utility to interact with Kaspersky Container Security (KCS).
It provides configuration management, a reusable API client, and commands to inspect a KCS deployment (health, integrations, images, clusters, scans).

## Features
- Save API base URL, token and optional CA certificate to a local YAML config (`$HOME/.kcskit/config`)
- Reusable API client that preserves the configured base URL and token
- Human-friendly tabbed-table output by default; raw JSON via `-o json`
- Global `-i` / `--invalid-cert` flag to skip TLS verification (lab/test only)
- `--ca_cert` accepts a PEM literal, a file path, or `-` to read from stdin
- Commands: config, config check, registries, images (list / scan), clusters

## Configuration

Save token and endpoint (endpoint must be a base URL including API path, e.g. `https://kcs.demo.lab/api/`):

```bash
kcskit config --token <token> --endpoint https://kcs.demo.lab/api/
```

CA certificate examples (`--ca_cert` accepts literal PEM, path, or `-` for stdin):

- From file:
  kcskit config --ca_cert /path/to/ca.pem --endpoint https://kcs.demo.lab/api/ --token kcs_...
- From stdin:
  cat /path/to/ca.pem | kcskit config --ca_cert - --endpoint https://kcs.demo.lab/api/ --token kcs_...
- Inline:
  kcskit config --ca_cert "$(cat /path/to/ca.pem)" --endpoint https://kcs.demo.lab/api/ --token kcs_...

Config YAML at `$HOME/.kcskit/config` contains `token`, `endpoint`, and optional `ca_cert`.

When `ca_cert` is present and `-i` is NOT used, the client will try to use the configured CA to validate TLS.

## Commands (examples)

- Help / basic:
  kcskit --help

- Config save:
  kcskit config --token <token> --endpoint https://kcs.demo.lab/api/ --ca_cert /path/to/ca.pem

- Health check (reads `/v1/core-health`):
  kcskit config check
  kcskit config check -o json
  kcskit config check -i

- Registries (list `/v1/integrations/image-registries`):
  kcskit registries list
  kcskit registries list -o json
  # Output columns: ID, Name, Type, Url

- Images list (`/v1/images/registry`) — flags map to query parameters:
  kcskit images list --registry <registry-id> --page 1 --limit 50 --sort name --by asc --name "Docker" --scopes scope1 --risks malware --output json
  # Flags:
  --page (int), --limit (int), --sort (name|riskRating), --by (asc|desc),
  --scopes (repeatable), --name, --registry, --repositoriesWith, --scannedAt, --risks (repeatable)
  # Output columns: ID, Name, Registry, Risk

- Images scan (`POST /v1/scans`):
  kcskit images scan --artifact nginx:latest --registry <registry-id>
  # Required flags: --artifact, --registry
  # Output columns: ID, Artifact, Scanner, Status

- Clusters list (`/v1/clusters`) — flags:
  kcskit clusters list --page 1 --limit 50 --sort clusterName --by asc --scopes scope1
  # Flags: --page, --limit, --sort (clusterName|orchestrator|namespaces|riskRating), --by, --scopes (repeatable)
  # Output columns: ID, Name, Orchestrator, Namespaces, Risk

Global flags:
- -i / --invalid-cert : ignore TLS validation for all commands
- -o / --output json   : print raw pretty JSON instead of tabbed table

## Project layout

- cmd/        — CLI commands (root, config, registries, images, clusters, ...)
- internal/
  - model/    — API models (config, health, registry, images, clusters, scans)
  - service/  — reusable API client and config file I/O
  - controller/— orchestration layer between cmd and service
- main.go

## Extending

- Add new cmd handlers that use `internal/service.NewClient(...)`.
- Prefer parsing API JSON into types in `internal/model` and present via tabwriter.
- Add unit tests using `httptest` for API client and temporary HOME for config I/O.

License: add a LICENSE file as desired.