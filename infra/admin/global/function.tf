resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action    = "sts:AssumeRole",
        Effect    = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_ssm_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMFullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_s3_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_ec2_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

data "archive_file" "lambda_code" {
  type        = "zip"
  source_file = "${path.root}/build/index.js"
  output_path = "${path.root}/lambda.zip"
}

resource "aws_lambda_function" "my_lambda" {
  function_name    = "terraplay-lambda"
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "index.handler"
  runtime          = "nodejs16.x"
  timeout          = 30
  memory_size      = 256

  filename         = "lambda.zip" 

  environment {
    variables = {
      EC2_INSTANCE_ID = "your-ec2-instance-id"
      BUCKET_NAME     = "your-s3-bucket-name"
    }
  }
}

resource "aws_lambda_function_url" "lambda_url" {
  function_name = aws_lambda_function.my_lambda.function_name
  authorization_type = "NONE"  # Use AWS_IAM for authenticated access
}

resource "aws_lambda_permission" "allow_public_invoke" {
  statement_id  = "AllowPublicInvoke"
  action        = "lambda:InvokeFunctionUrl"
  function_name = aws_lambda_function.my_lambda.function_name
  principal     = "*"
  function_url_auth_type = aws_lambda_function_url.lambda_url.authorization_type
}

output "lambda_function_url" {
  value = aws_lambda_function_url.lambda_url.function_url
}
