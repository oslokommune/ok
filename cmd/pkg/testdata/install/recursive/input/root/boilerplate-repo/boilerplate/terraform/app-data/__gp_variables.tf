{{ template "doNotEdit" . }}
{{- $assumableCdRole := and (hasKey .IamForCicd "AssumableCdRole") .IamForCicd.AssumableCdRole }}

{{- if $assumableCdRole }}

variable "assume_cd_role" {
  description = "Whether to assume the role used by continuous deployment in GitHub Actions"
  type        = string
  default     = "false"
}

{{- else }}
# x-boilerplate-delete
{{- end }}
