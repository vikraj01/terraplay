# -------------------------
# Project Information
# -------------------------
variable "project_name" {
  description = "The name of the project or application. Used as a prefix for naming resources to group them logically."
  type        = string
}

variable "managed_by" {
  description = "The team or individual responsible for managing the infrastructure. Typically used as a tag for resource identification."
  type        = string
}

# -------------------------
# VPC and Networking Configuration
# -------------------------
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

# -------------------------
# Security Group Configuration
# -------------------------
variable "security_group_description" {
  description = "A description for the security group"
  type        = string
}

variable "ingress_rules" {
  description = "A map defining ingress (inbound) rules for the security group. Includes the port range, protocol, CIDR blocks, and a description for each rule."
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
  description = "A map defining egress (outbound) rules for the security group. Includes the port range, protocol, CIDR blocks, and a description for each rule."
  type = map(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
    description = string
  }))
  default = {}
}

# -------------------------
# DynamoDB Table Configuration
# -------------------------
variable "table_name" {
  description = "The name of the DynamoDB table used for session tracking or other purposes."
  type        = string
}

variable "hash_key" {
  description = "The partition key for the DynamoDB table."
  type        = string
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

# -------------------------
# EC2 Instance Configuration
# -------------------------
variable "key_pair_name" {
  description = "The name of the SSH key pair used for connecting to the EC2 instances. Ensure the private key is stored securely."
  type        = string
}

variable "instance_type" {
  description = "The EC2 instance type for the server (e.g., t2.micro, m5.large), defining the compute power and pricing of the instance."
  type        = string
}

variable "ebs_volumes" {
  description = "A list of EBS volumes to attach to the EC2 instance. Each volume includes configuration for device name, size, and volume type (e.g., gp2, gp3)."
  type = list(object({
    device_name = string
    volume_size = number
    volume_type = string
  }))
  default = []
}

# -------------------------
# ECR Configuration
# -------------------------
variable "ecr_repository_name" {
  description = "The name of the Amazon Elastic Container Registry (ECR) repository where Docker images are stored for the bot server."
  type        = string
}

# -------------------------
# IAM Role and Trusted Entities Configuration
# -------------------------
variable "trusted_entities" {
  description = "A list of AWS services or accounts that are allowed to assume the IAM role. Common values include services like 'ec2.amazonaws.com' for EC2 or specific AWS account IDs."
  type        = list(string)
  default     = []
}

# -------------------------
# AWS Account and Region
# -------------------------
variable "account_id" {
  description = "The AWS account ID where resources will be deployed."
  type        = string
}

variable "region" {
  description = "The AWS region in which to deploy all resources."
  type        = string
}
