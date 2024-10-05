# -------------------------
# VPC Creation Module
# -------------------------
module "terraplay_vpc" {
  source = "../../modules/networking"

  vpc_config = {
    cidr_block = var.vpc_cidr
    name       = "${var.project_name}-vpc"
  }

  subnet_config = var.subnet_config
}

# -------------------------
# DynamoDB Table Module for Session Tracking
# -------------------------
module "session_table" {
  source                   = "../../modules/database/dyanmodb"
  table_name               = var.table_name
  hash_key                 = var.hash_key
  range_key                = var.range_key
  ttl_attribute            = var.ttl_attribute
  global_secondary_indexes = var.global_secondary_indexes
}

# -------------------------
# Field Definitions for the Sessions Table
# -------------------------
# Field Name      Data Type     Description
# -----------------------------------------------
# session_id      String (UUID) Unique ID for each session.
# user_id         String        Foreign key linking to the Users table.
# game_name       String        Name of the game (e.g., Minecraft, Minetest).
# status          String (ENUM) Session status (e.g., active, stopped, killed).
# start_time      Timestamp     When the session started.
# delete_time     Timestamp     When the session was deleted or destroyed.
# instance_id     String        The ID of the EC2 instance running the game.
# state_file      String (S3 URI) S3 path to the game state backup for the session.
# created_at      Timestamp     Timestamp when the session was created.
# updated_at      Timestamp     Timestamp when the session was last updated.

# -------------------------
# SSH Key Pair Module
# -------------------------
module "ssh_key" {
  source           = "../../modules/keys"
  private_key_path = "${path.module}/sensitive/my_private_key.pem"
  key_pair_name    = var.key_pair_name
}



# -------------------------
# ECR Repository for Bot Server
# -------------------------
module "global_ecr_repository" {
  source              = "../../modules/registery"
  ecr_repository_name = var.ecr_repository_name
}

# -------------------------
# Global S3 Bucket for Game Data Backup
# -------------------------
module "global_bucket" {
  source      = "../../modules/storage"
  bucket_name = "global-terraplay-bucketv1"
}

# -------------------------
# Secrets Manager For Environment Variables For The Bot
# -------------------------
resource "aws_secretsmanager_secret" "this" {
  name = "terraplay"
}



# -------------------------
# Server for Discord Bot (EC2 Instance)
# -------------------------
module "bot_server" {
  source          = "../../modules/compute"
  subnet_id       = module.terraplay_vpc.public_subnets["public"].subnet_id
  ami             = data.aws_ami.amazon_linux.id
  instance_name   = "nimbus_server-${terraform.workspace}"
  instance_type   = var.instance_type
  key_name        = module.ssh_key.aws_key_name
  security_groups = [module.bot_firewall.security_group_id]

  ebs_volumes = [
    {
      device_name = "/dev/sdf"
      volume_size = 50
      volume_type = "gp3"
    }
  ]

  user_data = file("${path.module}/scripts/user_data.sh")

  ec2_tags = {
    Name = "nimbus-${terraform.workspace}"
    Type = "Compute"
  }
}



# -------------------------
# IAM Role for Bot Server (EC2 Instance)
# -------------------------
resource "aws_iam_role" "bot_server_role" {
  name = "TerraplayBotServerRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# -------------------------
# IAM Policy for Bot Server - Access to DynamoDB and S3
# -------------------------
resource "aws_iam_policy" "bot_server_access_policy" {
  name = "bot-server-access-policy"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      # DynamoDb permissions
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem"
        ],
        Resource = "arn:aws:dynamodb:${var.region}:${var.account_id}:table/${module.session_table.dynamodb_table_name}"
      },
      # S3 Permissons
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        Resource = [
          "arn:aws:s3:::${module.global_bucket.bucket_name}",
          "arn:aws:s3:::${module.global_bucket.bucket_name}/*"
        ]
      },
      # Secrets Manager Permissions
      {
        Effect = "Allow",
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ],
        Resource = "arn:aws:secretsmanager:${var.region}:${var.account_id}:secret:${aws_secretsmanager_secret.this.name}"
      },
      {
        Effect   = "Allow",
        Action   = "secretsmanager:ListSecrets",
        Resource = "*"
      },
      {
        Effect   = "Allow",
        Action   = "ec2:*",
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "bot_server_role_policy_attachment" {
  role       = aws_iam_role.bot_server_role.name
  policy_arn = aws_iam_policy.bot_server_access_policy.arn
}

# -------------------------
# IAM Policy for EC2 Instance - EC2 Management Permissions / For Github Actions
# -------------------------
# resource "aws_iam_policy" "ec2_instance_policy" {
#   name = "ec2-instance-github-actions-policy"

#   policy = jsonencode({
#     Version = "2012-10-17",
#     Statement = [
#       {
#         Effect = "Allow",
#         Action = [
#           "ec2:StartInstances",
#           "ec2:StopInstances",
#           "ec2:RebootInstances"
#         ],
#         Resource = [
#           "arn:aws:ec2:${var.region}:${var.account_id}:instance/${module.bot_server.instance_id}"
#         ]
#       }
#     ]
#   })
# }

# resource "aws_iam_role_policy_attachment" "attach_custom_policy_to_role" {
#   role       = data.aws_iam_role.github_actions_role.name
#   policy_arn = aws_iam_policy.ec2_instance_policy.arn
# }




// Temporarily actions have all the power! but after finalization! following stricter policy tightening will be happen for all of them
