# 🛡️ kcskit v0.1.0-alpha — Kaspersky Container Security CLI

`kcskit` is a lightweight command-line utility for interacting with the Kaspersky Container Security (KCS) API. It simplifies managing and inspecting a KCS deployment directly from your terminal.

## 📚 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
  - [Global Flags](#global-flags)
  - [Command Examples](#command-examples)
- [Project Layout](#-project-layout)
- [Extending](#-extending)

## 🌟 Features

- 💾 Save API base URL, token and optional CA certificate to a local YAML config (`$HOME/.kcskit/config`)
- 🔌 Reusable API client that preserves the configured base URL and token
- 📊 Human-friendly tabbed-table output by default; raw JSON via `-o json`
- 🔐 Global `-i` / `--invalid-cert` flag to skip TLS verification (lab/test only)
- 📄 `--ca_cert` accepts a PEM literal, a file path, or `-` to read from stdin
- 🛠️ Commands: `config`, `config check`, `registries`, `images` (list / scan), `clusters`

## 📋 Prerequisites

- Access to a **Kaspersky Container Security (KCS)** instance
- An **API token** generated from your KCS user profile
- **Go 1.18+** (if building from source)

## ⬇️ Installation

### From Source

```bash
go install github.com/arturscheiner/kcskit@latest
```

## ⚙️ Configuration

Save token and endpoint (endpoint must be a base URL including API path, e.g. `https://kcs.demo.lab/api/`):

```bash
kcskit config --token <token> --endpoint https://kcs.demo.lab/api/
```

### CA Certificate Examples

`--ca_cert` accepts literal PEM, path, or `-` for stdin:

- **From file:**
  ```bash
  kcskit config --ca_cert /path/to/ca.pem --endpoint https://kcs.demo.lab/api/ --token kcs_...
  ```

- **From stdin:**
  ```bash
  cat /path/to/ca.pem | kcskit config --ca_cert - --endpoint https://kcs.demo.lab/api/ --token kcs_...
  ```

- **Inline:**
  ```bash
  kcskit config --ca_cert "$(cat /path/to/ca.pem)" --endpoint https://kcs.demo.lab/api/ --token kcs_...
  ```

Config YAML at `$HOME/.kcskit/config` contains `token`, `endpoint`, and optional `ca_cert`.

When `ca_cert` is present and `-i` is NOT used, the client will try to use the configured CA to validate TLS.

## 🖥️ Usage

### Basic Help

```bash
kcskit --help
```

### Configuration Commands

- **Save configuration:**
  ```bash
  kcskit config --token <token> --endpoint https://kcs.demo.lab/api/ --ca_cert /path/to/ca.pem
  ```

- **Health check** (reads `/v1/core-health`):
  ```bash
  kcskit config check
  kcskit config check -o json
  kcskit config check -i
  ```

### Registry Commands

- **List registries** (`/v1/integrations/image-registries`):
  ```bash
  kcskit registries list
  kcskit registries list -o json
  ```
  _Output columns: ID, Name, Type, Url_

### Image Commands

- **List images** (`/v1/images/registry`):
  ```bash
  kcskit images list --registry <registry-id> --page 1 --limit 50 --sort name --by asc --name "Docker" --scopes scope1 --risks malware --output json
  ```
  
  **Flags:**
  - `--page` (int)
  - `--limit` (int)
  - `--sort` (name|riskRating)
  - `--by` (asc|desc)
  - `--scopes` (repeatable)
  - `--name`
  - `--registry`
  - `--repositoriesWith`
  - `--scannedAt`
  - `--risks` (repeatable)
  
  _Output columns: ID, Name, Registry, Risk_

- **Scan images** (`POST /v1/scans`):
  ```bash
  kcskit images scan --artifact nginx:latest --registry <registry-id>
  ```
  
  **Required flags:** `--artifact`, `--registry`
  
  _Output columns: ID, Artifact, Scanner, Status_

### Cluster Commands

- **List clusters** (`/v1/clusters`):
  ```bash
  kcskit clusters list --page 1 --limit 50 --sort clusterName --by asc --scopes scope1
  ```
  
  **Flags:**
  - `--page`
  - `--limit`
  - `--sort` (clusterName|orchestrator|namespaces|riskRating)
  - `--by`
  - `--scopes` (repeatable)
  
  _Output columns: ID, Name, Orchestrator, Namespaces, Risk_

### Global Flags

- `-i` / `--invalid-cert` : Ignore TLS validation for all commands
- `-o` / `--output json` : Print raw pretty JSON instead of tabbed table

## 📁 Project Layout

```
- cmd/              — CLI commands (root, config, registries, images, clusters, ...)
- internal/
  - model/          — API models (config, health, registry, images, clusters, scans)
  - service/        — reusable API client and config file I/O
  - controller/     — orchestration layer between cmd and service
- main.go
```

## 🧩 Extending

- Add new cmd handlers that use `internal/service.NewClient(...)`
- Prefer parsing API JSON into types in `internal/model` and present via tabwriter
- Add unit tests using `httptest` for API client and temporary HOME for config I/O

## 📄 License

Add a LICENSE file as desired.