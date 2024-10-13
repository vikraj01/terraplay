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

  games_file      = file("${path.module}/games.yaml")
  games_yaml      = yamldecode(local.games_file)

  valid_games     = [for game in local.games_yaml["games"] : lower(game.key)]
  valid_game      = contains(local.valid_games, lower(local.workspace_game)) ? local.workspace_game : ""
}

