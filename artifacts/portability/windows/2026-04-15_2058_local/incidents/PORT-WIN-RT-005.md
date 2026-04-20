# Incident: PORT-WIN-RT-005

## Case
- ID: PORT-WIN-RT-005
- Description: Validate UTF-8 output

## Environment
- Run ID: 2026-04-15_2058_local
- Host: Windows
- Command: `.\yq.exe '.msg' <utf8.yml>`

## Expected
Output: cafe ðŸ˜Š

## Actual
Output: cafe ├░┼╕╦£┼á

## Exit code
`0`

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2058_local\logs\PORT-WIN-RT-005.stdout.log`
- `C:\Users\tqc\Downloads\cmpt473\yq-fork\artifacts\portability\windows\2026-04-15_2058_local\logs\PORT-WIN-RT-005.stderr.log`
