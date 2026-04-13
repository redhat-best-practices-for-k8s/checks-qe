# CLAUDE.md

This file provides guidance to Claude Code when working with the checks-qe codebase.

## What This Repo Does

checks-qe is a scenario-based QE test suite for the [checks](https://github.com/redhat-best-practices-for-k8s/checks) library. It validates each check by creating real Kubernetes resources, running discovery, executing the check function, and asserting the compliance result.

## Key Commands

```bash
go build -o checks-qe ./cmd/checks-qe/   # Build
./checks-qe list                           # List all 222 scenarios
./checks-qe run --skip-ocp --skip-probe --parallel 8 --verbose  # Run (K8s)
./checks-qe run --skip-probe --parallel 8 --verbose             # Run (OCP)
golangci-lint run --timeout=5m             # Lint
```

## Architecture

- **Scenarios** register via explicit `Register()` calls in `cmd/checks-qe/main.go` (no `init()` functions).
- Each scenario package has a `register.go` that calls `registerXxx()` from each file.
- **Builders** (`pkg/builder/`) construct K8s resources: Deployment, Service, Pod, PDB, CSV, CRD, StatefulSet, ResourceQuota, NetworkPolicy, RBAC.
- **PostDiscovery** hooks inject synthetic resources (nodes, PVs, CSVs, ClusterOperators) into `DiscoveredResources` before the check runs. This enables testing checks that inspect cluster-level resources without special infrastructure.
- **MockCertValidator** (`pkg/builder/cert_validator.go`) implements `checks.CertificationValidator` for certification scenarios.
- **K8sClientset** is injected into `DiscoveredResources` by the runner for checks that need it (scaling, logging, Helm version).

## Skip Conditions

Scenarios declare requirements via struct fields: `RequiresOCP`, `RequiresProbe`, `RequiresDualStack`, `RequiresMultus`. The runner auto-detects cluster capabilities and skips scenarios when requirements are not met.

## Adding Scenarios

1. Add a `registerXxx()` function in `scenarios/<category>/<file>.go`
2. Call it from `scenarios/<category>/register.go`
3. Use `scenario.VanillaDeploymentSetup` for simple cases
4. Use `PostDiscovery` to inject synthetic resources for checks that need cluster-level data
5. Run `go build ./...` and `golangci-lint run --timeout=5m` before committing

## CI

The CI workflow (`.github/workflows/ci.yml`) runs 4 jobs: lint, build, integration-k8s (Kind cluster), integration-ocp (CRC cluster). Integration jobs use `--parallel 8`.
