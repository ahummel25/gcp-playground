variable "github_token" {
  description = "GitHub personal access token for creating commits"
  type        = string
  sensitive   = true
}

variable "service_account_email" {
  description = "Service account email that needs access to the secret"
  type        = string
}

