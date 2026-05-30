# harbor-cleaner

Connect to a Harbor registry, list image tags matching the given regex filters, and delete them after confirmation.

## Usage

```
harbor-cleaner [global flags] <command> [command flags]
```

## Commands

| Command | Description |
|---------|-------------|
| `clean` | Delete tags matching regex filters (with confirmation) |
| `tags` | List all tags in a project |
| `repos` | List all repos in a project with storage size |
| `quotas` | Show quota usage of all projects |
| `version` | Print version and build info |

## Global Flags

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--endpoint` | `string` | Harbor API endpoint (e.g. `https://docker.example.com`) | |
| `--username` | `string` | Harbor username | |
| `--password` | `string` | Harbor password | |
| `--concurrency` | `int` | Number of concurrent delete operations | `4` |

## `clean` Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--project` | | `string` | Harbor project name or ID (required) |
| `--filter` | `-f` | `stringArray` | Regex filter for tag names (can be specified multiple times, required) |

## `tags` Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--project` | | `string` | Harbor project name or ID | |
| `--order` | | `string` | Order by field (`name` or `size`) | `name` |

## `repos` Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--project` | | `string` | Harbor project name or ID | |
| `--order` | | `string` | Order by field (`name` or `size`) | `name`

## `quotas` Flags

No additional flags.

## Examples

Delete tags matching regex filters:

```bash
harbor-cleaner \
  --endpoint https://docker.example.com \
  --username admin \
  --password Harbor12345 \
  clean \
  --project my-project \
  -f "^v1\\.0\\..*" \
  -f "^dev-.*" \
  --concurrency 8
```

List all tags in a project:

```bash
harbor-cleaner \
  --endpoint https://docker.example.com \
  --username admin \
  --password Harbor12345 \
  tags \
  --project my-project
```

List all repos in a project (ordered by size):

```bash
harbor-cleaner \
  --endpoint https://docker.example.com \
  --username admin \
  --password Harbor12345 \
  repos \
  --project my-project \
  --order size
```

Show quota usage of all projects:

```bash
harbor-cleaner \
  --endpoint https://docker.example.com \
  --username admin \
  --password Harbor12345 \
  quotas
```

## Build

```bash
make build                  # build for linux/darwin amd64/arm64
make install                # build and install to /usr/local/bin
make build VERSION=0.0.1    # build with a specific version
make install VERSION=0.0.1  # install with a specific version
```

