#!/usr/bin/env bash

set -euo pipefail

# Non-login shells (CI, some IDE terminals) may omit Go; try common locations.
if ! command -v go >/dev/null 2>&1; then
  for _go_dir in /usr/local/go/bin "${HOME}/go/bin" "${HOME}/.local/go/bin"; do
    if [[ -x "${_go_dir}/go" ]]; then
      export PATH="${_go_dir}:${PATH}"
      break
    fi
  done
fi
if ! command -v go >/dev/null 2>&1; then
  printf "error: go not found on PATH (install Go or add it to PATH, then retry).\n" >&2
  exit 2
fi

RUN_ID="${1:-$(date +"%Y-%m-%d_%H%M_local")}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

ARTIFACT_DIR="${ROOT_DIR}/artifacts/portability/linux/${RUN_ID}"
LOG_DIR="${ARTIFACT_DIR}/logs"
INCIDENT_DIR="${ARTIFACT_DIR}/incidents"
INPUT_DIR="${ARTIFACT_DIR}/inputs"
SUMMARY_CSV="${ARTIFACT_DIR}/summary.csv"
SUMMARY_MD="${ARTIFACT_DIR}/summary.md"

mkdir -p "${LOG_DIR}" "${INCIDENT_DIR}" "${INPUT_DIR}"

total=0
passed=0
failed=0
skipped=0
failed_case_ids=()

printf "case_id,status,exit_code,expected_exit,stdout_log,stderr_log\n" > "${SUMMARY_CSV}"

normalise_output() {
  local file_path="$1"
  tr -d '\r' < "${file_path}" | sed -e 's/[[:space:]]*$//'
}

