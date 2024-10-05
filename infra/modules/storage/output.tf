output "bucket_id" {
  description = "The ID of the created S3 bucket."
  value       = aws_s3_bucket.this.id
}

output "bucket_arn" {
  description = "The ARN of the created S3 bucket."
  value       = aws_s3_bucket.this.arn
}

output "bucket_name" {
  description = "Name of the bucket"
  value       = aws_s3_bucket.this.bucket
}
