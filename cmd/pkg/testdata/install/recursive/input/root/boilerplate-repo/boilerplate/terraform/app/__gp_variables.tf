{{ template "doNotEdit" . }}
{{- $mainContainerImage := not .ExampleImage.Enable }}
{{- $assumableCdRole := and (hasKey .IamForCicd "AssumableCdRole") .IamForCicd.AssumableCdRole }}
{{- $keepFile := or $mainContainerImage $assumableCdRole }}

{{- if $mainContainerImage }}

################################################################################
# Main container
################################################################################

variable "main_container_image_digest" {
  description = "The name of the main container"
  type        = string
}

variable "main_container_image_tag" {
  description = "The tag of the main container"
  type        = string
}

{{- end }}

{{- if $assumableCdRole }}

variable "assume_cd_role" {
  description = "Whether to assume the role used by continuous deployment in GitHub Actions"
  type        = string
  default     = "false"
}

{{- end -}}

{{- if not $keepFile }}
# x-boilerplate-delete
{{- end }}
