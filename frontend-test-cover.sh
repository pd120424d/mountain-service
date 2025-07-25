#!/bin/bash
set -e

if pwd | grep -q ui; then
  echo "Already in ui directory"
else
  echo "Navigating to ui directory"
  cd ui
fi

# Skip npm install if node_modules exists and has packages
if [ ! -d "node_modules" ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
  echo "Installing dependencies..."
  npm install
else
  echo "Dependencies already installed, skipping npm install"
fi

npm run test

COVERAGE_FILE="coverage/ui/coverage-summary.json"
if [ ! -f "$COVERAGE_FILE" ]; then
  echo "[FAILURE] Coverage summary file not found"
  exit 1
fi

COVERAGE=$(node -e "console.log(require('./$COVERAGE_FILE').total.statements.pct)")
THRESHOLD=70

echo "Total statement coverage: $COVERAGE%"
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "[FAILURE] Coverage $COVERAGE% is below threshold ($THRESHOLD%)."
  exit 1
else
  echo "[SUCCESS]Coverage $COVERAGE% is above threshold ($THRESHOLD%)."
fi
