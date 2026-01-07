variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region for the bucket"
  type        = string
}

variable "object_user_members" {
  description = "List of service account emails to grant roles/storage.objectUser on the bucket (e.g. the git-commit function's service account)"
  type        = list(string)
}
