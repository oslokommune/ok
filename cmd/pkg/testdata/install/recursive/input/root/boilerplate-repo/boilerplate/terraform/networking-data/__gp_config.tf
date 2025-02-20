{{ template "doNotEdit" . }}
################################################################################
# Config
################################################################################
locals {
  # To destroy this stack:
  # 1. Set skip_destroy = false
  # 2. Run terraform apply
  # 3. Run terraform destroy
  skip_destroy = true

  ################################################################################
  # Common
  ################################################################################
  {{ template "commonTags" }}

  {{ template "commonLocals" . }}

}
