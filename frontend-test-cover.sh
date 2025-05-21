#!/bin/bash
set -e

if pwd | grep -q ui; then
  echo "Already in ui directory"
else
  echo "Navigating to ui directory"
  cd ui
fi

npm ci
npm run test -- --watch=false --code-coverage

COVERAGE_FILE="coverage/ui/coverage-summary.json"
if [ ! -f "$COVERAGE_FILE" ]; then
  echo "❌ Coverage summary file not found"
  exit 1
fi

COVERAGE=$(node -e "console.log(require('./$COVERAGE_FILE').total.statements.pct)")
THRESHOLD=50

echo "Total statement coverage: $COVERAGE%"
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "❌ Coverage $COVERAGE% is below threshold ($THRESHOLD%)."
  exit 1
else
  echo "✅ Coverage $COVERAGE% is above threshold ($THRESHOLD%)."
fi
