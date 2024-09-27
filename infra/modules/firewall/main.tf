resource "aws_security_group" "this" {
  name        = var.name
  description = var.description
  vpc_id      = var.vpc_id

  tags = merge(
    var.sg_tags,
    {
      Name = var.name
    }
  )
}

resource "aws_vpc_security_group_egress_rule" "egress_rules" {
  for_each = { for rule_key, rule in var.egress_rules :
    rule_key => flatten([for cidr_block in rule.cidr_blocks : merge(rule, { cidr_block = cidr_block })])
  }
  security_group_id = aws_security_group.this.id
  ip_protocol       = each.value[0].protocol

  from_port = each.value[0].protocol != "-1" ? each.value[0].from_port : null
  to_port   = each.value[0].protocol != "-1" ? each.value[0].to_port : null

  cidr_ipv4   = each.value[0].cidr_block
  description = each.value[0].description
}

resource "aws_vpc_security_group_ingress_rule" "ingress_rules" {
  for_each = { for rule_key, rule in var.ingress_rules :
    rule_key => flatten([for cidr_block in rule.cidr_blocks : merge(rule, { cidr_block = cidr_block })])
  }
  security_group_id = aws_security_group.this.id
  ip_protocol       = each.value[0].protocol

  from_port = each.value[0].protocol != "-1" ? each.value[0].from_port : null
  to_port   = each.value[0].protocol != "-1" ? each.value[0].to_port : null

  cidr_ipv4   = each.value[0].cidr_block
  description = each.value[0].description
}
