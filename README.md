## Terraplay

### Terraform Init Alias

To streamline Terraform initialization with dynamic backend configurations, add this alias to your shell configuration (`~/.bashrc` or `~/.zshrc`):

```bash
alias terraform-init='export WORKSPACE=$(terraform workspace show) && terraform init --backend-config="./env/backend.conf" --backend-config="key=terraform.tfstate"'
```

### Usage:
Run `terraform-init` to initialize Terraform with dynamic workspace-based state management.

### Example `backend.conf`:

```hcl
bucket         = "your-backend-bucket-name"
region         = "your-region"
dynamodb_table = "terraform-state-lock-table"
```

This configuration allows you to manage Terraform state files based on the current workspace without hardcoding the file paths.

### Workspace Naming Convention:

- The workspace name should follow the pattern: `(project-id):(environment)`
  - Example workspaces:
    - `terraplay:global`
    - `terraplay:minecraft`
    - `terraplay:terraria`

### Terraform Apply Alias:

To dynamically load `tfvars` files based on the current workspace, use this alias:

```bash
alias terraform-apply='ENV=$(terraform workspace show); ENV=${ENV##*@}; terraform apply -var-file=env/$ENV.tfvars -var-file=env/common/terraform.tfvars'
```

### Usage:
- Run `terraform-apply` to apply Terraform configurations with workspace-specific and common variables.

This alias extracts the environment part of the workspace (i.e., everything after the `:`) and uses it to dynamically load the appropriate `.tfvars` files.




### Temproray
alias terraform-apply='export BACKEND_BUCKET="terraplay-backend-80f0b90026287b08"; export AWS_REGION="ap-south-1"; export DYNAMODB_TABLE="terraform-state-lock"; ENV=$(terraform workspace show); ENV=${ENV##*:}; terraform apply -var-file=env/$ENV.tfvars -var-file=env/common/terraform.tfvars'

alias terraform-apply='terraform apply -var="backend_bucket=$BACKEND_BUCKET" -var="aws_region=$AWS_REGION" -var="dynamodb_table=$DYNAMODB_TABLE" -var-file=env/${ENV}.tfvars -var-file=env/common/terraform.tfvars'


alias terraform-init='export WORKSPACE=$(terraform workspace show) && terraform init --backend-config="./env/backend.conf" --backend-config="key=env/${WORKSPACE}/terraform.tfstate"'





alias terraform-refresh='ENV=$(terraform workspace show); ENV=${ENV##*@}; terraform refresh -var-file=env/$ENV.tfvars -var-file=env/common/terraform.tfvars'

alias terraform-destory='terraform destory -var="backend_bucket=$BACKEND_BUCKET" -var="aws_region=$AWS_REGION" -var="dynamodb_table=$DYNAMODB_TABLE" -var-file=env/${ENV}.tfvars -var-file=env/common/terraform.tfvars'


alias terraform-destroy='ENV=$(terraform workspace show); ENV=${ENV##*@}; terraform destroy -var-file=env/$ENV.tfvars -var-file=env/common/terraform.tfvars'














[reddit-aws-host](https://www.reddit.com/r/aws/comments/fss6nx/considering_using_aws_to_host_a_minecraft_server/)
[How to make cracked minecraft server](https://youtu.be/iJiTsM2MT3c)


maybe i will use secrets-manager!

oh oh oh, i can create an oidc for the github actons and connect this to my AWS!
I can write one script to automate it powerfully! put it into scripts-pool!


ngrok config add-authtoken 2ml9xN7HSrWB9ZKSh5s4BVdxQwC_4BDYRg3i3QZdwbSsXrvYY

https://github.com/orgs/community/discussions/9752
ngrok http 8080 














// hardcoding value
// oidc