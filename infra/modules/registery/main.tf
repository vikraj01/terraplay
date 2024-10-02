resource "aws_ecr_repository" "ecr_repo" {
  name                 = var.ecr_repository_name
  image_tag_mutability = var.image_tag_mutability

  image_scanning_configuration {
    scan_on_push = var.enable_image_scanning
  }

  dynamic "encryption_configuration" {
    for_each = var.encryption_enabled ? [1] : []
    content {
      encryption_type = var.encryption_type
      kms_key         = var.kms_key_arn
    }
  }



  tags = var.tags
}

resource "aws_ecr_lifecycle_policy" "ecr_lifecycle_policy" {
  count      = var.lifecycle_policy_enabled ? 1 : 0
  repository = aws_ecr_repository.ecr_repo.name
  policy     = var.lifecycle_policy
}

resource "aws_ecr_repository_policy" "ecr_repo_policy" {
  count = var.repository_policy_text != "" ? 1 : 0

  repository = aws_ecr_repository.ecr_repo.name

  policy = var.repository_policy_text != "" ? var.repository_policy_text : null
}

resource "aws_ecr_repository_policy" "role_access_policy" {
  count = var.iam_role_arn != "" ? 1 : 0

  repository = aws_ecr_repository.ecr_repo.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          AWS = var.iam_role_arn
        },
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ],
        Resource = aws_ecr_repository.ecr_repo.arn
      }
    ]
  })
}
