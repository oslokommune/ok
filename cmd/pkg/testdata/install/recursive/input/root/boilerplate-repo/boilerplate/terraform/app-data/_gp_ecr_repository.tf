{{ if dig "Ecr" "Enable" false . -}}
{{ template "doNotEdit" . }}
resource "aws_ecr_repository" "app" {

  name                 = "${local.environment}-${local.main_container_name}"
  image_tag_mutability = "IMMUTABLE"

  encryption_configuration {
    encryption_type = "KMS"
  }

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = local.common_tags

}

resource "aws_ecr_lifecycle_policy" "app" {

  repository = aws_ecr_repository.app.name

  policy = jsonencode(
    {
      rules = [
        {
          rulePriority = 1
          description  = "Keep ${local.ecr_max_image_count} images, expire all others."
          selection = {
            tagStatus   = "any"
            countType   = "imageCountMoreThan"
            countNumber = local.ecr_max_image_count
          }
          action = {
            type = "expire"
          }
        }
      ]
    }
  )

}
{{- else }}
# x-boilerplate-delete
{{- end }}
