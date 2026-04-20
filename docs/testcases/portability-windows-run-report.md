# Windows Portability Run Report

## Latest run metadata

- Run ID: `2026-04-16_2149_local`
- Date: 2026-04-16
- Host: Windows 10 (`10.0.19045.6466`)
- Kernel/arch: `windows`, `amd64`
- Go version: `go1.26.2`
- Script: `scripts/portability/run_portability_windows.ps1`
- Artifacts directory: `artifacts/portability/windows/2026-04-16_2149_local/`

## Latest execution summary

- Total executed: `11`
- Passed: `9`
- Failed: `1`
- Skipped: `1`
- Pass rate: `81.82%`

Failed case IDs:

- `PORT-WIN-TST-002`

Skipped case IDs:

- `PORT-WIN-XTR-002`

## Progress from previous run

- Previous run ID: `2026-04-15_2145_local`
- Previous pass rate: `81.82%`
- Current pass rate: `81.82%`
- Improvements:
  - `PORT-WIN-TST-001` now passes after reader-close defect fixes in `pkg/yqlib`.
  - Runtime Windows smoke cases (`PORT-WIN-RT-001` through `PORT-WIN-RT-005`) are stable and passing.

## Remaining failure details (`PORT-WIN-TST-002`)

- Classification: **Windows-specific assertion mismatch in acceptance tests**.
- Why:
  - `./scripts/acceptance.sh` still has shell assertions that differ under Windows line ending/output behaviour.
  - In the latest run, `acceptance_tests/basic.sh` shows a value comparison mismatch in `testBasicNoExitStatus`.
- Evidence:
  - `artifacts/portability/windows/2026-04-16_2149_local/logs/PORT-WIN-TST-002.stdout.log`

## Post-run fix status (2026-04-16)

- Classification update: **legitimate code defect fixed**.
- Root cause:
  - Reader handles opened through buffered wrappers were not being closed in all execution paths.
  - Windows file-locking semantics exposed this as `TempDir` cleanup failures in `cmd` tests.
- Fix applied:
  - `pkg/yqlib/utils.go`: introduced a closeable buffered file reader and added guaranteed reader cleanup in document loading.
  - `pkg/yqlib/file_utils.go`: extended `SafelyCloseReader()` to close any `io.Closer`.
  - `pkg/yqlib/stream_evaluator.go`: ensured readers are closed on all evaluator paths.
- Verification performed:
  - `go test ./cmd -run "TestEvaluate(All|Sequence)_" -count=1` passes on Windows.
  - `bash ./scripts/test.sh` passes on Windows after the file-handle fix and script cleanup.
- Note:
  - The reader-close fix is confirmed by `PORT-WIN-TST-001` passing in run `2026-04-16_2149_local`.

## Generated evidence

- Latest run summary: `artifacts/portability/windows/2026-04-16_2149_local/summary.md`
- Latest machine-readable results: `artifacts/portability/windows/2026-04-16_2149_local/summary.csv`
- Latest per-case logs: `artifacts/portability/windows/2026-04-16_2149_local/logs/`
- Previous run summaries:
  - `artifacts/portability/windows/2026-04-15_2145_local/summary.md`
  - `artifacts/portability/windows/2026-04-15_2125_local/summary.md`
  - `artifacts/portability/windows/2026-04-15_2104_local/summary.md`

## Recommended next actions

- Keep the Windows runner and case-ID contract unchanged for trend continuity.
- Keep triaging `PORT-WIN-TST-002` acceptance assertions for Windows line ending/output compatibility.
- Keep `PORT-WIN-XTR-002` as optional and `SKIPPED` until TinyGo is intentionally installed on the Windows host.

## `PORT-WIN-TST-002` troubleshooting log (2026-04-17)

This section records the full investigation and fix attempts so the process is traceable.

### Starting point

- Latest full run at the time: `2026-04-16_2202_local`
- Status:
  - `PORT-WIN-TST-001`: `PASS`
  - `PORT-WIN-TST-002`: `FAIL`
- Failure evidence:
  - `artifacts/portability/windows/2026-04-16_2202_local/logs/PORT-WIN-TST-002.stdout.log`
  - `artifacts/portability/windows/2026-04-16_2202_local/logs/PORT-WIN-TST-002.stderr.log`

### Attempt 1: diagnose failure pattern

- Observed failures in `acceptance_tests/basic.sh` with errors such as:
  - `Error: no support for j output format`
  - assertions unexpectedly returning `null`/empty output
- Initial hypothesis: acceptance scripts might be invoking a stale/reduced `./yq` binary in Git Bash context.

### Attempt 2: ensure shell uses `yq.exe`

- Added a Windows bash launcher in `scripts/portability/run_portability_windows.ps1`:
  - creates `./yq` wrapper that forwards to `./yq.exe`
- Result:
  - `PORT-WIN-TST-001` remained `PASS`
  - `PORT-WIN-TST-002` still failed in full run
  - indicates launcher alone was not sufficient.

### Attempt 3: investigate script execution environment

- Checked script file types under Git Bash:
  - `acceptance_tests/basic.sh` and `scripts/acceptance.sh` had CRLF line terminators on this Windows checkout.
- Hypothesis: CRLF script execution in bash was causing assertion/command parsing inconsistencies.

### Fix applied

- Updated `scripts/acceptance.sh` to run acceptance tests with `igncr` only on Windows-like shells:
  - `MINGW*|MSYS*|CYGWIN*` -> `bash -o igncr "$test"`
  - all other OSes -> `bash "$test"`
- This keeps Linux/macOS behaviour unchanged while normalising CRLF handling on Windows.

### Validation after fix

- Rebuilt full binary and reran acceptance suite directly on Windows:
  - `go build -o yq.exe .`
  - `bash ./scripts/acceptance.sh`
- Result:
  - `acceptance_tests/basic.sh`: `OK`
  - full `scripts/acceptance.sh`: `exit_code 0` in local validation run
  - no remaining `TST-002` assertion failures in this isolated verification.

### Current interpretation

- `PORT-WIN-TST-002` root cause was a Windows shell/script line-ending execution mismatch (CRLF handling), not core yq data-processing logic.
- The isolated fix is in place and validated locally.
- `PORT-WIN-TST-002` is treated as passed for now based on the isolated validation result.
