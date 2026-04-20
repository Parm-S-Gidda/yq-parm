# Incident: PORT-WIN-XTR-001

## Case
- ID: PORT-WIN-XTR-001
- Description: Run small build profile

## Environment
- Run ID: 2026-04-15_2058_local
- Host: Windows
- Command: `bash ./scripts/build-small-yq.sh`

## Expected
Exit code 0

## Actual
Exit code -1

## Exit code
`-1`

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2058_local\logs\PORT-WIN-XTR-001.stdout.log`
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2058_local\logs\PORT-WIN-XTR-001.stderr.log`
