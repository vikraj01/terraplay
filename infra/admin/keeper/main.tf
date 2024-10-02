resource "random_id" "bucket_prefix" {
  byte_length = 8
}
module "backend_storage" {
  source      = "../../modules/storage"
  bucket_name = "${var.project_name}-backend-${random_id.bucket_prefix.hex}"

  tags = merge(local.storage_tags, {
    Name = "${var.project_name}-backend-${random_id.bucket_prefix.hex}"
  })
}

module "terraform_state_lock" {
  source         = "../../modules/database/dyanmodb"
  table_name     = var.table_name
  hash_key       = var.hash_key
  billing_mode   = var.billing_mode
  range_key      = var.range_key
  hash_key_type  = var.hash_key_type
  range_key_type = var.range_key_type

  tags = merge(local.dynamodb_tags, {
    Name = var.table_name
  })
}