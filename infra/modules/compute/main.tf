resource "aws_instance" "this" {
  ami                         = var.ami
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  key_name                    = var.key_name
  associate_public_ip_address = var.associate_public_ip_address

  vpc_security_group_ids = var.security_groups

  user_data = var.user_data

  root_block_device {
    volume_size = var.volume_size
    volume_type = var.volume_type
  }

  tags = merge(
    {
      Name = var.instance_name
    },
    var.ec2_tags
  )
}

resource "aws_eip" "this" {
  instance = aws_instance.this.id
  count    = var.assign_eip ? 1 : 0
}
