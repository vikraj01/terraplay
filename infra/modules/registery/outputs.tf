output "ecr_repository_uri" {
  value = aws_ecr_repository.ecr_repo.repository_url
}

output "ecr_repository_arn" {
  value = aws_ecr_repository.ecr_repo.arn
}
