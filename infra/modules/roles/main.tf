resource "aws_iam_role" "role" {
  name        = var.role_name
  description = var.role_description
  path        = var.role_path

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      for entity in var.trusted_entities : {
        Effect = "Allow"
        principal = {
          AWS = entity
        },
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "name" {
  count      = length(var.managed_policy_arns)
  role       = aws_iam_role.role.name
  policy_arn = var.managed_policy_arns[count.index]
}

resource "aws_iam_role_policy" "name" {
  for_each = var.inline_policies
  name     = each.key
  role     = aws_iam_role.role.name
  policy   = each.value
}
