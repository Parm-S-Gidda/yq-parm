# Alpine Linux Portability Run Report

## Latest run metadata

- Run ID: `2026-04-19_alpine_tinygo_enabled`
- Date: 2026-04-19
- Host: Alpine Linux v3.23 (`alpine`)
- Kernel/arch: `Linux 6.17.2-1-pve`, `amd64`
- Go version: `go1.25.9`
- TinyGo version: `0.40.1` (LLVM `20.1.1`)
- Script: `scripts/portability/run_portability_linux.sh`
- Artifacts directory: `artifacts/portability/linux/2026-04-19_alpine_tinygo_enabled/`

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

## Progress across Alpine runs

- Initial run (`2026-04-19_alpine_local`):
  - pass rate: `90.91%`
  - status: `PORT-LIN-XTR-002` skipped (`tinygo` not installed)
- Regression-check run after reverting `no_kyaml` fallback (`2026-04-19_alpine_rerun_after_revert`):
  - pass rate: `81.82%`
  - failed case: `PORT-LIN-XTR-001`
  - error:
    - `undefined: KYamlPreferences`
    - `undefined: ConfiguredKYamlPreferences`
- Current run (`2026-04-19_alpine_tinygo_enabled`):
  - pass rate: `100.00%`
  - `PORT-LIN-XTR-001` passes again with the fallback definitions present
  - `PORT-LIN-XTR-002` passes after TinyGo installation

## Root cause confirmation (`PORT-LIN-XTR-001`)

- Classification: **legitimate code defect**, reproducible on Alpine when fallback symbols are removed.
- Why:
  - `./scripts/build-small-yq.sh` builds with `yq_nokyaml`, excluding `pkg/yqlib/kyaml.go`.
  - Without fallback definitions in `pkg/yqlib/no_kyaml.go`, references to `KYamlPreferences` and `ConfiguredKYamlPreferences` fail during compile.
- Evidence:
  - `artifacts/portability/linux/2026-04-19_alpine_rerun_after_revert/logs/PORT-LIN-XTR-001.stderr.log`
- Verification:
  - With fallback definitions restored, `PORT-LIN-XTR-001` passes in `2026-04-19_alpine_tinygo_enabled`.

## TinyGo enablement notes on Alpine

- TinyGo binary path used:
  - `~/.cache/yq-portability/tools/tinygo/bin/tinygo`
- TinyGo validation output:
  - `Using tinygo: tinygo version 0.40.1 linux/amd64 (using go version go1.25.9 and LLVM version 20.1.1)`
- Runner invocation used to include TinyGo on PATH:
  - `PATH="$HOME/.cache/yq-portability/tools/tinygo/bin:$PATH" bash scripts/portability/run_portability_linux.sh 2026-04-19_alpine_tinygo_enabled`

## Generated evidence

- Latest run summary: `artifacts/portability/linux/2026-04-19_alpine_tinygo_enabled/summary.md`
- Latest machine-readable results: `artifacts/portability/linux/2026-04-19_alpine_tinygo_enabled/summary.csv`
- Latest per-case logs: `artifacts/portability/linux/2026-04-19_alpine_tinygo_enabled/logs/`
- Earlier Alpine run summaries:
  - `artifacts/portability/linux/2026-04-19_alpine_local/summary.md`
  - `artifacts/portability/linux/2026-04-19_alpine_rerun_after_revert/summary.md`
