# 1. VPC ID
# 2. Public subnets - subnet_key => { subnet_id, availability_zone }
# 3. Private subnets - subnet_key => { subnet_id, availability_zone }

locals {
  output_public_subnets = var.create_networking ? {
    for key in keys(local.public_subnets) : key => {
      subnet_id         = aws_subnet.this[key].id
      availability_zone = aws_subnet.this[key].availability_zone
    }
  } : {}

  output_private_subnets = var.create_networking ? {
    for key in keys(local.private_subnets) : key => {
      subnet_id         = aws_subnet.this[key].id
      availability_zone = aws_subnet.this[key].availability_zone
    }
  } : {}
}


output "vpc_id" {
  description = "The AWS ID from the created VPC"
  value       = var.create_networking ? aws_vpc.this[0].id : null
}

output "public_subnets" {
  description = "The ID and the availability zone of public subnets."
  value       = var.create_networking ? local.output_public_subnets : {}
}

output "private_subnets" {
  description = "The ID and the availability zone of private subnets."
  value       = var.create_networking ? local.output_private_subnets : {}
}
