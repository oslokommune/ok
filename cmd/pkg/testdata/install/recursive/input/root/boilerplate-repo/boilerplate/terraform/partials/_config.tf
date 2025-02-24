{{- define "commonTags" -}}
  common_tags = {
    "Team"        = local.team_name
    "Environment" = local.environment
    "CreatedBy"   = "ok-golden-path"
  }
{{- end -}}

{{- define "commonLocals" -}}
  account_id  = "{{ .AccountId }}"
  region      = "{{ .Region }}"
  team_name   = "{{ .Team }}"
  environment = "{{ .Environment }}"
{{- end -}}
