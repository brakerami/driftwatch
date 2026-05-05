# driftwatch

Lightweight daemon that detects config drift between running containers and their source manifests.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Point driftwatch at your manifests directory and let it run in the background:

```bash
driftwatch --manifests ./k8s --interval 30s
```

Example output:

```
[DRIFT] container=api-server field=image expected=nginx:1.25 got=nginx:1.23
[OK]    container=worker
[DRIFT] container=sidecar field=env.LOG_LEVEL expected=info got=debug
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manifests` | `./manifests` | Path to source manifest files |
| `--interval` | `60s` | How often to check for drift |
| `--output` | `text` | Output format: `text` or `json` |
| `--alert-webhook` | — | Webhook URL to POST drift events |

---

## How It Works

driftwatch polls your running containers via the Docker or Kubernetes API and compares their live configuration against the YAML/JSON manifests on disk. Any deviation in image tags, environment variables, resource limits, or labels is reported as drift.

---

## Requirements

- Go 1.21+
- Docker Engine or a Kubernetes cluster (in-cluster or kubeconfig)

---

## License

MIT © 2024 yourorg