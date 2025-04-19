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

# Define the list of packages to test
TARGETS=$(go list ./employee/internal/... 2>/dev/null)

echo "Running coverage for:"
echo "$TARGETS"

# Run go-acc with only those packages
go-acc $TARGETS \
  --ignore ".*_gomock.go" \
  --output coverage.out

# Show coverage summary without gomock
go tool cover -func=coverage.out | grep -v '_gomock.go'

# Check if coverage is above threshold
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
THRESHOLD=80.0

echo "Total coverage: $COVERAGE%"
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "❌ Coverage $COVERAGE% is below threshold ($THRESHOLD%). Failing."
  exit 1
else
  echo "✅ Coverage $COVERAGE% is above threshold ($THRESHOLD%)."
fi
