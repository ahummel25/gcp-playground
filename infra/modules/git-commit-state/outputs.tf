output "bucket_name" {
  description = "Name of the GCS bucket used for git-commit state"
  value       = google_storage_bucket.state.name
}
