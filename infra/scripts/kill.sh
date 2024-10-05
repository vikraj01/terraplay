WORKSPACE=$1
DEFAULT="default"

terraform workspace select $WORKSPACE
terraform destroy -auto-approve
terraform workspace select $DEFAULT

if [ "$WORKSPACE" != "$DEFAULT" ]; then
  terraform workspace delete $WORKSPACE
fi

echo "$WORKSPACE is successfully destroyed and deleted"
echo "workspace:$WORKSPACE"
