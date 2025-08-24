#!/bin/bash

# Setup script for Kubernetes Ingress with NGINX Ingress Controller
# This script installs the NGINX Ingress Controller and configure it for Duck DNS domain

set -e

echo "Setting up Kubernetes Ingress for Mountain Service"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed or not in PATH"
    exit 1
fi

# Check if we can connect to the cluster
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

print_status "Connected to Kubernetes cluster successfully"

# Install NGINX Ingress Controller
print_status "Installing NGINX Ingress Controller..."

# Check if ingress-nginx namespace already exists
if kubectl get namespace ingress-nginx &> /dev/null; then
    print_warning "ingress-nginx namespace already exists"
else
    print_status "Creating ingress-nginx namespace..."
    kubectl create namespace ingress-nginx
fi

# Install NGINX Ingress Controller using the official manifest
print_status "Applying NGINX Ingress Controller manifest..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

# Wait for the ingress controller to be ready
print_status "Waiting for NGINX Ingress Controller to be ready..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=300s

print_success "NGINX Ingress Controller is ready!"

# Get the external IP of the ingress controller
print_status "Getting external IP of the ingress controller..."
EXTERNAL_IP=""
while [ -z $EXTERNAL_IP ]; do
    print_status "Waiting for external IP..."
    EXTERNAL_IP=$(kubectl get svc ingress-nginx-controller -n ingress-nginx --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
    [ -z "$EXTERNAL_IP" ] && sleep 10
done

print_success "External IP: $EXTERNAL_IP"

# Apply the mountain-service namespace if it doesn't exist
print_status "Ensuring mountain-service namespace exists..."
kubectl apply -f namespaces.yaml

# Apply the updated frontend service (ClusterIP instead of NodePort)
print_status "Applying updated frontend service..."
kubectl apply -f frontend/frontend.yaml

# Apply the ingress configuration
print_status "Applying ingress configuration..."
kubectl apply -f ingress.yaml

print_success "Ingress setup completed
