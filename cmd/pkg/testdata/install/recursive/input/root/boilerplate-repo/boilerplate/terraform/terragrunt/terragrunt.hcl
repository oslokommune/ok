{{ if dig "Terragrunt" "Enable" false . -}}
{{ template "doNotEdit" . }}
dependencies {
  paths = [
{{- if dig "Terragrunt" "DependenciesPaths" false . }}
  {{- $deps := .Terragrunt.DependenciesPaths }}
  {{- range $_, $element := $deps }}
    "{{ $element }}",
  {{- end }}
{{- else }}
  {{- $deps := .DefaultDependenciesPaths }}
  {{- range $_, $element := $deps }}
  {{- if ne $element "" }}
    "{{ $element }}",
  {{- end }}
  {{- end }}
{{- end }}
  ]
}

# Merge dependencies
include "custom" {
  path = "terragrunt_custom.hcl"
}

terraform {
  source = "."
}
{{- else }}
# x-boilerplate-delete
{{- end }}
