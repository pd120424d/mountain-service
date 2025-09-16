{{- define "employee-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "employee-service.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s" $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "employee-service.labels" -}}
app.kubernetes.io/name: {{ include "employee-service.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}



{{- define "employee-service.serviceAccountName" -}}
{{- $sa := .Values.serviceAccount | default dict -}}
{{- if ($sa.create | default false) -}}
{{- default (printf "%s-admin-ops" (include "employee-service.fullname" .)) $sa.name -}}
{{- else -}}
{{- default "default" $sa.name -}}
{{- end -}}
{{- end -}}
