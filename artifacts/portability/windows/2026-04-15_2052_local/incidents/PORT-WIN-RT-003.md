# Incident: PORT-WIN-RT-003

## Case
- ID: PORT-WIN-RT-003
- Description: Extract JSON field

## Environment
- Run ID: 2026-04-15_2052_local
- Host: Windows
- Command: `.\yq.exe '.name' <basic.json>`

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
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2052_local\logs\PORT-WIN-RT-003.stdout.log`
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2052_local\logs\PORT-WIN-RT-003.stderr.log`
