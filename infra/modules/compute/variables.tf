variable "instance_type" {
  description = "The type of instance to create."
  type        = string
  default     = "t3.micro"
  validation {
    condition     = contains(["t3.micro", "t2.micro", "t3.small", "t3.medium", "m5.large", "m5.xlarge"], var.instance_type)
    error_message = "The instance_type must be one of 't3.micro', 't2.micro', 't3.small', 't3.medium', 'm5.large', or 'm5.xlarge'."
  }
}

variable "ami" {
  description = "The AMI to use for the instance."
  type        = string
}

variable "subnet_id" {
  description = "The ID of the subnet to launch the instance in."
  type        = string
}

variable "key_name" {
  description = "The key pair to use for the instance."
  type        = string
  default     = null
}

variable "associate_public_ip_address" {
  description = "Associate a public IP address with the instance."
  type        = bool
  default     = true
}

variable "security_groups" {
  description = "List of security group IDs to assign to the instance."
  type        = list(string)
  default     = []
}

variable "monitoring_enabled" {
  description = "Enable detailed monitoring for the instance."
  type        = bool
  default     = false
}

variable "user_data" {
  description = "User data to provide when launching the instance."
  type        = string
  default     = ""
}

variable "root_volume_size" {
  description = "The size of the root EBS volume."
  type        = number
  default     = 8
  validation {
    condition     = var.root_volume_size >= 8
    error_message = "The root volume size must be at least 8 GB."
  }
}

variable "root_volume_type" {
  description = "The type of the root EBS volume."
  type        = string
  default     = "gp2"
  validation {
    condition     = contains(["gp2", "gp3", "io1", "st1", "sc1", "standard"], var.root_volume_type)
    error_message = "The volume type must be one of 'gp2', 'gp3', 'io1', 'st1', 'sc1', or 'standard'."
  }
}

variable "ebs_volumes" {
  description = "List of additional EBS volumes to attach."
  type = list(object({
    device_name = string
    volume_size = number
    volume_type = string
  }))
  default = []
  validation {
    condition = alltrue([
      for vol in var.ebs_volumes : vol.volume_size >= 1
    ])
    error_message = "All EBS volumes must have a size of at least 1 GB."
  }
}

variable "http_tokens" {
  description = "Whether or not to require IMDSv2 tokens."
  type        = string
  default     = "optional"
  validation {
    condition     = contains(["required", "optional"], var.http_tokens)
    error_message = "http_tokens must be either 'required' or 'optional'."
  }
}

variable "http_put_response_hop_limit" {
  description = "The allowed number of hops for PUT requests to the instance metadata service."
  type        = number
  default     = 1
  validation {
    condition     = var.http_put_response_hop_limit >= 1 && var.http_put_response_hop_limit <= 64
    error_message = "The http_put_response_hop_limit must be between 1 and 64."
  }
}

variable "http_endpoint" {
  description = "Enable or disable the HTTP metadata endpoint on the instance."
  type        = string
  default     = "enabled"
  validation {
    condition     = contains(["enabled", "disabled"], var.http_endpoint)
    error_message = "The http_endpoint must be either 'enabled' or 'disabled'."
  }
}

variable "cpu_credits" {
  description = "The credit option for CPU usage (standard/unlimited). Only applies to burstable instances."
  type        = string
  default     = "standard"
  validation {
    condition     = contains(["standard", "unlimited"], var.cpu_credits)
    error_message = "The cpu_credits option must be either 'standard' or 'unlimited'."
  }
}

variable "assign_eip" {
  description = "Assign an Elastic IP to the instance."
  type        = bool
  default     = false
}

variable "instance_name" {
  description = "The name to assign to the EC2 instance."
  type        = string
}

variable "ec2_tags" {
  description = "Additional tags to assign to the instance."
  type        = map(string)
  default = {}
}
