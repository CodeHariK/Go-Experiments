resource "aws_vpc" "mtc_vpc" {
  cidr_block           = "10.123.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "dev"
  }
}

resource "aws_subnet" "mtc_public_subnet" {
  vpc_id     = aws_vpc.mtc_vpc.id
  cidr_block = "10.123.1.0/24"

  map_public_ip_on_launch = true
  availability_zone       = "ap-south-1a"

  tags = {
    Name = "dev-public"
  }
}

resource "aws_internet_gateway" "mtc_internet_gateway" {
  vpc_id = aws_vpc.mtc_vpc.id

  tags = {
    Name = "dev-igw"
  }
}

resource "aws_route_table" "mtc_public_rt" {
  vpc_id = aws_vpc.mtc_vpc.id

  tags = {
    Name = "dev-public_rt"
  }
}

resource "aws_route" "default_route" {
  route_table_id         = aws_route_table.mtc_public_rt.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.mtc_internet_gateway.id
}

resource "aws_route_table_association" "mtc_public_assoc" {
  subnet_id      = aws_subnet.mtc_public_subnet.id
  route_table_id = aws_route_table.mtc_public_rt.id
}

resource "aws_security_group" "mtc_sg" {
  name        = "dev_sg"
  description = "dev security group"
  vpc_id      = aws_vpc.mtc_vpc.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "mtc_auth" {
  key_name   = "terrakey"
  public_key = file("~/.ssh/terrakey.pub")
}

resource "aws_instance" "dev_node" {
  ami                         = data.aws_ami.server_ami.id
  instance_type               = "t2.micro"
  key_name                    = aws_key_pair.mtc_auth.id
  vpc_security_group_ids      = [aws_security_group.mtc_sg.id]
  subnet_id                   = aws_subnet.mtc_public_subnet.id
  associate_public_ip_address = true

  # user_data              = file("userdata.tpl")

  user_data = <<-EOF
            #!/bin/bash
            yum install -y nginx
            systemctl start nginx
            systemctl enable nginx
            echo "<html><body><h1>Hello from Nginx on AWS!</h1></body></html>" > /usr/share/nginx/html/index.html
            EOF

  root_block_device {
    volume_size = 10
  }

  tags = {
    Name = "HelloWorld"
  }

  provisioner "local-exec" {
    command = templatefile("${var.host_os}-ssh-config.tpl", {
      hostname     = self.public_ip
      user         = "ec2-user"
      identityfile = "~/.ssh/terrakey"
    })
    interpreter = var.host_os == "linux" ? ["bash", "-c"] : ["Powershell", "-Command"]
  }
}
