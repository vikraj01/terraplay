resource "aws_dynamodb_table" "dynamodb_table" {
  name         = var.table_name
  billing_mode = var.billing_mode
  hash_key     = var.hash_key

  attribute {
    name = var.hash_key
    type = var.hash_key_type
  }

  dynamic "attribute" {
    for_each = var.range_key != null ? [var.range_key] : []
    content {
      name = var.range_key
      type = var.range_key_type
    }
  }

  read_capacity  = var.read_capacity
  write_capacity = var.write_capacity

  dynamic "ttl" {
    for_each = var.ttl_attribute != null ? [var.ttl_attribute] : []
    content {
      attribute_name = var.ttl_attribute
      enabled        = true
    }
  }

  dynamic "global_secondary_index" {
    for_each = var.global_secondary_indexes
    content {
      name      = global_secondary_index.value.name
      hash_key  = global_secondary_index.value.hash_key
      range_key = global_secondary_index.value.range_key

      projection_type    = global_secondary_index.value.projection_type
      non_key_attributes = global_secondary_index.value.projection_type == "INCLUDE" ? global_secondary_index.value.non_key_attributes : []

      read_capacity  = global_secondary_index.value.read_capacity
      write_capacity = global_secondary_index.value.write_capacity
    }
  }

  tags = merge({
    Name = var.table_name
  }, var.tags)
}
