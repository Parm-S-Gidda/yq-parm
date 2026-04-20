# macOS Portability Run Report

## Latest run metadata

- Run ID: `2026-04-15_2235_local`
- Date: 2026-04-15
- Host: macOS
- Script: `scripts/portability/run_portability_macos.sh`
- Artifacts directory: `artifacts/portability/macos/2026-04-15_2235_local/`

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

- Previous run ID: `2026-04-12_1905_local`
- Previous pass rate: `90.91%`
- Current pass rate: `100.00%`
- Improvements:
  - `PORT-MAC-XTR-001` now passes after adding missing `KYamlPreferences` and `ConfiguredKYamlPreferences` definitions to `pkg/yqlib/no_kyaml.go`.
  - Runner now rebuilds full `./yq` before runtime checks, preventing stale small-build binaries from affecting later cases.
  - `PORT-MAC-XTR-002` now passes by using a local TinyGo toolchain and local Go 1.25 wrapper from cache paths outside the repository.

## Root cause and fix details (`PORT-MAC-XTR-001`)

- Classification: **legitimate code defect**, not a host-specific environment issue.
- Why:
  - `./scripts/build-small-yq.sh` builds with `-tags ... yq_nokyaml ...`.
  - Under `yq_nokyaml`, `pkg/yqlib/kyaml.go` is excluded (it normally defines `KYamlPreferences` and `ConfiguredKYamlPreferences`).
  - `pkg/yqlib/no_kyaml.go` still referenced `KYamlPreferences` in `NewKYamlEncoder` signature, and `pkg/yqlib/format.go`/`cmd/utils.go` still referenced `ConfiguredKYamlPreferences`.
  - Result was deterministic compile failure:
    - `undefined: KYamlPreferences`
    - `undefined: ConfiguredKYamlPreferences`
- Fix applied:
  - Added fallback `KYamlPreferences` type and `ConfiguredKYamlPreferences` variable in `pkg/yqlib/no_kyaml.go` for the `yq_nokyaml` build tag path.
- Verification:
  - `./scripts/build-small-yq.sh` now exits `0`.
  - Full mac portability rerun (`2026-04-12_1905_local`) passes `PORT-MAC-XTR-001`.

## Generated evidence

- Latest run summary: `artifacts/portability/macos/2026-04-15_2235_local/summary.md`
- Latest machine-readable results: `artifacts/portability/macos/2026-04-15_2235_local/summary.csv`
- Latest per-case logs: `artifacts/portability/macos/2026-04-15_2235_local/logs/`
- Previous run summaries:
  - `artifacts/portability/macos/2026-04-15_2200_local/summary.md`
  - `artifacts/portability/macos/2026-04-12_1815_local/summary.md`
  - `artifacts/portability/macos/2026-04-12_1735_local/summary.md`

## Recommended next actions

- Keep reusing the same case IDs and summary format across Linux, macOS, and Windows runs to compare portability outcomes consistently.
- Keep local TinyGo/Go shim assets outside the repository tree to avoid impacting `go install` compatibility tests.
