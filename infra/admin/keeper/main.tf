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
          "s3:DeleteObject",
          "s3:CreateBucket",       # Allow bucket creation
          "s3:DeleteBucket",       # Allow bucket deletion
          "s3:PutBucketPolicy",    # Allow modifying bucket policy
          "s3:PutBucketAcl"        # Allow setting bucket permissions
        ],
        Resource = [
          "arn:aws:s3:::${var.bucket_name}",
          "arn:aws:s3:::${var.bucket_name}/*"
        ]
      },
      # DynamoDB Permissions
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          "dynamodb:UpdateItem",
          "dynamodb:DescribeTable",  # Describe the table
          "dynamodb:CreateTable",    # Create tables
          "dynamodb:DeleteTable",    # Delete tables
          "dynamodb:UpdateTable"     # Modify table configurations
        ],
        Resource = "arn:aws:dynamodb:${var.region}:${var.account_id}:table/*"
      },
      # ECR Permissions
      {
        Effect = "Allow",
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:CreateRepository",   # Allow repository creation
          "ecr:DeleteRepository",   # Allow repository deletion
          "ecr:SetRepositoryPolicy" # Set repository policy
        ],
        Resource = "arn:aws:ecr:${var.region}:${var.account_id}:repository/*"
      },
      # EC2 Permissions
      {
        Effect = "Allow",
        Action = [
          "ec2:RunInstances",        # Launch EC2 instances
          "ec2:TerminateInstances",  # Terminate instances
          "ec2:DescribeInstances",   # Describe running instances
          "ec2:StartInstances",      # Start instances
          "ec2:StopInstances",       # Stop instances
          "ec2:RebootInstances",     # Reboot instances
          "ec2:DescribeImages",      # Describe AMIs (for selecting the right AMI)
          "ec2:DescribeSubnets",     # Get subnet details
          "ec2:DescribeSecurityGroups", # Describe security groups
          "ec2:DescribeAvailabilityZones", # Get availability zones
          "ec2:CreateSecurityGroup", # Create security groups
          "ec2:DeleteSecurityGroup", # Delete security groups
          "ec2:CreateTags",          # Tag resources
          "ec2:DeleteTags"           # Delete tags
        ],
        Resource = "*"
      },
      # IAM Permissions (for managing roles and policies)
      {
        Effect = "Allow",
        Action = [
          "iam:GetRole",             # Get information about roles
          "iam:PassRole",            # Pass role to other services (e.g., EC2)
          "iam:CreateRole",          # Create IAM roles
          "iam:DeleteRole",          # Delete IAM roles
          "iam:UpdateRole",          # Update roles
          "iam:AttachRolePolicy",    # Attach policies to roles
          "iam:DetachRolePolicy",    # Detach policies from roles
          "iam:CreatePolicy",        # Create new policies
          "iam:DeletePolicy"         # Delete policies
        ],
        Resource = "arn:aws:iam::${var.account_id}:role/*"
      },
      # Secrets Manager Permissions
      {
        Effect = "Allow",
        Action = [
          "secretsmanager:GetSecretValue",  # Get the secret's value
          "secretsmanager:DescribeSecret",  # Describe the secret
          "secretsmanager:CreateSecret",    # Create new secrets
          "secretsmanager:UpdateSecret",    # Update secrets
          "secretsmanager:DeleteSecret"     # Delete secrets
        ],
        Resource = "arn:aws:secretsmanager:${var.region}:${var.account_id}:secret/*"
      },
      {
        Effect = "Allow",
        Action = "secretsmanager:ListSecrets", # List all secrets
        Resource = "*"
      }
    ]
  })
}

