resource "random_id" "bucket_prefix" {
  byte_length = 8
}
module "backend_storage" {
  source      = "../../modules/storage"
  bucket_name = "${var.project_name}-backend-${random_id.bucket_prefix.hex}"

  tags = merge(local.storage_tags, {
    Name = "${var.project_name}-backend-${random_id.bucket_prefix.hex}"
  })
}

module "terraform_state_lock" {
  source         = "../../modules/database/dyanmodb"
  table_name     = var.table_name
  hash_key       = var.hash_key
  billing_mode   = var.billing_mode
  range_key      = var.range_key
  hash_key_type  = var.hash_key_type
  range_key_type = var.range_key_type

  tags = merge(local.dynamodb_tags, {
    Name = var.table_name
  })
}

module "bot_server" {
  source        = "../../modules/compute"
  ami           = data.aws_ami.amazon_linux.id
  instance_name = nimbus_bot_server
  instance_type = var.instance_type
  subnet_id     = ""
  security_groups = [  ]
  key_name = ""
  
}
# ource          = "../modules/compute"
#   count           = var.game != "" ? 1 : 0
#   ami             = data.aws_ami.amazon_linux.id
#   instance_name   = "${var.game}-server-${terraform.workspace}"
#   instance_type   = var.instance_type
#   subnet_id       = var.subnet_id
#   security_groups = var.security_group_ids
#   assign_eip      = var.assign_eip
#   key_name        = var.key_name

#   ebs_volumes = var.ebs_volumes

#   user_data = file("server/config/${lower(var.game)}.sh")

#   ec2_tags = merge(
#     {
#       Name = "${var.game}-server-${terraform.workspace}",
#       Game = var.game,
#       Type = "Compute"
#     },
#     var.ec2_tags
#   )