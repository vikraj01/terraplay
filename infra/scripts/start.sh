# 1. Take the incoming input
# 2. Build and workspace name
# 3. Terraform Init Using Alias or Just With the script
# 4. Terraform Apply Using Alias or Just With the script


GAME=$1
USER_ID=$2
WORKSPACE_NAME="$USER_ID@$GAME"

terraform init --backend-config="./env/backend.conf" --backend-config="key=terraform.tfstate"
terraform workspace select "$WORKSPACE_NAME" || terraform workspace new "$WORKSPACE_NAME"
terraform apply -var-file="env/${GAME}.tfvars" -var-file="env/common/terraform.tfvars" -auto-approve
terraform output -json