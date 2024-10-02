output "ecr_repository_uri" {
  value = module.ecr.ecr_repository_uri
}

output "ecr_repository_arn" {
  value = module.ecr.ecr_repository_arn
}

output "ecr_iam_role_arn" {
  value = module.ecr_iam_role.role_arn
}