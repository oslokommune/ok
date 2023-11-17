# This file contains Terraform variables of type local
locals {
  # Shared variable used by templates and modules:
  common_tags = {
    "Team"        = var.team_name
    "Environment" = var.environment
    "CreatedBy"   = "ok-golden-path"
  }
  
  hei = "hallo"

  # Add your own configuration here:
}
