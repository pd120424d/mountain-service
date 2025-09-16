# HTTPS (TLS) with Helm for mountain-service

This guide enables HTTPS for https://mountain-service.duckdns.org using Traefik + cert-manager via Helm only. It configures TLS on the Helm chart ingresses (frontend, docs-ui, docs-aggregator) and removes the need for raw k8s manifests.

## What’s included in this repo
- Charts updated to support TLS on Ingress templates (frontend, docs-ui, docs-aggregator)
- Values defaulted to enable TLS (websecure) and Let's Encrypt via cert-manager
- Docs-aggregator configured to advertise `https` in generated OpenAPI specs to avoid mixed content
- Optional HSTS header added in the frontend nginx config
- New infra Helm chart to create cert-manager ClusterIssuers (and optional Traefik redirect middleware): `api/charts/infra`

## 1) Prerequisites
- Traefik ingress controller installed and exposing `web` (80) and `websecure` (443) entrypoints
- DNS `mountain-service.duckdns.org` points to the Traefik LoadBalancer IP
- cert-manager installed in the cluster

Install cert-manager itself (once):
```
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace cert-manager --create-namespace \
  --version v1.14.5
```

## 2) Install infra chart (ClusterIssuers)
This chart installs Let’s Encrypt staging and prod ClusterIssuers.
Set your contact email:
```
helm upgrade --install infra ./charts/infra \
  -n mountain-service --create-namespace \
  --set email=pd120424d@etf.bg.ac.rs
```
Optional: also create an HTTP→HTTPS redirect middleware in the same namespace:
```
helm upgrade --install infra ./charts/infra \
  -n mountain-service \
  --set email=pd120424d@etf.bg.ac.rs \
  --set createRedirectMiddleware=true
```
Note: To use the middleware, add this annotation to each chart’s `ingress.annotations`:
```
traefik.ingress.kubernetes.io/router.middlewares: mountain-service-redirect-https@kubernetescrd
```
(Or set a different name via `--set middlewareName=...` and reference it accordingly.)

## 3) Enable TLS on app charts (values already prepared)
The following charts are configured to expose HTTPS on the same host with TLS from cert-manager (letsencrypt-prod):
- frontend
- docs-ui
- docs-aggregator

They use:
- `ingress.annotations.traefik.ingress.kubernetes.io/router.entrypoints: web,websecure`
- `ingress.annotations.traefik.ingress.kubernetes.io/router.tls: "true"`
- `ingress.annotations.cert-manager.io/cluster-issuer: letsencrypt-prod`
- `ingress.tls.secretName: mountain-service-tls`

Deploy/upgrade them (adjust release names to your setup):
```
helm upgrade --install frontend ./charts/frontend -n mountain-service
helm upgrade --install docs-ui ./charts/docs-ui -n mountain-service
helm upgrade --install docs-aggregator ./charts/docs-aggregator -n mountain-service
```

Notes:
- Backends (employee, urgency, activity, version) do not need external ingress; they’re proxied by the frontend and remain ClusterIP only.
- We leave their ingresses disabled by default.

## 4) Verify issuance and endpoints
```
# Certificate/secret (namespace: mountain-service)
kubectl -n mountain-service get certificate
kubectl -n mountain-service describe certificate mountain-service-tls
kubectl -n mountain-service get secret mountain-service-tls

# Smoke checks
curl -I https://mountain-service.duckdns.org/health
curl -I https://mountain-service.duckdns.org/api/v1/health
curl -I https://mountain-service.duckdns.org/api/v1/docs/swagger-config.json

# Optional: if redirect middleware added
curl -I http://mountain-service.duckdns.org/
```
Browser:
- Visit `https://mountain-service.duckdns.org`
- Confirm padlock and no mixed-content warnings
- Open `/docs`; specs should load over HTTPS

## 5) Notes and hardening
- HSTS is added at the frontend nginx level
- CORS defaults allow both http and https origins; after HTTPS is stable, you can remove the http origin in `api/shared/server/cors.go` for stricter security
- OAuth providers: ensure redirect URIs use `https`

## 6) Rollback
- Roll back to previous releases with `helm rollback`
- To temporarily go HTTP-only, remove TLS annotations and `ingress.tls` from chart values (not recommended)

## 7) Troubleshooting
- Issuance stuck:
  - `kubectl -n cert-manager get pods`
  - `kubectl -n cert-manager logs deploy/cert-manager -f`
  - `kubectl -n mountain-service describe challenge,order`
- 404 on HTTP-01 challenge: ensure Traefik class is `traefik` and `web` entrypoint is available cluster-wide; avoid conflicting ingresses for same host/path
- Mixed content in Docs: `EXTERNAL_SCHEME` is set to `https` in `docs-aggregator` values; re-deploy and clear browser cache
