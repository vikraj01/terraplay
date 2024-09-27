variable "name" {
  description = "The name of the security group"
  type        = string
}

variable "description" {
  description = "A description for the security group"
  type        = string
  default     = "Managed by Terraform"
}

variable "vpc_id" {
  description = "The VPC ID to create the security group in"
  type        = string
}

variable "ingress_rules" {
  description = <<EOT
List of ingress rules. Each ingress rule must be an object with the following attributes:
- from_port: (number) The starting port number for the rule
- to_port: (number) The ending port number for the rule
- protocol: (string) The protocol (e.g., "tcp", "udp", "-1" for all)
- cidr_blocks: (list of strings) List of CIDR blocks
- description: (string) Description of the rule
EOT
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
  description = <<EOT
List of egress rules. Each egress rule must be an object with the following attributes:
- from_port: (number) The starting port number for the rule
- to_port: (number) The ending port number for the rule
- protocol: (string) The protocol (e.g., "tcp", "udp", "-1" for all)
- cidr_blocks: (list of strings) List of CIDR blocks
- description: (string) Description of the rule
EOT
  type = map(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
    description = string
  }))
  default = {}
}

variable "sg_tags" {
  description = "A map of tags to assign to the security group"
  type        = map(string)
  default     = {}
}