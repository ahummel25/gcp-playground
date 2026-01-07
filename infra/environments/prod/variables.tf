variable "project_id" {
  description = "The project ID to deploy the Cloud Function to"
  type        = string
}

variable "service_account_email" {
  description = "The service account email used by Terraform"
  type        = string
}

variable "region" {
  description = "The region to deploy the Cloud Function in"
  type        = string
  default     = "us-central1"
}

variable "db_password" {
  description = "The password for the default database user"
  type        = string
  sensitive   = true
  default     = ""
}

variable "github_token" {
  description = "GitHub personal access token for creating commits"
  type        = string
  sensitive   = true
  default     = ""
}

