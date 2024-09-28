variable "game" {
  type        = string
  description = "The name of the game server to set up."
  default     = ""
  validation {
    condition     = contains(["minecraft", "terraria", "valheim", "minetest", ""], lower(var.game))
    error_message = "The game must be one of: Minecraft, Terraria, Valheim."
  }
}

variable "instance_type" {
  description = "The type of EC2 instance."
  type        = string
  default     = "t3.micro"
}

variable "subnet_id" {
  description = "The subnet ID to launch the instance in."
  type        = string
}

variable "security_group_ids" {
  description = "List of security group IDs."
  type        = list(string)
}

variable "assign_eip" {
  description = "Whether to assign an Elastic IP to the instance."
  type        = bool
  default     = false
}

variable "ec2_tags" {
  description = "Additional tags for the EC2 instance."
  type        = map(string)
  default     = {}
}

variable "ebs_volumes" {
  description = "List of additional EBS volumes to attach."
  type = list(object({
    device_name = string
    volume_size = number
    volume_type = string
  }))
  default = []
}


variable "key_name" {
  type    = string
  default = null
}
