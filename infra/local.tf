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
<<<<<<< Updated upstream
  valid_game      = contains(["minetest"], lower(local.workspace_game)) ? local.workspace_game : ""
=======
  valid_game      = contains(["minetest,terraria,valheim"], lower(local.workspace_game)) ? local.workspace_game : ""
>>>>>>> Stashed changes
}
