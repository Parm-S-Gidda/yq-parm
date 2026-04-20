# Arch Linux Portability Run Report

## Latest run metadata

- Run ID: `2026-04-19_1414_arch_local`
- Date: 2026-04-19
- Host: Arch Linux (`arch`, rolling)
- Kernel/arch: `Linux 6.19.11-arch1-1`, `amd64`
- Go version (host): `go1.26.2-X:nodwarf5`
- TinyGo path used for pass run: `~/.cache/yq-portability/tools/tinygo/bin/tinygo`
- TinyGo version (pass run): `0.40.1` (LLVM `20.1.1`)
- TinyGo Go toolchain path (pass run): `~/sdk/go1.25.9/bin/go`
- Script: `scripts/portability/run_portability_linux.sh`
- Artifacts directory: `artifacts/portability/linux/2026-04-19_1414_arch_local/`

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

## Progress across Arch runs

- Initial run (`2026-04-19_1354_arch_local`):
  - pass rate: `81.82%`
  - failed case: `PORT-LIN-XTR-001`
  - skipped case: `PORT-LIN-XTR-002` (`tinygo` not installed)
- Intermediate run after restoring `no_kyaml` fallback (`2026-04-19_1400_arch_local`):
  - pass rate: `90.91%`
  - `PORT-LIN-XTR-001` passes
  - `PORT-LIN-XTR-002` still skipped (`tinygo` not installed)
- Current run with TinyGo enabled (`2026-04-19_1414_arch_local`):
  - pass rate: `100.00%`
  - all `PORT-LIN-*` cases pass, including `PORT-LIN-XTR-002`

## Root cause confirmation (`PORT-LIN-XTR-001`)

- Classification: **legitimate code defect**, reproducible on Arch before fallback symbols are present.
- Why:
  - `./scripts/build-small-yq.sh` compiles with `yq_nokyaml`.
  - Under `yq_nokyaml`, `pkg/yqlib/kyaml.go` is excluded.
  - Without fallback definitions in `pkg/yqlib/no_kyaml.go`, references fail during compile.
- Error observed:
  - `undefined: KYamlPreferences`
  - `undefined: ConfiguredKYamlPreferences`
- Verification:
  - with fallback definitions restored in `pkg/yqlib/no_kyaml.go`, `PORT-LIN-XTR-001` passes on Arch.

## TinyGo enablement notes on Arch

- Arch repo TinyGo (`0.37.0`) is older and does not match this repo baseline (`go 1.25` in `go.mod`), so the default `tinygo` command is not sufficient for this test matrix.
- A local TinyGo `0.40.1` installation under `~/.cache/yq-portability/tools/tinygo/` was used for the passing run.
- TinyGo-compatible Go wrapper used for the pass run:
  - `YQ_TINYGO_GO_BIN="$HOME/sdk/go1.25.9/bin/go"`
- Passing invocation:
  - `PATH="$HOME/.cache/yq-portability/tools/tinygo/bin:$PATH" YQ_TINYGO_GO_BIN="$HOME/sdk/go1.25.9/bin/go" ./scripts/portability/run_portability_linux.sh 2026-04-19_1414_arch_local`

## Generated evidence

- Latest run summary: `artifacts/portability/linux/2026-04-19_1414_arch_local/summary.md`
- Latest machine-readable results: `artifacts/portability/linux/2026-04-19_1414_arch_local/summary.csv`
- Latest per-case logs: `artifacts/portability/linux/2026-04-19_1414_arch_local/logs/`
- Earlier Arch run summaries:
  - `artifacts/portability/linux/2026-04-19_1354_arch_local/summary.md`
  - `artifacts/portability/linux/2026-04-19_1400_arch_local/summary.md`
