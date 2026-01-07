variable "region" {
  description = "Region for shared infrastructure"
  type        = string
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
  description = "Common tags applied to shared resources"
  type        = map(string)
  default     = {}
}
