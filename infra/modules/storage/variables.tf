variable "bucket_name" {
  description = "The name of the S3 Bucket"
  type        = string
  validation {
    condition     = length(var.description) > 3 && length(var.description) < 64
    error_message = "value"
  }
}

variable "versioning" {
  description = "Enable versioning on the S3 bucket"
  type        = string
  default     = "Disabled"

  validation {
    condition     = var.versioning == "Disabled" || var.versioning == "Enabled"
    error_message = "versioning must be either 'Disabled' or 'Enabled'"
  }
}

variable "sse_algorithm" {
  description = "The server-side encryption algorithm to use, Must be AES256 or aws:kms"
  type        = string
  default     = ""
  validation {
    condition     = var.sse_algorithm == "AES256" || var.sse_algorithm == "aws:kms" || var.sse_algorithm == ""
    error_message = "sse_algorithm must be either 'AES256', 'aws:kms', or left empty for no encryption."
  }
}

variable "lifecycle_rules" {
  description = "A list of lifecycle rule configurations for the S3 bucket."
  type = list(object({
    id                                     = string
    enabled                                = bool
    prefix                                 = string
    transition_days                        = number
    storage_class                          = string
    expiration_days                        = number
    noncurrent_version_expiration_days     = number
    abort_incomplete_multipart_upload_days = number
  }))
  default = []
}

variable "bucket_policy" {
  description = "A JSON policy document for the S3 bucket."
  type        = string
  default     = ""
}

variable "tags" {
  description = "A map of tags to assign to the bucket."
  type        = map(string)
  default     = {}
}
