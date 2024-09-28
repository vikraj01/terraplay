locals {
  default_tags = {
    Project   = var.project_name
    ManagedBy = var.managed_by
    WorkSpace = terraform.workspace
  }
}
locals {
  workspace_name  = terraform.workspace
  split_workspace = split("@", local.workspace_name)
  workspace_game  = length(local.split_workspace) > 1 ? local.split_workspace[1] : ""
  valid_game      = contains(["minecraft", "terraria", "valheim", "minetest"], lower(local.workspace_game)) ? local.workspace_game : ""
}
