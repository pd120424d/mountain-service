#!/bin/bash
set -e

# Install go-acc if needed
if ! command -v go-acc &> /dev/null; then
  echo "Installing go-acc..."
  go install github.com/ory/go-acc@latest
  export PATH="$PATH:$(go env GOPATH)/bin"
fi

# Move into the api folder
cd api

SERVICES=("employee" "urgency")
THRESHOLD=80.0
OVERALL_SUCCESS=true

echo "Running coverage tests for all services..."
echo "Coverage threshold: $THRESHOLD%"
echo "=================================="

for SERVICE in "${SERVICES[@]}"; do
  echo ""
  echo "Testing service: $SERVICE"
  echo "----------------------------"

  # Check if service has internal packages
  TARGETS=$(go list ./$SERVICE/internal/... 2>/dev/null || echo "")

  if [ -z "$TARGETS" ]; then
    echo "WARNING: No internal packages found for $SERVICE, skipping..."
    continue
  fi

  echo "Running coverage for:"
  echo "$TARGETS"

  # Run go-acc with only those packages
  COVERAGE_FILE="coverage-$SERVICE.out"
  if go-acc $TARGETS \
    --ignore ".*_gomock.go" \
    --output "$COVERAGE_FILE"; then

    # Show coverage summary without gomock
    echo ""
    echo "Coverage details for $SERVICE:"
    go tool cover -func="$COVERAGE_FILE" | grep -v '_gomock.go'

    # Check if coverage is above threshold
    COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print substr($3, 1, length($3)-1)}')

    echo ""
    echo "📊 $SERVICE coverage: $COVERAGE%"

    if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
      echo "❌ $SERVICE coverage $COVERAGE% is below threshold ($THRESHOLD%)"
      OVERALL_SUCCESS=false
    else
      echo "✅ $SERVICE coverage $COVERAGE% meets threshold ($THRESHOLD%)"
    fi

    # Clean up individual coverage file
    rm -f "$COVERAGE_FILE"
  else
    echo "❌ Failed to run coverage for $SERVICE"
    OVERALL_SUCCESS=false
  fi

  echo "----------------------------"
done

echo ""
echo "=================================="
if [ "$OVERALL_SUCCESS" = true ]; then
  echo "SUCCESS: All services passed coverage requirements!"
  exit 0
else
  echo "FAILURE: One or more services failed coverage requirements!"
  exit 1
fi
