variable "region" {
  description = "The region to deploy the Cloud Function in"
  type        = string
  default     = "us-central1"
}

variable "db_password" {
  description = "The password for the default database user"
  type        = string
  sensitive   = true
}
