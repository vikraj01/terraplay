# 1. Take the incoming input
# 2. Build and workspace name
# 3. Terraform Init Using Alias or Just With the script
# 4. Terraform Apply Using Alias or Just With the script


CONTEXT=$1

RANDOM_ID=$(openssl rand -base64 8 | tr -d '=+/')
WORKSPACE_NAME="$RANDOM_ID@$CONTEXT"

ls -la
pwd
echo "$CONTEXT"
echo "$WORKSPACE_NAME"

terraform init --backend-config="./env/backend.conf" --backend-config="key=terraform.tfstate"
terraform workspace select "$WORKSPACE_NAME" || terraform workspace new "$WORKSPACE_NAME"
terraform apply -var-file="env/${CONTEXT}.tfvars" -var-file="env/common/terraform.tfvars" -auto-approve
# terraform apply -var-file="env/${CONTEXT}.tfvars" -var-file="env/common/terraform.tfvars"
