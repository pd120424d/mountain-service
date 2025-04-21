#!/bin/bash
set -e

cd ui
npm ci
npm run test -- --watch=false --code-coverage

COVERAGE_FILE="coverage/ui/coverage-summary.json"
if [ ! -f "$COVERAGE_FILE" ]; then
  echo "❌ Coverage summary file not found"
  exit 1
fi

COVERAGE=$(node -e "console.log(require('./$COVERAGE_FILE').total.statements.pct)")
THRESHOLD=55

echo "Total statement coverage: $COVERAGE%"
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "❌ Coverage $COVERAGE% is below threshold ($THRESHOLD%)."
  exit 1
else
  echo "✅ Coverage $COVERAGE% is above threshold ($THRESHOLD%)."
fi
