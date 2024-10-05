# -------------------------
# Backend Storage Module for S3
# -------------------------
module "backend_storage" {
  source      = "../../modules/storage"
  bucket_name = var.bucket_name

  tags = merge(local.storage_tags, {
    Name      = var.bucket_name
    ManagedBy = var.managed_by
  })
}

# -------------------------
# DynamoDB Table for Terraform State Lock
# -------------------------
module "terraform_state_lock" {
  source         = "../../modules/database/dyanmodb"
  table_name     = var.table_name
  hash_key       = var.hash_key
  billing_mode   = var.billing_mode
  range_key      = var.range_key
  hash_key_type  = var.hash_key_type
  range_key_type = var.range_key_type

  tags = merge(local.dynamodb_tags, {
    Name      = var.table_name,
    ManagedBy = var.managed_by
  })
}

# -------------------------
# IAM OpenID Connect Provider for GitHub Actions
# -------------------------
resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

# -------------------------
# IAM Role for GitHub Actions
# -------------------------
resource "aws_iam_role" "github_actions_role" {
  name = "GithubActionsRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" : "repo:${var.github_repo}:*"
          }
        }
      }
    ]
  })
}

# -------------------------
# IAM Role Policy for GitHub Actions Role
# -------------------------
resource "aws_iam_role_policy" "github_actions_policy" {
  role = aws_iam_role.github_actions_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        Resource = [
          "arn:aws:s3:::${var.bucket_name}",
          "arn:aws:s3:::${var.bucket_name}/*"
        ]
      },
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          "dynamodb:UpdateItem"
        ],
        Resource = "arn:aws:dynamodb:${var.region}:${var.account_id}:table/${var.table_name}"
      },
      {
        Effect = "Allow",
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ],
        Resource = "arn:aws:ecr:${var.region}:${var.account_id}:repository/*"
      }
    ]
  })
}
