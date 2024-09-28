output "dynamodb_table_name" {
  description = "The name of the DynamoDB table for Terraform state locking."
  value       = module.terraform_state_lock.dynamodb_table_name
}

output "dynamodb_table_arn" {
  description = "The ARN of the DynamoDB table for Terraform state locking."
  value       = module.terraform_state_lock.dynamodb_table_arn
}

output "dynamodb_table_read_capacity" {
  description = "The read capacity units of the DynamoDB table."
  value       = module.terraform_state_lock.dynamodb_table_read_capacity
}

output "dynamodb_table_write_capacity" {
  description = "The write capacity units of the DynamoDB table."
  value       = module.terraform_state_lock.dynamodb_table_write_capacity
}

output "dynamodb_table_billing_mode" {
  description = "The billing mode for the DynamoDB table."
  value       = module.terraform_state_lock.dynamodb_table_billing_mode
}

output "bucket_id" {
  description = "The ID of the created S3 bucket."
  value       = module.backend_storage.bucket_id
}

output "bucket_arn" {
  description = "The ARN of the created S3 bucket."
  value       = module.backend_storage.bucket_arn
}
