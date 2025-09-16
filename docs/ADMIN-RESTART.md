# Admin: Safe Kubernetes restarts from Admin Panel

This document explains how to enable secure rollout restarts of Mountain Service components from the Admin panel.

## What’s included
- Backend endpoint in employee-service:
  - POST /api/v1/admin/k8s/restart with JSON body {"deployment": "<name>"}
  - Admin‑only (JWT admin middleware)
  - In‑cluster call to Kubernetes API to emulate `kubectl rollout restart`
  - Strict allowlist of deployment names
- Helmized RBAC + ServiceAccount for employee-service (no raw k8s manifests)
- Admin UI restart section to trigger restarts safely

## Enable via Helm
In the employee-service chart, the ServiceAccount and RBAC are controlled by values and enabled by default. If you manage overrides, ensure these are set:

values.yaml (or your override file):

serviceAccount:
  create: true
  name: ""
rbac:
  create: true
  allowedDeployments:
    - employee-service
    - urgency-service
    - activity-service
    - version-service
    - docs-aggregator
    - docs-ui

Then upgrade/install:

helm upgrade --install employee-service api/charts/employee-service \
  -n mountain-service \
  -f your-values.yaml

3) Test the endpoint (requires admin JWT)
   # Acquire admin token via /api/v1/login (username=admin, ADMIN_PASSWORD)
   export TOKEN=...  # paste Bearer token
   curl -i -X POST \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     https://mountain-service.duckdns.org/api/v1/admin/k8s/restart \
     -d '{"deployment":"employee-service"}'

Expected 200 OK with JSON { message: "restart triggered", deployment, namespace, at }

## Security notes
- Endpoint is behind AdminMiddleware and requires a valid admin JWT
- Employee pod uses a dedicated ServiceAccount created by the Helm chart with a minimal Role granting patch on deployments in the release namespace
- Allowlist blocks restart of arbitrary objects; extend the list in code when you add new services


## RBAC: what it is and why we need it
RBAC (Role-Based Access Control) is Kubernetes’ permission model. Instead of giving the Pod full cluster credentials, we attach a minimal-permissions ServiceAccount to employee-service and bind it to a Role that can only patch Deployments in the namespace. This follows the principle of least privilege and avoids any external IAM or GitHub secrets.

What Helm creates in the employee-service chart when enabled:
- ServiceAccount: identity used by the Pod to talk to the kube‑apiserver (in-cluster token/CA mounted at /var/run/secrets/kubernetes.io/serviceaccount)
- Role (namespace-scoped): allows only the “patch” verb on “deployments”
- RoleBinding: binds the Role to the ServiceAccount

Why we need it:
- The restart endpoint performs a strategic‑merge PATCH on a Deployment to set kubectl.kubernetes.io/restartedAt, which triggers a rollout restart. Without RBAC, the request would be unauthorized.
- Scoping to the namespace and limiting to PATCH on Deployments reduces blast radius and keeps the system compliant and auditable.

How to verify:
- kubectl -n <ns> get sa,role,rolebinding | grep employee-service
- kubectl -n <ns> auth can-i patch deployment --as=system:serviceaccount:<ns>:<sa-name>

Disable/override if needed:
- Set serviceAccount.create=false and/or rbac.create=false in values if you provide your own SA/RBAC. Ensure that the Deployment’s serviceAccountName matches the SA you manage.

## Add another service to the allowlist
- File: api/employee/internal/handler/handler.go
- Locate allowed := map[string]bool{ ... }
- Add: "your-deployment": true
- Upgrade employee-service Helm release

## Troubleshooting
- 500 k8s auth not available: ensure the Helm values enable serviceAccount.create and rbac.create
- 500 k8s ca not available: check default serviceaccount token/CA volume mounts are present in the pod
- 500 k8s api error: Inspect logs in employee-service to see the exact status; confirm the Role/RoleBinding exist in the namespace and the deployment name is correct


