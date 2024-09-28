resource "aws_instance" "this" {
  ami                         = var.ami
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  key_name                    = var.key_name
  associate_public_ip_address = var.associate_public_ip_address
  monitoring                  = var.monitoring_enabled
  vpc_security_group_ids      = var.security_groups

  # Pass optional user data for EC2 instance initialization
  user_data = var.user_data

  # Root block device (boot volume)
  root_block_device {
    volume_size = var.root_volume_size
    volume_type = var.root_volume_type
  }

  # Add additional EBS volumes
  dynamic "ebs_block_device" {
    for_each = var.ebs_volumes
    content {
      device_name = ebs_block_device.value.device_name
      volume_size = ebs_block_device.value.volume_size
      volume_type = ebs_block_device.value.volume_type
    }
  }

  # Metadata options (e.g., for IAM roles, security enhancements)
  metadata_options {
    http_tokens               = var.http_tokens
    http_put_response_hop_limit = var.http_put_response_hop_limit
    http_endpoint              = var.http_endpoint
  }

  # Credit specification for burstable instance types (e.g., t3, t2)
  credit_specification {
    cpu_credits = var.cpu_credits
  }

  tags = merge(
    {
      Name = var.instance_name
    },
    var.ec2_tags
  )
}

# Optionally assign an Elastic IP
resource "aws_eip" "this" {
  instance = aws_instance.this.id
  count    = var.assign_eip ? 1 : 0
}

