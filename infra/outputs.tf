# # Output for public subnets, only if we are in the global workspace
# output "public_subnets" {
#   description = "A list of public subnets in the VPC"
#   value       = terraform.workspace == "terraplay@global" ? module.terraplay-vpc[0].public_subnets : null
# }

# # Output for security group ID, only if we are in the global workspace
# output "security_group_id" {
#   description = "The ID of the security group for the game server"
#   value       = terraform.workspace == "terraplay@global" ? module.game-server-firewall[0].security_group_id : null
# }

# # Output for VPC ID, only if we are in the global workspace
# output "vpc" {
#   description = "The VPC ID"
#   value       = terraform.workspace == "terraplay@global" ? module.terraplay-vpc[0].vpc_id : null
# }

# # Output for AWS key name, only if we are in the global workspace
# output "aws_key_name" {
#   description = "The AWS key name used for SSH"
#   value       = terraform.workspace == "terraplay@global" ? module.ssh_key[0].aws_key_name : null
# }
