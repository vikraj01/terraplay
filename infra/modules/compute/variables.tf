variable "ami" {
  description = "AMI ID to be used for the instance"
  type        = string
}

variable "instance_type" {
  description = "The type of the instance (e.g., t2.micro)"
  type        = string
  default     = "t2.micro"
}

variable "subnet_id" {
  description = "The subnet ID where the EC2 instance will be deployed"
  type        = string
}

variable "key_name" {
  description = "The key pair to use for the instance"
  type        = string
  default     = null
}

variable "security_groups" {
  description = "List of security group IDs to attach to the instance"
  type        = list(string)
  default     = []
}

variable "associate_public_ip_address" {
  description = "Whether to associate a public IP with the instance"
  type        = bool
  default     = true
}


variable "volume_size" {
  description = "The size of the EBS volume in GB"
  type        = number
  default     = 8
}

variable "volume_type" {
  description = "The type of the EBS volume (gp2, io1, etc.)"
  type        = string
  default     = "gp2"
}

variable "ec2_tags" {
  description = "A map of tags to assign to the resource"
  type        = map(string)
  default     = {}
}

variable "instance_name" {
  description = "The Name tag to assign to the instance"
  type        = string
}

variable "assign_eip" {
  description = "Whether to assign an Elastic IP to the instance"
  type        = bool
  default     = false
}

variable "user_data" {
  type = string
}