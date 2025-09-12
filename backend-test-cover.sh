#!/bin/bash
set -e

# Default to pure-Go builds to avoid CGO flakiness (e.g., mattn/go-sqlite3 compile)
export CGO_ENABLED=${CGO_ENABLED:-0}

# Install go-acc if needed
command -v go-acc >/dev/null || {
  echo "Installing go-acc..."
  go install github.com/ory/go-acc@latest
  export PATH="$PATH:$(go env GOPATH)/bin"
}

cd api

SERVICES=("activity" "activity-readmodel-updater" "employee" "urgency") # Add "activity-readmodel-updater" when the tests are implemented and have enough coverage
THRESHOLD=75.0
OVERALL_SUCCESS=true

echo "Running backend coverage tests (threshold: $THRESHOLD%)"
echo "=================================================="

for SERVICE in "${SERVICES[@]}"; do
  echo "Testing $SERVICE..."

    TARGETS=$(go list ./$SERVICE/internal/... 2>/dev/null || echo "")
    [ -z "$TARGETS" ] && { echo "WARNING: No packages found for $SERVICE, skipping"; continue; }
    if [ -n "${VERBOSE:-}" ]; then
      echo "Targets: $TARGETS"
    fi

    TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-600}"
    COVERAGE_FILE="coverage-$SERVICE.out"
    if [ -n "${VERBOSE:-}" ]; then
      echo "Running: go-acc $TARGETS --ignore '.*_gomock.go' --output '$COVERAGE_FILE' -- -v (timeout ${TIMEOUT_SECONDS}s)"
      timeout "${TIMEOUT_SECONDS}s" go-acc $TARGETS --ignore ".*_gomock.go" --output "$COVERAGE_FILE" -- -v
    else
      timeout "${TIMEOUT_SECONDS}s" go-acc $TARGETS --ignore ".*_gomock.go" --output "$COVERAGE_FILE" >/dev/null 2>&1
    fi

    if [ -f "$COVERAGE_FILE" ]; then
      COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print substr($3, 1, length($3)-1)}')

      if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
        echo "FAILURE: $SERVICE: $COVERAGE% (below $THRESHOLD%)"
        OVERALL_SUCCESS=false
      else
        echo "SUCCESS: $SERVICE: $COVERAGE%"
      fi
      rm -f "$COVERAGE_FILE"
    else
      echo "FAILURE: $SERVICE: Failed to run tests"
      OVERALL_SUCCESS=false
    fi
done

echo "=================================================="
if [ "$OVERALL_SUCCESS" = true ]; then
  echo "SUCCESS: All backend services passed coverage requirements!"
  exit 0
else
  echo "FAILURE: Some services failed coverage requirements!"
  exit 1
fi
