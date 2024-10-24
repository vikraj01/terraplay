data "terraform_remote_state" "global" {
  backend = "s3"
  config = {
    bucket         = "terraplay-state-storage-v1c2946c2056f6cd0"
    region         = "ap-south-1"
    dynamodb_table = "terraform_state_lock"
    key            = "env:/global/terraform.tfstate"
  }
}

locals {
  global_vpc_id          = try(data.terraform_remote_state.global.outputs.vpc_id, null)
  global_public_subnet_0 = try(data.terraform_remote_state.global.outputs.subnets["public_subnets"].subnet_id, null)
  global_security_group  = try(data.terraform_remote_state.global.outputs.security_group_ids[local.valid_game], null)
  ssh_key_name           = try(data.terraform_remote_state.global.outputs.aws_key_name, null)
}

module "game_server" {
  source             = "./server"
  subnet_id          = local.global_public_subnet_0 != null ? local.global_public_subnet_0 : null
  security_group_ids = local.global_security_group != null ? [local.global_security_group] : []
  key_name           = local.ssh_key_name

  ebs_volumes = [
    {
      device_name = "/dev/sdf"
      volume_size = 100
      volume_type = "gp3"
    },
    {
      device_name = "/dev/sdg"
      volume_size = 200
      volume_type = "gp3"
    }
  ]
  game = local.valid_game
}


