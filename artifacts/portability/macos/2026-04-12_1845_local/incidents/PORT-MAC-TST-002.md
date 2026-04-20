# Incident: PORT-MAC-TST-002

## Case
- ID: PORT-MAC-TST-002
- Description: Run acceptance test script

## Environment
- Run ID: 2026-04-12_1845_local
- Host: macOS
- Command: `./scripts/acceptance.sh`

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
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1845_local/logs/PORT-MAC-TST-002.stdout.log`
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1845_local/logs/PORT-MAC-TST-002.stderr.log`
