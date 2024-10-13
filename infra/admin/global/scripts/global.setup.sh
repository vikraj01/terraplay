#!/usr/bin/bash

WORKSPACE="global"

terraform init
terraform workspace select "$WORKSPACE" || terraform workspace new "$WORKSPACE"
terraform apply -var-file="terraform.tfvars" -auto-approve