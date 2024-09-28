variable "vpc_config" {
  type = object({
    name       = string
    cidr_block = string
  })

  validation {
    condition     = can(cidrnetmask(var.vpc_config.cidr_block))
    error_message = "The vpc_config option must contain a valid CIDR block"
  }
}

variable "subnet_config" {
  type = map(object({
    cidr_block = string
    public     = optional(bool, false)
    az         = string
  }))

  validation {
    condition = alltrue([
      for config in var.subnet_config : can(cidrnetmask(config.cidr_block))
    ])
    error_message = "The subnet config option must contain a valid CIDR block"
  }
}

variable "tags" {
  description = "A map of tags to apply to all resources created in this module"
  type        = map(string)
  default     = {}
}
