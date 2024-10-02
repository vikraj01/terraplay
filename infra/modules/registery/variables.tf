variable "ecr_repository_name" {
  description = "Name of the ECR repository"
  type        = string
}

variable "enable_image_scanning" {
  description = "Enable automated image scanning"
  type        = bool
  default     = true
}

variable "image_tag_mutability" {
  description = "Specifies whether image tags are mutable"
  type        = string
  default     = "IMMUTABLE"
}

variable "lifecycle_policy_enabled" {
  description = "Enable lifecycle policy for the repository"
  type        = bool
  default     = false
}

variable "lifecycle_policy" {
  description = "JSON string for the ECR lifecycle policy (use file function to load from an external file)"
  type        = string
  default     = ""
}

variable "encryption_enabled" {
  description = "Enable encryption for the repository"
  type        = bool
  default     = true
}

variable "encryption_type" {
  description = "Encryption configuration for the ECR repository. Can be AES256 or KMS."
  type        = string
  default     = "AES256"
}

variable "kms_key_arn" {
  description = "Optional KMS key ARN if using KMS encryption"
  type        = string
  default     = ""
}

variable "iam_role_arn" {
  description = "Optional IAM Role ARN for accessing the ECR repository"
  type        = string
  default     = ""
}

variable "repository_policy_text" {
  description = "JSON string for the ECR repository policy to allow cross-account or external access"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags to apply to all resources created by this module"
  type        = map(string)
  default     = {}
}
