# Incident: PORT-MAC-RT-003

## Case
- ID: PORT-MAC-RT-003
- Description: Extract JSON field

## Environment
- Run ID: 2026-04-12_1735_local
- Host: macOS
- Command: `./yq '.name' "/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1735_local/inputs/basic.json"`

## Expected
Output: yq

## Actual
Output: "yq"

## Exit code
`0`

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1735_local/logs/PORT-MAC-RT-003.stdout.log`
- `/Users/tran/Desktop/CMPT 473/yq/artifacts/portability/macos/2026-04-12_1735_local/logs/PORT-MAC-RT-003.stderr.log`
