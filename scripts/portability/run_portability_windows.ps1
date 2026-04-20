param(
  [string]$RunId = (Get-Date -Format "yyyy-MM-dd_HHmm_local")
)

$ErrorActionPreference = "Stop"
$PSNativeCommandUseErrorActionPreference = $false
[Console]::InputEncoding = [System.Text.UTF8Encoding]::new($false)
[Console]::OutputEncoding = [System.Text.UTF8Encoding]::new($false)
$OutputEncoding = [Console]::OutputEncoding

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Resolve-Path (Join-Path $scriptDir "..\..")
$bt = [char]96

$preferredToolDirs = @(
  "C:\Program Files\Go\bin",
  "C:\Program Files\Git\cmd",
  "C:\Program Files\Git\bin"
)
foreach ($toolDir in $preferredToolDirs) {
  if ((Test-Path $toolDir) -and -not (($env:Path -split ";") -contains $toolDir)) {
    $env:Path = "$toolDir;$env:Path"
  }
}

$artifactDir = Join-Path $rootDir "artifacts\portability\windows\$RunId"
$logDir = Join-Path $artifactDir "logs"
$incidentDir = Join-Path $artifactDir "incidents"
$inputDir = Join-Path $artifactDir "inputs"
$summaryCsv = Join-Path $artifactDir "summary.csv"
$summaryMd = Join-Path $artifactDir "summary.md"

New-Item -ItemType Directory -Force -Path $logDir, $incidentDir, $inputDir | Out-Null

$total = 0
$passed = 0
$failed = 0
$skipped = 0
$failedCaseIds = [System.Collections.Generic.List[string]]::new()

"case_id,status,exit_code,expected_exit,stdout_log,stderr_log" | Set-Content -Path $summaryCsv -Encoding UTF8

function Write-Utf8NoBomFile {
  param(
    [string]$Path,
    [string]$Content
  )
  [System.IO.File]::WriteAllText($Path, $Content, [System.Text.UTF8Encoding]::new($false))
}

function Normalise-Output {
  param([string]$FilePath)
  $raw = Get-Content -Path $FilePath -Raw -Encoding UTF8
  $raw = $raw -replace "`r", ""
  $lines = $raw -split "`n" | ForEach-Object { $_.TrimEnd() }
  while ($lines.Count -gt 0 -and $lines[-1] -eq "") {
    $lines = $lines[0..($lines.Count - 2)]
  }
  return ($lines -join "`n")
}

function Write-IncidentReport {
  param(
    [string]$CaseId,
    [string]$Description,
    [string]$CommandText,
    [string]$ExpectedText,
    [string]$ActualText,
    [int]$ExitCode
  )

  @"
# Incident: $CaseId

## Case
- ID: $CaseId
- Description: $Description

## Environment
- Run ID: $RunId
- Host: Windows
- Command: $bt$CommandText$bt

## Expected
$ExpectedText

## Actual
$ActualText

## Exit code
$bt$ExitCode$bt

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- $bt$logDir\$CaseId.stdout.log$bt
- $bt$logDir\$CaseId.stderr.log$bt
"@ | Set-Content -Path (Join-Path $incidentDir "$CaseId.md") -Encoding UTF8
}

function Append-SummaryRow {
  param(
    [string]$CaseId,
    [string]$Status,
    [string]$ExitCode,
    [string]$ExpectedExit,
    [string]$StdoutLog,
    [string]$StderrLog
  )
  "$CaseId,$Status,$ExitCode,$ExpectedExit,$StdoutLog,$StderrLog" | Add-Content -Path $summaryCsv -Encoding UTF8
}

function Skip-Case {
  param(
    [string]$CaseId,
    [string]$Description,
    [string]$Reason
  )

  $script:total++
  $script:skipped++
  $stdoutLog = Join-Path $logDir "$CaseId.stdout.log"
  $stderrLog = Join-Path $logDir "$CaseId.stderr.log"

  $Reason | Set-Content -Path $stdoutLog -Encoding UTF8
  "" | Set-Content -Path $stderrLog -Encoding UTF8
  Append-SummaryRow $CaseId "SKIPPED" "NA" "NA" $stdoutLog $stderrLog
  Write-Host "[SKIPPED] $CaseId - $Description ($Reason)"
}

function Run-Case {
  param(
    [string]$CaseId,
    [string]$Description,
    [string]$CommandText,
    [int]$ExpectedExit,
    [scriptblock]$Command,
    [string]$ExpectedOutput = ""
  )

  $script:total++
  $stdoutLog = Join-Path $logDir "$CaseId.stdout.log"
  $stderrLog = Join-Path $logDir "$CaseId.stderr.log"
  $status = "PASS"
  $expectedText = "Exit code $ExpectedExit"
  $actualText = "Exit code 0"
  $exitCode = 0

  Push-Location $rootDir
  try {
    $previousEap = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
      & $Command 1> $stdoutLog 2> $stderrLog
      $exitCode = if ($null -ne $LASTEXITCODE) { $LASTEXITCODE } else { 0 }
    } catch {
      $_ | Out-String | Set-Content -Path $stderrLog -Encoding UTF8
      $exitCode = if ($null -ne $LASTEXITCODE -and $LASTEXITCODE -ne 0) { $LASTEXITCODE } else { 1 }
    } finally {
      $ErrorActionPreference = $previousEap
    }
  } finally {
    Pop-Location
  }

  if ($exitCode -ne $ExpectedExit) {
    $status = "FAIL"
    $actualText = "Exit code $exitCode"
  }

  if ($status -eq "PASS" -and $ExpectedOutput -ne "") {
    $actualOutput = Normalise-Output -FilePath $stdoutLog
    if ($actualOutput -ne $ExpectedOutput) {
      $status = "FAIL"
      $expectedText = "Output: $ExpectedOutput"
      $actualText = "Output: $actualOutput"
    }
  }

  if ($status -eq "PASS") {
    $script:passed++
    Write-Host "[PASS] $CaseId - $Description"
  } else {
    $script:failed++
    $failedCaseIds.Add($CaseId) | Out-Null
    Write-IncidentReport $CaseId $Description $CommandText $expectedText $actualText $exitCode
    Write-Host "[FAIL] $CaseId - $Description"
  }

  Append-SummaryRow $CaseId $status "$exitCode" "$ExpectedExit" $stdoutLog $stderrLog
}

