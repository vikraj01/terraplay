# Terraplay

Terraplay is a tool that combines a Discord bot (written in Go) and Terraform to manage game servers. It automates the process of provisioning and managing game environments by using Terraform workspaces, creating isolated infrastructure for each game session.

## Features

### Keeper Infrastructure
- **State Management**: Uses a remote state file to track the infrastructure.
- **Session Locking**: Prevents multiple operations on the same game environment by locking Terraform sessions.
- **Secure Authentication**: Uses **OpenID Connect (OIDC)** for GitHub Actions to avoid using hardcoded access keys.

### Global Infrastructure
- **Networking**: Uses a Virtual Private Cloud (VPC) and public subnet for game servers.
- **Game Configuration (`game.yml`)**: Stores game-specific configurations like firewall rules and storage paths.
- **SSH Access**: Single SSH key for all servers (can be customized for security).
- **Discord Bot VM**: Manages Terraform workflows and stores secrets securely.
- **Game Data Backup**: Automatically backs up game data to S3 after each session or when stopping the server.

### Game Infrastructure
Each game session is handled using a **Terraform workspace**. The workspace names follow the format `<random_id>@<game_name>`, allowing easy identification and parameterization for different games.

- **Dynamic Resource Allocation**: Networking, firewalls, and other resources are set up based on the game and session.
- **Custom Game Environments**: Each game environment is customized based on game-specific scripts.

## Commands

Terraplay is controlled via simple Discord commands:

1. **`!create <game>`**  
   Starts a new game session by provisioning the required infrastructure.
   
2. **`!list-sessions <all, running, halted, terminated>`**  
   Lists all game sessions, filtered by their current status.
   
3. **`!stop <sessionId>`**  
   Stops a game session and backs up the data to S3.
   
4. **`!destroy <sessionId>`**  
   Destroys the infrastructure related to a game session.
   
5. **`!restart <sessionId>`**  
   Restarts a stopped session, restoring game data from the backup.

6. **`!list-games`**  
   Shows all available games that can be provisioned.
