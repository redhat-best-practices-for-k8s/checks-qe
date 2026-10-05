# checks-qe

[![CI](https://github.com/redhat-best-practices-for-k8s/checks-qe/actions/workflows/ci.yml/badge.svg)](https://github.com/redhat-best-practices-for-k8s/checks-qe/actions/workflows/ci.yml)
[![Build and push image](https://github.com/redhat-best-practices-for-k8s/checks-qe/actions/workflows/checks-qe-image.yaml/badge.svg)](https://github.com/redhat-best-practices-for-k8s/checks-qe/actions/workflows/checks-qe-image.yaml)

Scenario-based QE testing for the [checks](https://github.com/redhat-best-practices-for-k8s/checks) library. Each scenario creates real Kubernetes resources, runs the checks library's autodiscovery, executes a single check, and asserts the result matches expectations.

**Coverage:** 222 scenarios covering all 105 checks across 9 categories.

## Quick Start

### Prerequisites

- Go 1.27+
- Access to a Kubernetes or OpenShift cluster (`KUBECONFIG` set or `~/.kube/config` present)

### Build

```bash
go build -o checks-qe ./cmd/checks-qe/
```

### Check version

```bash
./checks-qe version
```

### List scenarios

```bash
./checks-qe list
```

### Run all scenarios

```bash
./checks-qe run --parallel 8 --verbose
```

### Run against vanilla Kubernetes (skip OCP and probe-dependent scenarios)

```bash
./checks-qe run --skip-ocp --skip-probe --parallel 8 --verbose
```

### Run a single category

```bash
./checks-qe run --category access-control --verbose
```

### Run a specific scenario

```bash
./checks-qe run --scenario "networking/dual-stack" --verbose
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `--parallel N` | Number of scenarios to run concurrently (default: 4) |
| `--verbose` | Show resource details on all results |
| `--skip-ocp` | Skip scenarios requiring OpenShift |
| `--skip-probe` | Skip scenarios requiring probe pods |
| `--category X` | Run only scenarios in this category |
| `--scenario X` | Run only scenarios matching this name pattern |
| `--kubeconfig` | Path to kubeconfig (default: `$KUBECONFIG` or `~/.kube/config`) |
| `--mode` | Execution mode: `direct` (default) or `operator` (via bps-operator) |
| `--operator-namespace` | Namespace where bps-operator is deployed (operator mode only, default: `bps-operator-system`) |

## Auto-Skip Capabilities

Scenarios automatically skip when cluster infrastructure is unavailable:

| Condition | Detection | Example scenarios |
|-----------|-----------|-------------------|
| Dual-stack | `kubernetes` service IP families | IPv6 ICMP, dual-stack service |
| Multus CNI | `k8s.cni.cncf.io/v1` API group | Multus connectivity checks |
| Probe pods | `--skip-probe` flag | ICMP, process, SSH, platform node checks |
| OpenShift | API detection / `--skip-ocp` | OCP lifecycle, cluster operator health |

## Project Structure

```
cmd/checks-qe/       CLI entry point
pkg/builder/          Resource builders (Deployment, Service, CSV, CRD, etc.)
pkg/cluster/          Cluster client and resource operations
pkg/discovery/        Kubernetes resource discovery
pkg/runner/           Scenario execution engine
pkg/scenario/         Scenario registration and types
scenarios/            Scenario definitions by category
  accesscontrol/
  certification/
  lifecycle/
  manageability/
  networking/
  observability/
  operator/
  performance/
  platform/
```

## Adding a New Scenario

1. Create or edit a file in `scenarios/<category>/`.
2. Define a `registerXxx()` function that calls `scenario.Register(...)`.
3. Add the call to the package's `register.go` `Register()` function.
4. Build and run `./checks-qe list` to verify.
