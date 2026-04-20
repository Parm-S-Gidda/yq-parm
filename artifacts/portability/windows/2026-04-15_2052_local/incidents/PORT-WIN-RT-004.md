# Incident: PORT-WIN-RT-004

## Case
- ID: PORT-WIN-RT-004
- Description: Evaluate stdin pipeline

## Environment
- Run ID: 2026-04-15_2052_local
- Host: Windows
- Command: `'k: v' | .\yq.exe '.k' -`

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
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2052_local\logs\PORT-WIN-RT-004.stdout.log`
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2052_local\logs\PORT-WIN-RT-004.stderr.log`
