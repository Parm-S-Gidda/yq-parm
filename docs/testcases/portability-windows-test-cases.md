# Windows Portability Test Cases for yq

## Scope

This document defines the first execution batch for portability testing on Windows.

- Host OS: Windows 10/11 (record exact version per run)
- Architecture: amd64 or arm64 (record exact value per run)
- Shell: `powershell`
- Project path: repository root

## Evidence Output

Store logs under:

- `artifacts/portability/windows/<run-id>/`

Suggested `run-id` format:

- `YYYY-MM-DD_HHMM_local`

Each test case should generate:

- stdout and stderr log files
- exit status
- short pass/fail verdict

## Pre-check Cases

### PORT-WIN-ENV-001: Capture environment

- **Objective:** Record host/compiler details for traceability.
- **Command:**
  - `cmd /c ver`
  - `go version`
  - `go env GOOS GOARCH`
  - `git rev-parse --short HEAD`
- **Expected:** Commands succeed and values are written to the run log.

### PORT-WIN-BLD-001: Build yq locally

- **Objective:** Confirm baseline build portability on Windows.
- **Command:** `go build ./...`
- **Expected:** Exit code `0`.

## Runtime Smoke Cases

### PORT-WIN-RT-001: Version command

- **Objective:** Confirm built binary can execute basic command.
- **Command:** `go build -o .\yq.exe .; .\yq.exe --version`
- **Expected:** Exit code `0`, version text printed.

### PORT-WIN-RT-002: YAML scalar evaluation

- **Objective:** Validate basic YAML parsing and expression execution.
- **Input:** create `basic.yml` with `a: 1`
- **Command:** `.\yq.exe '.a' <input>`
- **Expected output:** `1`
- **Expected exit code:** `0`

### PORT-WIN-RT-003: JSON input handling

- **Objective:** Validate JSON decode path and value extraction.
- **Input:** create `basic.json` with `{"name":"yq","ok":true}`
- **Command:** `.\yq.exe '.name' <input>`
- **Expected output:** `"yq"`
- **Expected exit code:** `0`

### PORT-WIN-RT-004: Pipe input (stdin path)

- **Objective:** Validate stdin workflow in shell pipelines.
- **Command:** `"k: v`n" | .\yq.exe '.k' -`
- **Expected output:** `v`
- **Expected exit code:** `0`

### PORT-WIN-RT-005: UTF-8 handling

- **Objective:** Validate Unicode portability on Windows.
- **Input:** create `utf8.yml` with `msg: "cafe 😊"` (UTF-8 without BOM)
- **Command:** `.\yq.exe '.msg' <input>`
- **Expected output:** `cafe 😊`
- **Expected exit code:** `0`

## Project Test Script Cases

### PORT-WIN-TST-001: Project Go tests

- **Objective:** Confirm behaviour tests run on Windows.
- **Command:** `bash ./scripts/test.sh`
- **Expected:** Exit code `0`.

### PORT-WIN-TST-002: Acceptance tests

- **Objective:** Confirm shell acceptance suite on Windows.
- **Command:** `bash ./scripts/acceptance.sh`
- **Expected:** Exit code `0`.

## Extreme Configuration Cases

### PORT-WIN-XTR-001: Small build profile

- **Objective:** Validate reduced feature build profile.
- **Command:** `bash ./scripts/build-small-yq.sh`
- **Expected:** Exit code `0`.

### PORT-WIN-XTR-002: TinyGo profile (optional)

- **Objective:** Validate alternative compiler profile when available.
- **Command:** `bash ./scripts/build-tinygo-yq.sh`
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
