output "dynamodb_table_name" {
  description = "The name of the DynamoDB table."
  value       = aws_dynamodb_table.dynamodb_table.name
}

output "dynamodb_table_arn" {
  description = "The ARN of the DynamoDB table."
  value       = aws_dynamodb_table.dynamodb_table.arn
}

output "dynamodb_table_read_capacity" {
  description = "The read capacity units of the DynamoDB table."
  value       = aws_dynamodb_table.dynamodb_table.read_capacity
}

output "dynamodb_table_write_capacity" {
  description = "The write capacity units of the DynamoDB table."
  value       = aws_dynamodb_table.dynamodb_table.write_capacity
}

output "dynamodb_table_billing_mode" {
  description = "The billing mode for the DynamoDB table."
  value       = aws_dynamodb_table.dynamodb_table.billing_mode
}
