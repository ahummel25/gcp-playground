output "bucket_name" {
  description = "Name of the function storage bucket"
  value       = google_storage_bucket.function_bucket.name
}

