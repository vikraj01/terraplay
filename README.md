## Terraplay

### Terraform Init Alias

To streamline Terraform initialization with dynamic backend configurations, add this alias to your shell configuration (`~/.bashrc` or `~/.zshrc`):

```bash
alias terraform-init='export WORKSPACE=$(terraform workspace show) && terraform init --backend-config="./env/backend.conf" --backend-config="key=${WORKSPACE}/terraform.tfstate"'
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
alias terraform-apply='ENV=$(terraform workspace show); ENV=${ENV##*:}; terraform apply -var-file=env/$ENV.tfvars -var-file=env/common/terraform.tfvars'
```

### Usage:
- Run `terraform-apply` to apply Terraform configurations with workspace-specific and common variables.

This alias extracts the environment part of the workspace (i.e., everything after the `:`) and uses it to dynamically load the appropriate `.tfvars` files.

