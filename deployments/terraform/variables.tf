variable "aws_region" {
  description = "AWS region for infrastructure deployment."
  type        = string
  default     = "us-east-1"
}

variable "db_username" {
  description = "PostgreSQL database username."
  type        = string
  default     = "aegis"
}

variable "db_password" {
  description = "PostgreSQL database password."
  type        = string
  default     = "supersecret123"
  sensitive   = true
}

variable "db_name" {
  description = "Application database name."
  type        = string
  default     = "aegis_iam"
}

variable "instance_type" {
  description = "EC2 instance type for API server."
  type        = string
  default     = "t3.micro"
}
