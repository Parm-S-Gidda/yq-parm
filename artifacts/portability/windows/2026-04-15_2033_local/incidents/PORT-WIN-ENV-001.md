# Incident: PORT-WIN-ENV-001

## Case
- ID: PORT-WIN-ENV-001
- Description: Capture Windows environment details

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
- $logDir\PORT-WIN-ENV-001.stdout.log
- $logDir\PORT-WIN-ENV-001.stderr.log
