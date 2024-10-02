locals {
  common_ingress_rules = {
    http = {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
      description = "Allow HTTP traffic from anywhere"
    }
    https = {
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
      description = "Allow HTTPS traffic from anywhere"
    }
    ssh = {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
      description = "Allow SSH access from anywhere"
    }
  }
  common_egress_rules = {
    all_outbound = {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
      description = "Allow all outbound traffic"
    }
  }
  firewall = yamldecode(file("${path.module}/../../games.yaml")).games

  firewall_rules = {
    for game, rules in local.firewall : game => {
      ingress_rules = merge(
        local.common_ingress_rules,
        { for rule in rules.ingress_rules : keys(rule)[0] => rule[keys(rule)[0]] }
      )
      egress_rules = merge(
        local.common_egress_rules,
      )
      description : rules.description
    }
  }


  games = [
    for key, _ in local.firewall : key
  ]
}

module "server_firewall" {
  for_each      = local.firewall_rules
  source        = "../../modules/firewall"
  vpc_id        = module.terraplay_vpc.vpc_id
  name          = "${each.key}-firewall-rule"
  description   = each.value.description
  ingress_rules = each.value.ingress_rules
  egress_rules  = each.value.egress_rules
}

module "bot_firewall" {
  source = "../../modules/firewall"
  vpc_id = module.terraplay_vpc.vpc_id
  name   = "nimbus-firewall"
  ingress_rules = local.common_ingress_rules
  egress_rules = local.common_egress_rules
}