module "game_server" {
  source          = "../modules/compute"
  count           = var.game != "" ? 1 : 0
  ami             = data.aws_ami.amazon_linux.id
  instance_name   = "${var.game}-server-${terraform.workspace}"
  instance_type   = var.instance_type
  subnet_id       = var.subnet_id
  security_groups = var.security_group_ids
  assign_eip      = var.assign_eip
  key_name        = var.key_name

  ebs_volumes = var.ebs_volumes

  user_data = file("game/config/${lower(var.game)}.sh")

  ec2_tags = merge(
    {
      Name = "${var.game}-server-${terraform.workspace}",
      Game = var.game,
      Type = "Compute"
    },
    var.ec2_tags
  )
}
