# harbor-cleaner

Connect to a Harbor registry, list image tags matching the given regex filters, and delete them after confirmation.

## Usage

```
harbor-cleaner [flags]
```

## Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--endpoint` | | `string` | Harbor API endpoint (e.g. `https://docker.example.com`) | |
| `--username` | | `string` | Harbor username | |
| `--password` | | `string` | Harbor password | |
| `--project` | | `string` | Harbor project name or ID | |
| `-f` | `-f` | `stringArray` | Regex filter for tag names (can be specified multiple times) | |
| `--concurrency` | | `int` | Number of concurrent delete operations | `4` |
| `-h` | `-h` | | Print help | |

## Example

```bash
harbor-cleaner \
  --endpoint https://docker.example.com \
  --username admin \
  --password Harbor12345 \
  --project my-project \
  -f "^v1\\.0\\..*" \
  -f "^dev-.*" \
  --concurrency 8
```

