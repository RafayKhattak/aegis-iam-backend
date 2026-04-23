output "api_public_ip" {
  description = "Public IP address of the Aegis API EC2 instance."
  value       = aws_instance.api_server.public_ip
}

output "api_url" {
  description = "Swagger UI URL for the deployed API."
  value       = "http://${aws_instance.api_server.public_ip}:8080/swagger/index.html"
}

output "rds_endpoint" {
  description = "Endpoint address for the PostgreSQL RDS instance."
  value       = aws_db_instance.postgres.endpoint
}
