output "dev_ip" {
  description = "Public Ip address"
  value       = aws_instance.dev_node.public_ip
}
