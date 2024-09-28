module "terraplay-vpc" {
  count  = var.create_networking ? 1 : 0
  source = "./modules/networking"
  vpc_config = {
    cidr_block = var.vpc_cidr
    name       = "${var.project_name}-vpc"
  }
  subnet_config = var.subnet_config
}

module "game-server-firewall" {
  count         = var.create_firewall ? 1 : 0
  source        = "./modules/firewall"
  vpc_id        = module.terraplay-vpc[0].vpc_id
  name          = "game-server-firewall"
  description   = var.security_group_description
  ingress_rules = var.ingress_rules
  egress_rules  = var.egress_rules
}

module "ssh_key" {
  count = var.create_key ? 1 : 0
  source            = "./modules/keys"
  private_key_path  = "${path.module}/sensitive/my_private_key.pem"
  key_pair_name = var.key_pair_name
}



