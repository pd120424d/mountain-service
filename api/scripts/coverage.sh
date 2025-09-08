#!/bin/bash

# Backend Coverage Script for Mountain Service
# Runs Go tests with coverage for all backend services, excluding generated files
# Usage: ./backend-coverage [service_name|all] [--html] [--threshold=N]

set -e

# Force pure-Go builds during coverage to avoid CGO toolchain stalls (e.g., sqlite3 CGO)
export CGO_ENABLED=${CGO_ENABLED:-0}

# Default values
SERVICE=""
GENERATE_HTML=false
THRESHOLD=75
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --html)
            GENERATE_HTML=true
            shift
            ;;
        --threshold=*)
            THRESHOLD="${1#*=}"
            shift
            ;;
        --help|-h)
            echo "Backend Coverage Script for Mountain Service"
            echo ""
            echo "Usage: $0 [service_name|all] [options]"
            echo ""
            echo "Services:"
            echo "  activity              Run coverage for activity service"
            echo "  urgency               Run coverage for urgency service"
            echo "  employee              Run coverage for employee service"
            # echo "  activity-readmodel    Run coverage for activity-readmodel-updater"
            echo "  all                   Run coverage for all services (default)"
            echo ""
            echo "Options:"
            echo "  --html                Generate HTML coverage report"
            echo "  --threshold=N         Set coverage threshold (default: 75)"
            echo "  --help, -h            Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 activity --html"
            echo "  $0 all --threshold=80"
            echo "  $0 urgency"
            exit 0
            ;;
        *)
            if [[ -z "$SERVICE" ]]; then
                SERVICE="$1"
            fi
            shift
            ;;
    esac
done

# Default to 'all' if no service specified
SERVICE=${SERVICE:-"all"}

echo "Mountain Service Backend Coverage Report"
echo "============================================="

# Function to run coverage for a single service
run_service_coverage() {
    local service_name=$1
    local service_path=$2

    echo ""
    echo "Testing $service_name service..."
    echo "-----------------------------------"

    if [[ ! -d "$service_path" ]]; then
        echo "Service directory not found: $service_path"
        return 1
    fi

    cd "$service_path"

    local coverage_file="coverage.out"
    local coverage_filtered="coverage.filtered.out"
    local coverage_html="coverage.html"

    # Run tests and generate coverage profile
    echo "Running tests..."
    if ! go test ./... -coverprofile=$coverage_file -coverpkg=./...; then
        echo "Tests failed for $service_name"
        cd "$PROJECT_ROOT"
        return 1
    fi

    if [[ ! -f "$coverage_file" ]]; then
        echo "No coverage file generated for $service_name"
        cd "$PROJECT_ROOT"
        return 1
    fi

    # Filter out generated files
    echo "Filtering out generated files (mocks, main.go, cmd/, docs/)..."
    grep -v "_gomock.go" $coverage_file | grep -v "main.go" | grep -v "/cmd/" | grep -v "/docs/" > $coverage_filtered || cp $coverage_file $coverage_filtered

    if [[ ! -s "$coverage_filtered" ]]; then
        echo "No coverage data after filtering, using original"
        cp $coverage_file $coverage_filtered
    fi

    # Generate coverage report
    echo ""
    echo "Coverage report for $service_name:"
    go tool cover -func=$coverage_filtered

    # Calculate coverage percentage
    local coverage_percent=$(go tool cover -func=$coverage_filtered | grep "total:" | awk '{print $3}' | sed 's/%//')

    if [[ -z "$coverage_percent" ]]; then
        echo "Could not calculate coverage for $service_name"
        cd "$PROJECT_ROOT"
        return 1
    fi

    echo ""
    echo "$service_name Coverage: ${coverage_percent}%"

    # Generate HTML report if requested
    if [[ "$GENERATE_HTML" == "true" ]]; then
        go tool cover -html=$coverage_filtered -o $coverage_html
        echo "HTML report: $service_path/$coverage_html"
    fi

    # Check threshold
    if (( $(echo "$coverage_percent >= $THRESHOLD" | bc -l) )); then
        echo "$service_name meets threshold: ${coverage_percent}% >= ${THRESHOLD}%"
        cd "$PROJECT_ROOT"
        return 0
    else
        echo "$service_name below threshold: ${coverage_percent}% < ${THRESHOLD}%"
        cd "$PROJECT_ROOT"
        return 1
    fi
}

# Define available services
declare -A SERVICES
SERVICES[activity]="$PROJECT_ROOT/activity"
SERVICES[urgency]="$PROJECT_ROOT/urgency"
SERVICES[employee]="$PROJECT_ROOT/employee"
# SERVICES[activity-readmodel]="$PROJECT_ROOT/activity-readmodel-updater"

# Main execution logic
FAILED_SERVICES=()
PASSED_SERVICES=()

if [[ "$SERVICE" == "all" ]]; then
    echo "Running coverage for all backend services..."
    echo "Threshold: ${THRESHOLD}%"

    for service_name in "${!SERVICES[@]}"; do
        if run_service_coverage "$service_name" "${SERVICES[$service_name]}"; then
            PASSED_SERVICES+=("$service_name")
        else
            FAILED_SERVICES+=("$service_name")
        fi
    done
else
    if [[ -n "${SERVICES[$SERVICE]}" ]]; then
        echo "Running coverage for $SERVICE service..."
        echo "Threshold: ${THRESHOLD}%"

        if run_service_coverage "$SERVICE" "${SERVICES[$SERVICE]}"; then
            PASSED_SERVICES+=("$SERVICE")
        else
            FAILED_SERVICES+=("$SERVICE")
        fi
    else
        echo "Unknown service: $SERVICE"
        echo "Available services: ${!SERVICES[*]}"
        exit 1
    fi
fi

# Final summary
echo ""
echo "Final Results"
echo "================"

if [[ ${#PASSED_SERVICES[@]} -gt 0 ]]; then
    echo "Passed (${#PASSED_SERVICES[@]}): ${PASSED_SERVICES[*]}"
fi

if [[ ${#FAILED_SERVICES[@]} -gt 0 ]]; then
    echo "Failed (${#FAILED_SERVICES[@]}): ${FAILED_SERVICES[*]}"
    echo ""
    echo "Tip: Use --html flag to generate detailed HTML reports"
    exit 1
else
    echo ""
    echo "All services meet the coverage threshold of ${THRESHOLD}%!"
    exit 0
fi
