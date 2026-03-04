{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "desec-webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "desec-webhook.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "desec-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Selector labels must remain stable to avoid immutable selector changes.
*/}}
{{- define "desec-webhook.matchLabels" -}}
app: {{ include "desec-webhook.name" . }}
release: {{ .Release.Name }}
{{- end -}}

{{/*
Standard and legacy labels used across resources.
*/}}
{{- define "desec-webhook.labels" -}}
{{ include "desec-webhook.matchLabels" . }}
chart: {{ include "desec-webhook.chart" . }}
heritage: {{ .Release.Service }}
helm.sh/chart: {{ include "desec-webhook.chart" . }}
app.kubernetes.io/name: {{ include "desec-webhook.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "desec-webhook.selfSignedIssuer" -}}
{{- printf "%s-selfsign" (include "desec-webhook.fullname" .) -}}
{{- end -}}

{{- define "desec-webhook.rootCAIssuer" -}}
{{- printf "%s-ca" (include "desec-webhook.fullname" .) -}}
{{- end -}}

{{- define "desec-webhook.rootCACertificate" -}}
{{- printf "%s-ca" (include "desec-webhook.fullname" .) -}}
{{- end -}}

{{- define "desec-webhook.servingCertificate" -}}
{{- printf "%s-webhook-tls" (include "desec-webhook.fullname" .) -}}
{{- end -}}