Write-Utf8NoBomFile -Path (Join-Path $inputDir "basic.yml") -Content "a: 1`n"
Write-Utf8NoBomFile -Path (Join-Path $inputDir "basic.json") -Content '{"name":"yq","ok":true}'
Write-Utf8NoBomFile -Path (Join-Path $inputDir "utf8.yml") -Content 'msg: "cafe 😊"'

$hasBash = $null -ne (Get-Command bash -ErrorAction SilentlyContinue)

Run-Case "PORT-WIN-ENV-001" "Capture Windows environment details" "cmd /c ver; go version; go env GOOS GOARCH; git rev-parse --short HEAD" 0 {
  cmd /c ver
  go version
  go env GOOS GOARCH
  git rev-parse --short HEAD
}

Run-Case "PORT-WIN-BLD-001" "Build yq packages locally" "go build ./..." 0 {
  go build ./...
}

Run-Case "PORT-WIN-RT-001" "Run yq version command" "go build -o .\yq.exe .; .\yq.exe --version" 0 {
  go build -o .\yq.exe .
  if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
  }
  .\yq.exe --version
}

Run-Case "PORT-WIN-RT-002" "Evaluate YAML scalar" ".\yq.exe '.a' <basic.yml>" 0 {
  .\yq.exe '.a' (Join-Path $inputDir "basic.yml")
} "1"

Run-Case "PORT-WIN-RT-003" "Extract JSON field" ".\yq.exe '.name' <basic.json>" 0 {
  .\yq.exe '.name' (Join-Path $inputDir "basic.json")
} '"yq"'

Run-Case "PORT-WIN-RT-004" "Evaluate stdin pipeline" "'k: v' | .\yq.exe '.k' -" 0 {
  "k: v`n" | .\yq.exe '.k' -
} "v"

Run-Case "PORT-WIN-RT-005" "Validate UTF-8 output" ".\yq.exe '.msg' <utf8.yml>" 0 {
  .\yq.exe '.msg' (Join-Path $inputDir "utf8.yml")
} "cafe 😊"

if ($hasBash) {
  Run-Case "PORT-WIN-TST-001" "Run project Go tests script" "bash ./scripts/test.sh" 0 {
    bash ./scripts/test.sh
  }
  Run-Case "PORT-WIN-TST-002" "Run acceptance test script" "bash ./scripts/acceptance.sh" 0 {
    bash ./scripts/acceptance.sh
  }
  Run-Case "PORT-WIN-XTR-001" "Run small build profile" "bash ./scripts/build-small-yq.sh" 0 {
    bash ./scripts/build-small-yq.sh
  }
} else {
  Skip-Case "PORT-WIN-TST-001" "Run project Go tests script" "bash is not installed"
  Skip-Case "PORT-WIN-TST-002" "Run acceptance test script" "bash is not installed"
  Skip-Case "PORT-WIN-XTR-001" "Run small build profile" "bash is not installed"
}

if ($null -ne (Get-Command tinygo -ErrorAction SilentlyContinue)) {
  if ($hasBash) {
    Run-Case "PORT-WIN-XTR-002" "Run TinyGo build profile" "bash ./scripts/build-tinygo-yq.sh" 0 {
      bash ./scripts/build-tinygo-yq.sh
    }
  } else {
    Skip-Case "PORT-WIN-XTR-002" "Run TinyGo build profile" "tinygo is installed but bash is not installed"
  }
} else {
  Skip-Case "PORT-WIN-XTR-002" "Run TinyGo build profile" "tinygo is not installed"
}

$passRate = "0.00"
if ($total -gt 0) {
  $passRate = "{0:N2}" -f (($passed / $total) * 100)
}

if ($failedCaseIds.Count -gt 0) {
  $failedCaseIdsText = ($failedCaseIds -join " ")
} else {
  $failedCaseIdsText = "none"
}

@"
# Windows Portability Run Summary

- Run ID: $bt$RunId$bt
- Artifacts: $bt$artifactDir$bt
- Total: $total
- Passed: $passed
- Failed: $failed
- Skipped: $skipped
- Pass rate: $passRate%
- Failed case IDs: $failedCaseIdsText
"@ | Set-Content -Path $summaryMd -Encoding UTF8

Write-Host ""
Write-Host "Completed Windows portability run."
Write-Host "Summary: $summaryMd"
Write-Host "CSV: $summaryCsv"

if ($failed -gt 0) {
  exit 1
}
