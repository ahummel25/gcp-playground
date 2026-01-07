resource "google_storage_bucket" "function_bucket" {
  name     = "${var.project_id}-functions"
  location = var.region
}

