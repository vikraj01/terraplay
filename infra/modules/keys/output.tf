output "public_key" {
  value       = tls_private_key.ssh_key.public_key_openssh
  description = "The public key in OpenSSH format"
}

output "key" {
    value = tls_private_key.ssh_key.id
    description = "This is the ssh key id"
}

output "aws_key_name" {
  value = aws_key_pair.ssh_key_pair.key_name
}