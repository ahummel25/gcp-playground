resource "google_secret_manager_secret" "github_token" {
  secret_id = "github_token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "github_token" {
  secret      = google_secret_manager_secret.github_token.id
  secret_data = var.github_token
}

resource "google_secret_manager_secret_iam_member" "github_token_accessor" {
  secret_id = google_secret_manager_secret.github_token.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${var.service_account_email}"
}

