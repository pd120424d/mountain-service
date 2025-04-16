#!/bin/bash
set -e

COVERAGE_THRESHOLD=60
PKG="./employee"

cd api  # move into the module root (where go.mod is located)

echo "Running tests with coverage for $PKG..."

go test -coverprofile=coverage.out $PKG/...
go tool cover -func=coverage.out

COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "Total coverage: $COVERAGE%"

COVERAGE_INT=${COVERAGE%.*}

if [ "$COVERAGE_INT" -lt "$COVERAGE_THRESHOLD" ]; then
    echo "❌ Coverage $COVERAGE% is below threshold ($COVERAGE_THRESHOLD%). Failing."
    exit 1
fi

echo "✅ Coverage check passed."
