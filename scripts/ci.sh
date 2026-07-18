#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

qadc_required="${QADC_REQUIRED:-false}"
vuln_required="${VULN_REQUIRED:-false}"

ruby scripts/check-go-baseline.rb
ruby scripts/check-go-baseline_test.rb
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go test ./example -run 'Test(CoreBinDoesNotLoadServerDefaults|MultipleHTTPListenerExample|ServerHTTPProductionBoundaryExample|GRPCMetadataExample)' -count=1

if command -v qadc >/dev/null 2>&1; then
  output="$(GOWORK=off qadc gate --profile backend-service --format json)"
  printf 'ci_result step=qadc status=pass score=%s grade=%s\n' \
    "$(printf '%s' "$output" | jq -r '.quality_score.score')" \
    "$(printf '%s' "$output" | jq -r '.quality_score.grade')"
elif [[ "$qadc_required" == "true" ]]; then
  echo "ci_result step=qadc status=fail reason=qadc_missing" >&2
  exit 1
else
  echo "ci_result step=qadc status=coverage_degraded reason=qadc_missing"
fi

if command -v govulncheck >/dev/null 2>&1; then
  GOWORK=off govulncheck ./...
  echo "ci_result step=govulncheck status=pass"
elif [[ "$vuln_required" == "true" ]]; then
  echo "ci_result step=govulncheck status=fail reason=govulncheck_missing" >&2
  exit 1
else
  echo "ci_result step=govulncheck status=coverage_degraded reason=govulncheck_missing"
fi
