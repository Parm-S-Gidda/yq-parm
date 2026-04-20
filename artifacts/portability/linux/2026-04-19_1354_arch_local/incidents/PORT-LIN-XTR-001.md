# Incident: PORT-LIN-XTR-001

## Case
- ID: PORT-LIN-XTR-001
- Description: Run small build profile

## Environment
- Run ID: 2026-04-19_1354_arch_local
- Host: Linux
- Command: `./scripts/build-small-yq.sh`

## Expected
Exit code 0

## Actual
Exit code 1

## Exit code
`1`

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- `/home/tqc/dev/cmpt473/yq-fork/artifacts/portability/linux/2026-04-19_1354_arch_local/logs/PORT-LIN-XTR-001.stdout.log`
- `/home/tqc/dev/cmpt473/yq-fork/artifacts/portability/linux/2026-04-19_1354_arch_local/logs/PORT-LIN-XTR-001.stderr.log`
