# Incident: PORT-WIN-TST-002

## Case
- ID: PORT-WIN-TST-002
- Description: Run acceptance test script

## Environment
- Run ID: 2026-04-15_2104_local
- Host: Windows
- Command: `bash ./scripts/acceptance.sh`

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
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2104_local\logs\PORT-WIN-TST-002.stdout.log`
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2104_local\logs\PORT-WIN-TST-002.stderr.log`
