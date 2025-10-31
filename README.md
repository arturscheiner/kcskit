# 🛡️ kcskit v0.1.5 — Kaspersky Container Security CLI

`kcskit` is a small CLI for interacting with a Kaspersky Container Security (KCS) API. It provides lightweight commands to configure the client, inspect registry and image inventory, create scan jobs, and query clusters. Output is human-friendly by default and can emit JSON or be sent to the AI assistant integration.

Note: kcskit is a proof-of-work (PoW) project created to demonstrate the potential of building a CLI tool that interacts with Kaspersky Container Security APIs. It is not an official Kaspersky product and is intended for demonstration and integration experiments only.

## 📚 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Project Layout](#-project-layout)
- [Extending](#-extending)

## 🌟 Features

- Persistent local configuration stored in `$HOME/.kcskit/config` (token, endpoint, optional ca_cert)
- Reusable API client (`internal/service`) that honors configured endpoint and token
- Tabbed-table human output by default, `-o json` for pretty JSON, `-o ai` to send results to the AI model
- Global `-i` / `--invalid-cert` to skip TLS verification (use only in test/lab)
- `--ca_cert` accepts a PEM literal, a file path, or `-` to read from stdin
- Commands implemented: `config`, `config check`, `registries list`, `images list`, `images scan`, `clusters list`, `cicd list`

## 📋 Prerequisites

- Access to a KCS instance
- An API token with the required permissions
- Go 1.18+ if building from source

## ⬇️ Installation

From source:

```bash
go install github.com/arturscheiner/kcskit@latest
```

## ⚙️ Configuration

Save token and endpoint (endpoint should include the API base path, e.g. `https://kcs.demo.lab/api/`):

```bash
kcskit config --token <token> --endpoint https://kcs.demo.lab/api/ --ca_cert /path/to/ca.pem
```

`--ca_cert` examples:

- From file:

```bash
kcskit config --ca_cert /path/to/ca.pem --endpoint https://kcs.demo.lab/api/ --token kcs_...
```

- From stdin:

```bash
cat /path/to/ca.pem | kcskit config --ca_cert - --endpoint https://kcs.demo.lab/api/ --token kcs_...
```

When a `ca_cert` is configured (and `-i` is not used), the client uses it to validate TLS.

## 🖥️ Usage

Run the basic help to see top-level commands and flags:

```bash
kcskit --help
```

Global flags available to most commands:

- `-i`, `--invalid-cert` : ignore TLS validation (lab/test only)
- `-o`, `--output` : `json` (pretty JSON), `ai` (send results to the AI assistant), or omitted for tabbed table

### Configuration commands

- Save configuration:

```bash
kcskit config --token <token> --endpoint https://kcs.demo.lab/api/ --ca_cert /path/to/ca.pem
```

- Health check (reads `/v1/core-health`):

```bash
kcskit config check
kcskit config check -o json
kcskit config check -i
```

### Registries

- List configured image registries (`GET /v1/integrations/image-registries`):

```bash
kcskit registries list --page 1 --limit 50 --sort name --by asc
kcskit registries list -o json
kcskit registries list -o ai
```

Flags for `registries list`:

- `--page` (int) — page number (default: 1)
- `--limit` (int) — items per page (default: 50)
- `--sort` — sort field (updatedAt|name|description|type|url|createdAt|status)
- `--by` — sort order (asc|desc)

Output columns (default tabbed table): `ID`, `Name`, `Type`, `URL`, `Auth`

### Images

- List images for a registry (`GET /v1/images/registry`):

```bash
kcskit images list --registry <registry-id> --page 1 --limit 50 --sort name --by asc
kcskit images list -o json
```

Flags for `images list` include `--page`, `--limit`, `--sort`, `--by`, `--scopes` (repeatable), `--name`, `--registry`, `--repositoriesWith`, `--scannedAt`, `--risks` (repeatable).

Output columns (default): `ID`, `Name`, `Registry`, `Risk`

- Create a scan job (`POST /v1/scans`):

```bash
kcskit images scan --artifact nginx:latest --registry <registry-id>
```

Notes for `images scan`:

- `--artifact` and `--registry` are required flags.
- If the artifact value does not include a tag or digest (for example `nginx`), the CLI will append `:latest` automatically before sending the request (so `nginx` → `nginx:latest`).

Output columns: `ID`, `Artifact`, `Scanner`, `Status` (or `-o json` / `-o ai`).

### Clusters

- List clusters (`GET /v1/clusters`):

```bash
kcskit clusters list --page 1 --limit 50 --sort clusterName --by asc
```

Flags: `--page`, `--limit`, `--sort`, `--by`, `--scopes`.

Output columns: `ID`, `Name`, `Orchestrator`, `Namespaces`, `Risk`

### CI/CD scans

- List CI/CD scans (`GET /v1/scans/ci-cd`):

```bash
kcskit cicd list --page 1 --limit 50 --sort createdAt --by desc
```

## 📁 Project Layout

```
- cmd/              — CLI commands (root, config, registries, images, clusters, cicd, ...)
- internal/
  - model/          — API models (config, health, registry, images, clusters, scans)
  - service/        — reusable API client and config file I/O
  - controller/     — orchestration layer between cmd and service
- main.go
```

## 🧩 Extending

- Add new command handlers in `cmd/` that call helper functions in `internal/controller` and `internal/service`.
- Parse API JSON into types in `internal/model` and present via `text/tabwriter` for consistent output.
- New features should include small unit tests; use `httptest` to mock API responses and a temporary HOME for config I/O.

## 📄 License

Add a LICENSE file as desired.