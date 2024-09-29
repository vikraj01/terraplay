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



variable "create_networking" {
  description = "Boolean to determine whether to set up networking (VPC, subnets, gateway, etc.)"
  type        = bool
  default     = false
}

variable "create_firewall" {
  description = "Boolean to determine whether to set up networking (VPC, subnets, gateway, etc.)"
  type        = bool
  default     = false
}

variable "networking_tags" {
  description = "A map of tags to apply to all resources created in this module"
  type        = map(string)
  default     = {}
}



variable "create_key" {
  type = bool
  default = false
}

variable "key_pair_name" {
  type = string
  
}