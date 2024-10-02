variable "role_name" {
  description = "Name of the IAM Role"
  type        = string
}

variable "role_description" {
  description = "Description for the IAM role"
  type        = string
  default     = "IAM role created by Terraform"
}

variable "role_path" {
  description = "The path for the IAM role"
  type        = string
  default     = "/"
}

variable "trusted_entities" {
  description = "List of trusted entities (AWS principals) for the assume role policy"
  type        = list(string)
  default     = []
}

variable "managed_policy_arns" {
  description = "List of manged policy ARNs to attach to the role"
  type        = list(string)
  default     = []
}

variable "inline_policies" {
  description = "Map of inline policy names to their JSON policies"
  type        = map(string)
  default = {
  }
}

variable "tags" {
  description = "A map of tags to assign to the IAM role"
  type        = map(string)
  default     = {}
}
