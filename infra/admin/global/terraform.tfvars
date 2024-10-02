project_name = "terraplay"
managed_by   = "Terraform"
vpc_cidr     = "10.0.0.0/16"
subnet_config = {
  public = {
    public     = true
    cidr_block = "10.0.1.0/24"
    az         = "ap-south-1a"
  }
  private = {
    cidr_block = "10.0.2.0/24"
    az         = "ap-south-1b"
  }
}
ingress_rules = {
  http = {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow HTTP traffic from anywhere"
  },
  ssh = {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow SSH access from anywhere (Adjust IPs if needed)"
  },
  game_server = {
    from_port   = 25565
    to_port     = 25565
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow Minecraft traffic from anywhere"
  }
}

egress_rules = {
  all_outbound = {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }
}

security_group_description = "This firewall for terraplay game servers"

table_name    = "SessionsTable"
hash_key      = "session_id"
range_key     = "user_id"
ttl_attribute = "ttl"

global_secondary_indexes = [
  {
    name               = "user_id-index"
    hash_key           = "user_id"
    range_key          = null
    range_key_type     = "S"
    hash_key_type      = "S"
    projection_type    = "ALL"
    non_key_attributes = [] 
    read_capacity      = 5
    write_capacity     = 5
  }
]

key_pair_name = "terraplay-key-pair"
instance_type="t2.micro"




trusted_entities=["ec2.amazonaws.com"]
ecr_repository_name = "global_terraplay_ecr"