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
  source = "../../modules/database/dyanmodb"

  table_name = var.table_name

  hash_key = var.hash_key

  range_key = var.range_key

  ttl_attribute = var.ttl_attribute

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
  source = "../../modules/keys"

  private_key_path = "${path.module}/sensitive/my_private_key.pem"

  key_pair_name = var.key_pair_name
}




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

# module "ec2_role_with_ecr_access" {
#   source              = "../../modules/roles"
#   role_name           = "ec2-role-with-ecr-access"
#   trusted_entities    = var.trusted_entities
#   managed_policy_arns = [data.aws_iam_policy.ecr_full_access.arn]
# }

module "global_ecr_repository" {
  source              = "../../modules/registery"
  ecr_repository_name = var.ecr_repository_name
  # iam_role_arn        = module.ec2_role_with_ecr_access.role_arn
  # depends_on = [ module.ec2_role_with_ecr_access ]
}



resource "aws_s3_bucket" "global_bucket" {
  bucket = "global-bucket-893606"

}

output "global_bucket_name" {
  value = aws_s3_bucket.global_bucket.bucket
}

output "global_bucket_id" {
  value = aws_s3_bucket.global_bucket.id
}
