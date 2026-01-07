variable "name" {
  description = "Name of the Cloud Function"
  type        = string
}

variable "description" {
  description = "Description of the Cloud Function"
  type        = string
  default     = ""
}

variable "entry_point" {
  description = "Entry point function name"
  type        = string
}

variable "source_dir" {
  description = "Path to the function source directory"
  type        = string
}

variable "runtime" {
  description = "Runtime for the Cloud Function"
  type        = string
  default     = "go125"
}

variable "region" {
  description = "GCP region for the function"
  type        = string
}

variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "bucket_name" {
  description = "Name of the GCS bucket for function source"
  type        = string
}

variable "service_account_email" {
  description = "Service account email for the function"
  type        = string
}

variable "invoker_service_account_email" {
  description = "Service account email that can invoke the function"
  type        = string
}

variable "environment_variables" {
  description = "Environment variables for the function"
  type        = map(string)
  default     = {}
}

variable "max_instance_count" {
  description = "Maximum number of function instances"
  type        = number
  default     = 2
}

variable "max_instance_request_concurrency" {
  description = "Maximum number of concurrent requests per instance"
  type        = number
  default     = 1
}

variable "available_memory" {
  description = "Available memory for the function"
  type        = string
  default     = "256M"
}

variable "timeout_seconds" {
  description = "Timeout in seconds for the function"
  type        = number
  default     = 60
}

variable "schedule_config" {
  description = "Optional Cloud Scheduler configuration"
  type = object({
    schedule    = string
    time_zone   = string
    description = string
    job_name    = string
  })
  default = null
}

