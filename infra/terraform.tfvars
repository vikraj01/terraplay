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
