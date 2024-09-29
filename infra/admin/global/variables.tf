variable "project_name" {
  type = string
}

variable "managed_by" {
  type = string
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

variable "table_name" {
  description = "The name of the DynamoDB table."
  type        = string
}

variable "hash_key" {
  type        = string
  description = "The partition key for the DynamoDB table."
}

variable "range_key" {
  description = "The optional sort key for the DynamoDB table."
  type        = string
}

variable "ttl_attribute" {
  description = "The attribute for Time to Live (TTL) expiration."
  type        = string
}

variable "global_secondary_indexes" {
  description = "A list of global secondary index configurations."
  type = list(object({
    name               = string
    hash_key           = string
    hash_key_type      = string
    range_key          = string
    range_key_type     = string
    read_capacity      = number
    write_capacity     = number
    projection_type    = string
    non_key_attributes = list(string)
  }))
}

variable "key_pair_name" {
  
}