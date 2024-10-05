variable "region" {
  description = "The AWS region in which to deploy all resources"
  type        = string
}

variable "table_name" {
  description = "The name of the DynamoDB table for storing Terraform state locks"
  type        = string
}

variable "hash_key" {
  description = "The partition key for the DynamoDB table"
  type        = string
}

variable "hash_key_type" {
  description = "The data type for the hash key (S: string, N: number, B: binary)"
  type        = string
  default     = "S"
}

variable "range_key" {
  description = "The optional sort key for the DynamoDB table"
  type        = string
  default     = null
}

variable "range_key_type" {
  description = "The data type for the range key (S: string, N: number, B: binary)"
  type        = string
  default     = "S"
}

variable "billing_mode" {
  description = "The billing mode for the DynamoDB table (PROVISIONED or PAY_PER_REQUEST)"
  type        = string
  default     = "PROVISIONED"
}

variable "project_name" {
  description = "The name of the project or application, used for resource grouping and identification"
  type        = string
}

variable "managed_by" {
  description = "The team or individual responsible for managing these resources"
  type        = string
  default     = "Terraform"
}

variable "instance_type" {
  description = "The type of EC2 instance to be used in the project (e.g., t2.micro, m5.large)"
  type        = string
  default     = "t2.micro"
}

variable "account_id" {
  description = "The AWS account ID where resources will be deployed"
  type        = string
}

variable "github_repo" {
  description = "The GitHub repository (in the format 'owner/repo') that will use OIDC to assume the AWS IAM role"
  type        = string
}
