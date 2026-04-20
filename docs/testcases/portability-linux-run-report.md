# Linux Portability Run Report

## Latest run metadata

- Run ID: `2026-04-14_linux_tinygo_enabled`
- Date: 2026-04-14
- Host: Linux Mint 22.2 (`linuxmint`, `noble`)
- Kernel/arch: `Linux 6.8.0-106-generic`, `amd64`
- Go version: `go1.25.4`
- Script: `scripts/portability/run_portability_linux.sh`
- Artifacts directory: `artifacts/portability/linux/2026-04-14_linux_tinygo_enabled/`

## Latest execution summary

- Total executed: `11`
- Passed: `11`
- Failed: `0`
- Skipped: `0`
- Pass rate: `100.00%`

Failed case IDs:

- none

Skipped case IDs:

- none

## Progress from previous run

- Previous run ID: `2026-04-14_linux_local_rerun`
- Previous pass rate: `90.91%`
- Current pass rate: `100.00%`
- Improvements:
  - `PORT-LIN-XTR-001` now passes after adding missing `KYamlPreferences` and `ConfiguredKYamlPreferences` definitions to `pkg/yqlib/no_kyaml.go`.
  - `PORT-LIN-XTR-002` now passes after installing TinyGo and using a TinyGo-compatible Go toolchain (`go1.25.4`).

## Root cause and fix details (`PORT-LIN-XTR-001`)

- Classification: **legitimate code defect**, not a host-specific environment issue.
- Why:
  - `./scripts/build-small-yq.sh` builds with `-tags ... yq_nokyaml ...`.
  - Under `yq_nokyaml`, `pkg/yqlib/kyaml.go` is excluded (it normally defines `KYamlPreferences` and `ConfiguredKYamlPreferences`).
  - `pkg/yqlib/no_kyaml.go` still referenced `KYamlPreferences` in `NewKYamlEncoder` signature, and `pkg/yqlib/format.go` still referenced `ConfiguredKYamlPreferences`.
  - Result was deterministic compile failure:
    - `undefined: KYamlPreferences`
    - `undefined: ConfiguredKYamlPreferences`
- Fix applied:
  - Added fallback `KYamlPreferences` type and `ConfiguredKYamlPreferences` variable in `pkg/yqlib/no_kyaml.go` for the `yq_nokyaml` build tag path.
- Verification:
  - `./scripts/build-small-yq.sh` now exits `0`.
  - Full Linux portability run (`2026-04-14_linux_tinygo_enabled`) passes `PORT-LIN-XTR-001`.

## Generated evidence

- Latest run summary: `artifacts/portability/linux/2026-04-14_linux_tinygo_enabled/summary.md`
- Latest machine-readable results: `artifacts/portability/linux/2026-04-14_linux_tinygo_enabled/summary.csv`
- Latest per-case logs: `artifacts/portability/linux/2026-04-14_linux_tinygo_enabled/logs/`
- Previous run summaries:
  - `artifacts/portability/linux/2026-04-14_linux_local_rerun/summary.md`
  - `artifacts/portability/linux/2026-04-14_linux_local/summary.md`

## Recommended next actions

- Document Linux TinyGo prerequisites (`tinygo`, plus Go `1.19` through `1.25`) in the portability setup notes.
- Use the Windows run report (`docs/testcases/portability-windows-run-report.md`) to track and reduce the remaining `PORT-WIN-TST-001` instability.
