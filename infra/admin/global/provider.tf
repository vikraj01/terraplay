terraform {
  required_version = "~> 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    null = {
      source = "hashicorp/null"
      version = "3.2.3"
    }
  }

  backend "s3" {
    bucket         = "terraplay-keeper-backend-da28ee014ea0433f"
    region         = "ap-south-1"
    dynamodb_table = "terraform-state-lock"
    key            = "terraform.tfstate"
  }
}

provider "aws" {
  region = "ap-south-1"
  default_tags {
    tags = local.default_tags
  }
}
