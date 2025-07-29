#!/bin/bash
set -e

# Navigate to UI directory if not already there
[[ $(pwd) == *ui ]] || cd ui

# Install dependencies only if needed
if [ ! -d "node_modules" ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
  echo "Installing dependencies..."
  npm install
fi

echo "Running frontend tests..."
npm run test

COVERAGE_FILE="coverage/ui/coverage-summary.json"
THRESHOLD=70

if [ ! -f "$COVERAGE_FILE" ]; then
  echo "Coverage summary file not found"
  exit 1
fi

COVERAGE=$(node -e "console.log(require('./$COVERAGE_FILE').total.statements.pct)")

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "Frontend coverage: $COVERAGE% (below $THRESHOLD%)"
  exit 1
else
else
  echo "Frontend coverage: $COVERAGE%"
fi
