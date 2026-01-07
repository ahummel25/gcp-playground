output "function_url" {
  description = "URL of the deployed Cloud Function"
  value       = google_cloudfunctions2_function.function.url
}

output "function_name" {
  description = "Name of the Cloud Function"
  value       = google_cloudfunctions2_function.function.name
}

