# driftcheck

> Lightweight utility to detect configuration drift between running containers and their source docker-compose or Kubernetes manifests.

---

## Installation

```bash
go install github.com/yourusername/driftcheck@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/driftcheck/releases).

---

## Usage

Point `driftcheck` at your manifest file and it will compare the declared configuration against what is currently running.

**Docker Compose:**
```bash
driftcheck --file docker-compose.yml
```

**Kubernetes:**
```bash
driftcheck --file deployment.yaml --context my-k8s-context
```

**Example output:**
```
[DRIFT] web: image mismatch — expected nginx:1.25, got nginx:1.23
[DRIFT] api: env var PORT missing from running container
[OK]    db: no drift detected
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to manifest file | `docker-compose.yml` |
| `--context` | Kubernetes context to use | current context |
| `--output` | Output format: `text`, `json` | `text` |
| `--quiet` | Only report drifted resources | `false` |

---

## How It Works

`driftcheck` reads your manifest file, queries the Docker daemon or Kubernetes API for the current state of each resource, and reports any fields that differ from the declared configuration.

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

MIT © 2024 yourusername