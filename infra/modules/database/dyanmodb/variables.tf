variable "billing_mode" {
  type    = string
  default = "PROVISIONED"

  validation {
    condition     = contains(["PROVISIONED"], var.billing_mode)
    error_message = "The billing_mode must be either 'PROVISIONED' or 'PAY_PER_REQUEST'."
  }
}

variable "table_name" {
  description = "The name of the DynamoDB table."
  type        = string

  validation {
    condition     = length(var.table_name) > 0
    error_message = "The table_name must not be empty."
  }
}

variable "hash_key" {
  description = "The partition key for the DynamoDB table."
  type        = string

  validation {
    condition     = length(var.hash_key) > 0
    error_message = "The hash_key must not be empty."
  }
}

variable "hash_key_type" {
  description = "The data type for the hash key (S: string, N: number, B: binary)."
  type        = string
  default     = "S"

  validation {
    condition     = contains(["S", "N", "B"], var.hash_key_type)
    error_message = "The hash_key_type must be either 'S', 'N', or 'B'."
  }
}

variable "range_key" {
  description = "The optional sort key for the DynamoDB table."
  type        = string
  default     = null
}

variable "range_key_type" {
  description = "The data type for the range key (S: string, N: number, B: binary)."
  type        = string
  default     = "S"

  validation {
    condition     = contains(["S", "N", "B"], var.range_key_type)
    error_message = "The range_key_type must be either 'S', 'N', or 'B'."
  }
}

variable "read_capacity" {
  description = "The number of read capacity units for the table."
  type        = number
  default     = 5

  validation {
    condition     = var.read_capacity > 0
    error_message = "The read_capacity must be greater than 0."
  }
}

variable "write_capacity" {
  description = "The number of write capacity units for the table."
  type        = number
  default     = 5

  validation {
    condition     = var.write_capacity > 0
    error_message = "The write_capacity must be greater than 0."
  }
}

variable "ttl_attribute" {
  description = "The attribute for Time to Live (TTL) expiration."
  type        = string
  default     = null
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
  default = []

  validation {
    condition     = alltrue([for idx in var.global_secondary_indexes : contains(["ALL", "KEYS_ONLY", "INCLUDE"], idx.projection_type)])
    error_message = "Each global_secondary_index must have a valid projection_type (ALL, KEYS_ONLY, INCLUDE)."
  }
}

variable "tags" {
  description = "A map of tags to assign to the dyanmodb table."
  type        = map(string)
  default     = {}
}