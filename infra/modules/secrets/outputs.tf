output "secret_id" {
  description = "ID of the GitHub token secret"
  value       = google_secret_manager_secret.github_token.secret_id
}

