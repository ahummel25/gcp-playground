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

variable "git_commit_mode" {
  description = "Git commit function mode: 'scheduled' (once daily, always commit) or 'random' (4x daily, state-driven commit/skip)"
  type        = string
  default     = "scheduled"
}

variable "environment" {
  description = "Environment name (e.g. prod, staging)"
  type        = string
}

variable "project_name" {
  description = "GCP project name or ID"
  type        = string
}

variable "tags" {
  description = "Common tags applied to resources"
  type        = map(string)
  default     = {}
}

