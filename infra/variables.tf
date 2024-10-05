
variable "project_name" {
  description = "The name of the project or application this resource is associated with, used to group resources and facilitate resource management and cost allocation."
  type        = string
  default     = "terraplay"
}
variable "managed_by" {
  description = "The team or individual responsible for managing this resource, ensuring accountability and maintenance oversight."
  type        = string
  default     = "Terraform"
}

variable "bucket" {
  description = "The bucket name to store the state"
  type = string
}

variable "key" {
  description = "The key that identify the state"
  type = string
}

variable "dynamodb_table"{
  description = "The key to lock the state"
  type = string
}

variable "region" {
  description = "The AWS region in which to deploy all resources"
  type        = string
  default     = "ap-south-1"
}
variable "global_workspace_name"{
  description = "This is the name of the global workspace"
  type = string
}
