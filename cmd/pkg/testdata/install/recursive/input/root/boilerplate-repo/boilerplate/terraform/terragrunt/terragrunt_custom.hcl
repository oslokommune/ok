{{ if dig "Terragrunt" "Enable" false . -}}
dependencies {
  paths = []
}
{{- else }}
# x-boilerplate-delete
{{- end }}
