variable "region" {
  description = "The AWS region in which to deploy all resources"
  type        = string
}

variable "table_name" {
  description = "The name of the DynamoDB table"
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
  description = "The name of the project or application this resource is associated with, used to group resources and facilitate resource management and cost allocation."
  type        = string
}

variable "managed_by" {
  description = "The team or individual responsible for managing this resource, ensuring accountability and maintenance oversight."
  type        = string
  default     = "Terraform"
}