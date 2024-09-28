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

variable "vpc_cidr" {
  description = "The CIDR block for the VPC. This defines the IP address range for the VPC."
  type        = string
}
variable "subnet_config" {
  description = "A map of objects that define the configuration for each subnet, including CIDR block, availability zone (AZ), and an optional flag indicating if it's public."
  type = map(object({
    cidr_block = string
    public     = optional(bool, false)
    az         = string
  }))
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

variable "security_group_description" {
  description = "A description for the security group"
  type        = string
}

variable "ingress_rules" {
  type = map(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
    description = string
  }))
  default = {}
}

variable "egress_rules" {
  type = map(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
    description = string
  }))
  default = {}
}

variable "game" {
  description = "The name of the game"
  type        = string
  default     = ""
}
