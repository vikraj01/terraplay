output "security_group_id" {
  description = "The ID of the created security group"
  value       = aws_security_group.this.id
}

output "ingress_rule_ids" {
  description = "A list of IDs for the ingress rules created"
  value       = [for rule in aws_vpc_security_group_ingress_rule.ingress_rules : rule.id]
}

output "egress_rule_ids" {
  description = "A list of IDs for the egress rules created"
  value       = [for rule in aws_vpc_security_group_egress_rule.egress_rules : rule.id]
}