locals {
  storage_tags = {
    Type = "Storage"
  }

  dynamodb_tags = {
    Type = "Database"
  }
}


locals {
  default_tags = {
    Project   = var.project_name
    ManagedBy = var.managed_by
    WorkSpace = terraform.workspace
  }
}
