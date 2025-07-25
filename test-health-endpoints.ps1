#!/usr/bin/env pwsh

# Test script to verify all health endpoints work correctly
# This script tests the health endpoints that were renamed from /ping to /api/v1/health

Write-Host "Testing Health Endpoints" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green

# Define the services and their ports
$services = @{
    "employee-service" = 8082
    "urgency-service" = 8083
    "activity-service" = 8084
    "version-service" = 8090
}

$allHealthy = $true

foreach ($service in $services.GetEnumerator()) {
    $serviceName = $service.Key
    $port = $service.Value
    $url = "http://localhost:$port/api/v1/health"
    
    Write-Host "`nTesting $serviceName on port $port..." -ForegroundColor Yellow
    
    try {
        $response = Invoke-WebRequest -Uri $url -Method GET -TimeoutSec 10
        
        if ($response.StatusCode -eq 200) {
            $content = $response.Content | ConvertFrom-Json
            
            if ($content.message -eq "Service is healthy" -and $content.service) {
                Write-Host "✓ $serviceName: HEALTHY" -ForegroundColor Green
                Write-Host "  Response: $($response.Content)" -ForegroundColor Gray
            } else {
                Write-Host "✗ $serviceName: UNHEALTHY - Unexpected response format" -ForegroundColor Red
                Write-Host "  Response: $($response.Content)" -ForegroundColor Gray
                $allHealthy = $false
            }
        } else {
            Write-Host "✗ $serviceName: UNHEALTHY - Status Code: $($response.StatusCode)" -ForegroundColor Red
            $allHealthy = $false
        }
    }
    catch {
        Write-Host "✗ $serviceName: UNREACHABLE - $($_.Exception.Message)" -ForegroundColor Red
        $allHealthy = $false
    }
}

Write-Host "`n=========================" -ForegroundColor Green
if ($allHealthy) {
    Write-Host "All health endpoints are working correctly! ✓" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some health endpoints are not working properly! ✗" -ForegroundColor Red
    exit 1
}
