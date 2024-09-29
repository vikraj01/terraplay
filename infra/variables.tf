variable "region" {
  description = "The AWS region in which to deploy all resources"
  type        = string
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