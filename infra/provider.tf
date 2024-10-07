terraform {
  required_version = "~> 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "terraplay-state-storage-v1c2946c2056f6cd0"
    region         = "ap-south-1"
    dynamodb_table = "terraform_state_lock"
    key            = "terraform.tfstate"
  }
}

provider "aws" {
  region = "ap-south-1"
  default_tags {
    tags = local.default_tags
  }
}
