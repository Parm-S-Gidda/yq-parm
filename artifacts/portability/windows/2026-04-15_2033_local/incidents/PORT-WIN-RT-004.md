# Incident: PORT-WIN-RT-004

## Case
- ID: PORT-WIN-RT-004
- Description: Evaluate stdin pipeline

## Environment
- Run ID: 2026-04-15_2033_local
- Host: Windows
- Command: $CommandText

## Expected
Exit code 0

## Actual
Exit code 1

## Exit code
$ExitCode

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- $logDir\PORT-WIN-RT-004.stdout.log
- $logDir\PORT-WIN-RT-004.stderr.log
