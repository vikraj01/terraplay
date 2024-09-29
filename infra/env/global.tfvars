vpc_cidr = "10.0.0.0/16"
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

security_group_description = "This firewall for terraplay game servers"
create_firewall            = true
create_key = true
create_networking = true
key_pair_name= "terraplay-key-pair"





# Ingress rules to allow traffic for HTTP, SSH, and Minetest server (TCP and UDP)
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
  minetest_server_tcp = {
    from_port   = 30000
    to_port     = 30000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow Minetest TCP traffic from anywhere"
  },
  minetest_server_udp = {
    from_port   = 30000
    to_port     = 30000
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow Minetest UDP traffic from anywhere"
  }
}

# Egress rules to allow all outbound traffic
egress_rules = {
  all_outbound = {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }
}

