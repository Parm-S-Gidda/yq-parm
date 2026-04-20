# macOS Portability Test Cases for yq

## Scope

This document defines the first execution batch for portability testing on macOS.

- Host OS: macOS (record exact version per run)
- Architecture: arm64 or amd64 (record exact value per run)
- Shell: `zsh`
- Project path: repository root

## Evidence Output

Store logs under:

- `artifacts/portability/macos/<run-id>/`

Suggested `run-id` format:

- `YYYY-MM-DD_HHMM_local`

Each test case should generate:

- stdout and stderr log files
- exit status
- short pass/fail verdict

## Pre-check Cases

### PORT-MAC-ENV-001: Capture environment

- **Objective:** Record host/compiler details for traceability.
- **Command:**
  - `sw_vers`
  - `uname -a`
  - `go version`
  - `go env GOOS GOARCH`
  - `git rev-parse --short HEAD`
- **Expected:** Commands succeed and values are written to the run log.

### PORT-MAC-BLD-001: Build yq locally

- **Objective:** Confirm baseline build portability on macOS.
- **Command:** `go build ./...`
- **Expected:** Exit code `0`.

## Runtime Smoke Cases

### PORT-MAC-RT-001: Version command

- **Objective:** Confirm built binary can execute basic command.
- **Command:** `./yq --version`
- **Expected:** Exit code `0`, version text printed.

### PORT-MAC-RT-002: YAML scalar evaluation

- **Objective:** Validate basic YAML parsing and expression execution.
- **Input:**
  - `printf 'a: 1\n' > /tmp/port-mac-basic.yml`
- **Command:** `./yq '.a' /tmp/port-mac-basic.yml`
- **Expected output:** `1`
- **Expected exit code:** `0`

### PORT-MAC-RT-003: JSON input handling

- **Objective:** Validate JSON decode path and value extraction.
- **Input:**
  - `printf '{"name":"yq","ok":true}\n' > /tmp/port-mac-basic.json`
- **Command:** `./yq '.name' /tmp/port-mac-basic.json`
- **Expected output:** `"yq"`
- **Expected exit code:** `0`

### PORT-MAC-RT-004: Pipe input (stdin path)

- **Objective:** Validate stdin workflow in shell pipelines.
- **Command:** `printf 'k: v\n' | ./yq '.k' -`
- **Expected output:** `v`
- **Expected exit code:** `0`

### PORT-MAC-RT-005: UTF-8 handling

- **Objective:** Validate Unicode portability.
- **Input:**
  - `printf 'msg: "cafe 😊"\n' > /tmp/port-mac-utf8.yml`
- **Command:** `./yq '.msg' /tmp/port-mac-utf8.yml`
- **Expected output:** `cafe 😊`
- **Expected exit code:** `0`

## Project Test Script Cases

### PORT-MAC-TST-001: Project Go tests

- **Objective:** Confirm behaviour tests run on macOS.
- **Command:** `./scripts/test.sh`
- **Expected:** Exit code `0`.

### PORT-MAC-TST-002: Acceptance tests

- **Objective:** Confirm shell acceptance suite on macOS.
- **Command:** `./scripts/acceptance.sh`
- **Expected:** Exit code `0`.

## Extreme Configuration Cases

### PORT-MAC-XTR-001: Small build profile

- **Objective:** Validate reduced feature build profile.
- **Command:** `./scripts/build-small-yq.sh`
- **Expected:** Exit code `0`.

### PORT-MAC-XTR-002: TinyGo profile (optional)

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
