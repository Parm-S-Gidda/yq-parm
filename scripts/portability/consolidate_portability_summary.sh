#!/usr/bin/env bash
# Print a cross-OS markdown table from portability summary.md files.
# Usage:
#   bash scripts/portability/consolidate_portability_summary.sh <mac_summary.md> [linux_summary.md] [windows_summary.md]
# Any omitted path is skipped (use "" to skip an OS while passing later paths).

set -euo pipefail

extract_field() {
  local file="$1"
  local label="$2"
  grep -E "^- ${label}:" "${file}" | head -1 | sed -E "s/^- ${label}:[[:space:]]*//"
}

print_row() {
  local os="$1"
  local file="$2"
  if [[ -z "${file}" || ! -f "${file}" ]]; then
    printf "| %s | — | — | — | — | — |\n" "${os}"
    return
  fi
  local total passed failed skipped rate
  total="$(extract_field "${file}" "Total")"
  passed="$(extract_field "${file}" "Passed")"
  failed="$(extract_field "${file}" "Failed")"
  skipped="$(extract_field "${file}" "Skipped")"
  rate="$(extract_field "${file}" "Pass rate")"
  printf "| %s | %s | %s | %s | %s | %s |\n" "${os}" "${total}" "${passed}" "${failed}" "${skipped}" "${rate}"
}

if [[ "${#}" -lt 1 ]]; then
  printf "usage: %s <mac_summary.md> [linux_summary.md] [windows_summary.md]\n" "$(basename "$0")" >&2
  exit 2
fi

printf "## Portability summary (by OS)\n\n"
printf "| OS | Total | Passed | Failed | Skipped | Pass rate |\n"
printf "|----|-------|--------|--------|---------|----------|\n"
print_row "macOS" "${1:-}"
print_row "Linux" "${2:-}"
print_row "Windows" "${3:-}"
printf "\n"
