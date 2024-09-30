#!/bin/bash

# Step 1: List all Terraform workspaces except 'default'
WORKSPACES=$(terraform workspace list | grep -v "default" | sed 's/*//g' | tr -d ' ')

# Step 2: If there are no other workspaces, exit
if [ -z "$WORKSPACES" ]; then
  echo "No workspaces to destroy, only 'default' exists."
  exit 0
fi

# Step 3: Loop through each workspace and destroy it
for WORKSPACE in $WORKSPACES; do
  echo "Processing workspace: $WORKSPACE"

  # Switch to the workspace
  echo "Switching to workspace: $WORKSPACE"
  terraform workspace select "$WORKSPACE"
  if [ $? -ne 0 ]; then
    echo "Error: Failed to switch to workspace $WORKSPACE"
    exit 1
  fi

  # Destroy resources in the workspace
  echo "Destroying resources in workspace: $WORKSPACE"
  terraform destroy -auto-approve
  if [ $? -ne 0 ]; then
    echo "Error: Failed to destroy resources in workspace $WORKSPACE"
    exit 1
  fi

  # Switch back to the default workspace
  echo "Switching back to default workspace"
  terraform workspace select default
  if [ $? -ne 0 ]; then
    echo "Error: Failed to switch back to default workspace"
    exit 1
  fi

  # Delete the workspace
  echo "Deleting workspace: $WORKSPACE"
  terraform workspace delete "$WORKSPACE"
  if [ $? -ne 0 ]; then
    echo "Error: Failed to delete workspace $WORKSPACE"
    exit 1
  fi

  echo "Workspace $WORKSPACE destroyed and deleted successfully!"
done

echo "All non-default workspaces have been destroyed and deleted."
