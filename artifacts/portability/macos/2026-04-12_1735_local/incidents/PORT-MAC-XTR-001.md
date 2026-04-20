# Incident: PORT-MAC-XTR-001

## Case
- ID: PORT-MAC-XTR-001
- Description: Run small build profile

## Environment
- Run ID: 2026-04-12_1735_local
- Host: macOS
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
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1735_local/logs/PORT-MAC-XTR-001.stdout.log`
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1735_local/logs/PORT-MAC-XTR-001.stderr.log`
