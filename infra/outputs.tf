# output "public_subnets" {
#   description = "A list of public subnets in the VPC"
#   value       = module.terraplay-vpc[0].public_subnets
# }

# output "security_group_id" {
#   description = "The ID of the security group for the game server"
#   value       = module.game-server-firewall[0].security_group_id
# }

# output "sg_vpc_id" {
#   value = module.game-server-firewall[0].sg_vpc
# }

# output "vpc" {
#   value = module.terraplay-vpc[0].vpc_id
# }

# output "aws_key_name" {
#   value = module.ssh_key[0].aws_key_name
# }
