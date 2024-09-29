output "instance_id" {
  description = "The ID of the created EC2 instance"
  value       = module.game_server.instance_id
}

output "public_ip" {
  description = "The public IP of the EC2 instance"
  value       = module.game_server.public_ip
}

output "private_ip" {
  description = "The private IP of the EC2 instance"
  value       = module.game_server.private_ip
}