locals {
  public_subnets = {
    for key, config in var.subnet_config : key => config if config.public
  }

  private_subnets = {
    for key, config in var.subnet_config : key => config if !config.public
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "this" {
  count      = var.create_networking ? 1 : 0
  cidr_block = var.vpc_config.cidr_block
  tags = merge(var.tags, {
    Name = var.vpc_config.name
    Type = "Networking"
  })
}

resource "aws_subnet" "this" {
  for_each = var.create_networking ? var.subnet_config : {}
  vpc_id   = aws_vpc.this[0].id
  availability_zone = each.value.az
  cidr_block        = each.value.cidr_block
  tags = merge(var.tags, {
    Name = each.key
    Type = "Networking"
  })

  lifecycle {
    precondition {
      condition     = contains(data.aws_availability_zones.available.names, each.value.az)
      error_message = "Invalid AZ"
    }
  }

  depends_on = [aws_vpc.this]
}

resource "aws_internet_gateway" "this" {
  for_each = var.create_networking && length(local.public_subnets) > 0 ? { "default" = 1 } : {}
  vpc_id   = aws_vpc.this[0].id
  tags     = merge(var.tags, {
    Type = "Networking"
  })

  depends_on = [aws_vpc.this]
}

resource "aws_route_table" "public_rtb" {
  for_each = var.create_networking && length(local.public_subnets) > 0 ? { "default" = 1 } : {}
  vpc_id   = aws_vpc.this[0].id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this["default"].id
  }

  tags = merge(var.tags, {
    Type = "Networking"
  })

  depends_on = [aws_vpc.this, aws_internet_gateway.this]
}

resource "aws_route_table_association" "public_rtb" {
  for_each       = var.create_networking ? local.public_subnets : {}
  subnet_id      = aws_subnet.this[each.key].id
  route_table_id = aws_route_table.public_rtb["default"].id

  depends_on = [aws_vpc.this, aws_route_table.public_rtb]
}