write_incident_report() {
  local case_id="$1"
  local description="$2"
  local command_text="$3"
  local expected_text="$4"
  local actual_text="$5"
  local exit_code="$6"

  cat > "${INCIDENT_DIR}/${case_id}.md" <<EOF
# Incident: ${case_id}

## Case
- ID: ${case_id}
- Description: ${description}

## Environment
- Run ID: ${RUN_ID}
- Host: Linux
- Command: \`${command_text}\`

## Expected
${expected_text}

## Actual
${actual_text}

## Exit code
\`${exit_code}\`

## Severity
medium

## Reproduction
Run the same command from repository root and compare outputs in:
- \`${LOG_DIR}/${case_id}.stdout.log\`
- \`${LOG_DIR}/${case_id}.stderr.log\`
EOF
}

skip_case() {
  local case_id="$1"
  local description="$2"
  local reason="$3"
  local stdout_log="${LOG_DIR}/${case_id}.stdout.log"
  local stderr_log="${LOG_DIR}/${case_id}.stderr.log"

  total=$((total + 1))
  skipped=$((skipped + 1))

  printf "%s\n" "${reason}" > "${stdout_log}"
  : > "${stderr_log}"
  printf "%s,%s,%s,%s,%s,%s\n" "${case_id}" "SKIPPED" "NA" "NA" "${stdout_log}" "${stderr_log}" >> "${SUMMARY_CSV}"
  printf "[SKIPPED] %s - %s (%s)\n" "${case_id}" "${description}" "${reason}"
}

run_case() {
  local case_id="$1"
  local description="$2"
  local command_text="$3"
  local expected_exit="$4"
  local expected_output="${5:-}"
  local stdout_log="${LOG_DIR}/${case_id}.stdout.log"
  local stderr_log="${LOG_DIR}/${case_id}.stderr.log"
  local exit_code=0
  local status="PASS"
  local expected_text="Exit code ${expected_exit}"
  local actual_text="Exit code 0"

  total=$((total + 1))

  set +e
  (
    cd "${ROOT_DIR}" && eval "${command_text}"
  ) > "${stdout_log}" 2> "${stderr_log}"
  exit_code=$?
  set -e

  if [[ "${exit_code}" -ne "${expected_exit}" ]]; then
    status="FAIL"
    expected_text="Exit code ${expected_exit}"
    actual_text="Exit code ${exit_code}"
  fi

  if [[ "${status}" == "PASS" && -n "${expected_output}" ]]; then
    local actual_output
    actual_output="$(normalise_output "${stdout_log}")"
    if [[ "${actual_output}" != "${expected_output}" ]]; then
      status="FAIL"
      expected_text="Output: ${expected_output}"
      actual_text="Output: ${actual_output}"
    fi
  fi

  if [[ "${status}" == "PASS" ]]; then
    passed=$((passed + 1))
    printf "[PASS] %s - %s\n" "${case_id}" "${description}"
  else
    failed=$((failed + 1))
    failed_case_ids+=("${case_id}")
    write_incident_report "${case_id}" "${description}" "${command_text}" "${expected_text}" "${actual_text}" "${exit_code}"
    printf "[FAIL] %s - %s\n" "${case_id}" "${description}"
  fi

  printf "%s,%s,%s,%s,%s,%s\n" "${case_id}" "${status}" "${exit_code}" "${expected_exit}" "${stdout_log}" "${stderr_log}" >> "${SUMMARY_CSV}"
}

printf "a: 1\n" > "${INPUT_DIR}/basic.yml"
printf '{"name":"yq","ok":true}\n' > "${INPUT_DIR}/basic.json"
printf 'msg: "cafe 😊"\n' > "${INPUT_DIR}/utf8.yml"

# Linux environment capture (parallel to PORT-MAC-ENV-001 / sw_vers)
run_case "PORT-LIN-ENV-001" "Capture Linux environment details" "if test -f /etc/os-release; then cat /etc/os-release; fi && uname -a && go version && go env GOOS GOARCH && git rev-parse --short HEAD" 0
run_case "PORT-LIN-BLD-001" "Build yq packages locally" "go build ./..." 0
run_case "PORT-LIN-RT-001" "Run yq version command" "go build -o ./yq . && ./yq --version" 0
run_case "PORT-LIN-RT-002" "Evaluate YAML scalar" "./yq '.a' \"${INPUT_DIR}/basic.yml\"" 0 "1"
run_case "PORT-LIN-RT-003" "Extract JSON field" "./yq '.name' \"${INPUT_DIR}/basic.json\"" 0 "\"yq\""
run_case "PORT-LIN-RT-004" "Evaluate stdin pipeline" "printf 'k: v\n' | ./yq '.k' -" 0 "v"
run_case "PORT-LIN-RT-005" "Validate UTF-8 output" "./yq '.msg' \"${INPUT_DIR}/utf8.yml\"" 0 "cafe 😊"
run_case "PORT-LIN-TST-001" "Run project Go tests script" "./scripts/test.sh" 0
run_case "PORT-LIN-TST-002" "Run acceptance test script" "./scripts/acceptance.sh" 0
run_case "PORT-LIN-XTR-001" "Run small build profile" "./scripts/build-small-yq.sh" 0

if command -v tinygo >/dev/null 2>&1; then
  run_case "PORT-LIN-XTR-002" "Run TinyGo build profile" "./scripts/build-tinygo-yq.sh" 0
else
  skip_case "PORT-LIN-XTR-002" "Run TinyGo build profile" "tinygo is not installed"
fi

pass_rate="0.00"
if [[ "${total}" -gt 0 ]]; then
  pass_rate="$(awk "BEGIN { printf \"%.2f\", (${passed}/${total}) * 100 }")"
fi

{
  printf "# Linux Portability Run Summary\n\n"
  printf -- "- Run ID: \`%s\`\n" "${RUN_ID}"
  printf -- "- Artifacts: \`%s\`\n" "${ARTIFACT_DIR}"
  printf -- "- Total: %d\n" "${total}"
  printf -- "- Passed: %d\n" "${passed}"
  printf -- "- Failed: %d\n" "${failed}"
  printf -- "- Skipped: %d\n" "${skipped}"
  printf -- "- Pass rate: %s%%\n" "${pass_rate}"
  if [[ "${#failed_case_ids[@]}" -gt 0 ]]; then
    printf -- "- Failed case IDs: %s\n" "${failed_case_ids[*]}"
  else
    printf -- "- Failed case IDs: none\n"
  fi
} > "${SUMMARY_MD}"

printf "\nCompleted Linux portability run.\n"
printf "Summary: %s\n" "${SUMMARY_MD}"
printf "CSV: %s\n" "${SUMMARY_CSV}"

if [[ "${failed}" -gt 0 ]]; then
  exit 1
fi
