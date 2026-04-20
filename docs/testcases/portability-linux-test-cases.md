# Linux Portability Test Cases for yq

## Scope

This document defines the first execution batch for portability testing on Linux.

- Host OS: Linux (record exact distro/version per run)
- Architecture: amd64 or arm64 (record exact value per run)
- Shell: `zsh`
- Project path: repository root

## Evidence Output

Store logs under:

- `artifacts/portability/linux/<run-id>/`

Suggested `run-id` format:

- `YYYY-MM-DD_HHMM_local`

Each test case should generate:

- stdout and stderr log files
- exit status
- short pass/fail verdict

## Pre-check Cases

### PORT-LIN-ENV-001: Capture environment

- **Objective:** Record host/compiler details for traceability.
- **Command:**
  - `cat /etc/os-release` (when available)
  - `uname -a`
  - `go version`
  - `go env GOOS GOARCH`
  - `git rev-parse --short HEAD`
- **Expected:** Commands succeed and values are written to the run log.

### PORT-LIN-BLD-001: Build yq locally

- **Objective:** Confirm baseline build portability on Linux.
- **Command:** `go build ./...`
- **Expected:** Exit code `0`.

## Runtime Smoke Cases

### PORT-LIN-RT-001: Version command

- **Objective:** Confirm built binary can execute basic command.
- **Command:** `./yq --version`
- **Expected:** Exit code `0`, version text printed.

### PORT-LIN-RT-002: YAML scalar evaluation

- **Objective:** Validate basic YAML parsing and expression execution.
- **Input:**
  - `printf 'a: 1\n' > /tmp/port-lin-basic.yml`
- **Command:** `./yq '.a' /tmp/port-lin-basic.yml`
- **Expected output:** `1`
- **Expected exit code:** `0`

### PORT-LIN-RT-003: JSON input handling

- **Objective:** Validate JSON decode path and value extraction.
- **Input:**
  - `printf '{"name":"yq","ok":true}\n' > /tmp/port-lin-basic.json`
- **Command:** `./yq '.name' /tmp/port-lin-basic.json`
- **Expected output:** `"yq"`
- **Expected exit code:** `0`

### PORT-LIN-RT-004: Pipe input (stdin path)

- **Objective:** Validate stdin workflow in shell pipelines.
- **Command:** `printf 'k: v\n' | ./yq '.k' -`
- **Expected output:** `v`
- **Expected exit code:** `0`

### PORT-LIN-RT-005: UTF-8 handling

- **Objective:** Validate Unicode portability.
- **Input:**
  - `printf 'msg: "cafe 😊"\n' > /tmp/port-lin-utf8.yml`
- **Command:** `./yq '.msg' /tmp/port-lin-utf8.yml`
- **Expected output:** `cafe 😊`
- **Expected exit code:** `0`

## Project Test Script Cases

### PORT-LIN-TST-001: Project Go tests

- **Objective:** Confirm behaviour tests run on Linux.
- **Command:** `./scripts/test.sh`
- **Expected:** Exit code `0`.

### PORT-LIN-TST-002: Acceptance tests

- **Objective:** Confirm shell acceptance suite on Linux.
- **Command:** `./scripts/acceptance.sh`
- **Expected:** Exit code `0`.

## Extreme Configuration Cases

### PORT-LIN-XTR-001: Small build profile

- **Objective:** Validate reduced feature build profile.
- **Command:** `./scripts/build-small-yq.sh`
- **Expected:** Exit code `0`.

### PORT-LIN-XTR-002: TinyGo profile (optional)

- **Objective:** Validate alternative compiler profile when available.
- **Command:** `./scripts/build-tinygo-yq.sh`
- **Expected:** Exit code `0` when TinyGo is installed; otherwise mark as `SKIPPED`.

## Incident Reporting Template

For each failed test, record:

- Test case ID
- Environment details
- Exact command
- Expected result
- Actual result
- Exit code
- Severity (`low`/`medium`/`high`)
- Reproduction notes

## Summary Fields to Report

At the end of a run, report:

- total executed
- passed
- failed
- skipped
- pass rate
- failed case IDs
