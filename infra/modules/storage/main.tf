resource "aws_s3_bucket" "this" {
  bucket = var.bucket_name

  dynamic "server_side_encryption_configuration" {
    for_each = var.sse_algorithm != "" ? [var.sse_algorithm] : []
    content {
      rule {
        apply_server_side_encryption_by_default {
          sse_algorithm = server_side_encryption_configuration.value
        }
      }
    }
  }

  tags = var.tags
}

resource "aws_s3_bucket_policy" "name" {
  count  = var.bucket_policy != "" ? 1 : 0
  bucket = aws_s3_bucket.this.id
  policy = var.bucket_policy
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.bucket
  versioning_configuration {
    status = var.versioning
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  bucket = aws_s3_bucket.this.id
  count  = length(var.lifecycle_rules) > 0 ? 1 : 0
  dynamic "rule" {
    for_each = var.lifecycle_rules
    content {
      id     = rule.value.id
      status = rule.value.enabled ? "Enabled" : "Disabled"

      filter {
        prefix = rule.value.prefix
      }

      transition {
        days          = rule.value.transition_days
        storage_class = rule.value.storage_class
      }

      expiration {
        days = rule.value.expiration_days
      }

      noncurrent_version_expiration {
        noncurrent_days = rule.value.noncurrent_version_expiration_days
      }

      abort_incomplete_multipart_upload {
        days_after_initiation = rule.value.abort_incomplete_multipart_upload_days
      }
    }
  }
}
