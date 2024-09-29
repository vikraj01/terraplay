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
