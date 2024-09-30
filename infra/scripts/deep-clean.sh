#!/bin/bash

WORKSPACES=$(terraform workspace list | grep -v "default" | sed 's/*//g' | tr -d ' ')

if [ -z "$WORKSPACES" ]; then
  echo "No workspaces to destroy, only 'default' exists."
  exit 0
fi

for WORKSPACE in $WORKSPACES; do
  if [ "$WORKSPACE" = "global" ]; then
    echo "Skipping global workspace."
    continue
  fi

  echo "Processing workspace: $WORKSPACE"
  terraform workspace select "$WORKSPACE"
  if [ $? -ne 0 ]; then
    echo "Error: Failed to switch to workspace $WORKSPACE"
    exit 1
  fi

  echo "Destroying resources in workspace: $WORKSPACE"
  terraform destroy -auto-approve
  if [ $? -ne 0 ]; then
    echo "Error: Failed to destroy resources in workspace $WORKSPACE"
    exit 1
  fi

  echo "Switching back to default workspace"
  terraform workspace select default
  if [ $? -ne 0 ]; then
    echo "Error: Failed to switch back to default workspace"
    exit 1
  fi

  echo "Deleting workspace: $WORKSPACE"
  terraform workspace delete "$WORKSPACE"
  if [ $? -ne 0 ]; then
    echo "Error: Failed to delete workspace $WORKSPACE"
    exit 1
  fi

  echo "Workspace $WORKSPACE destroyed and deleted successfully!"
done

echo "All non-default, non-global workspaces have been destroyed and deleted."
