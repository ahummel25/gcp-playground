resource "google_storage_bucket" "state" {
  name     = "${var.project_id}-git-commit-state"
  location = var.region
}

resource "google_storage_bucket_iam_member" "object_users" {
  for_each = toset(var.object_user_members)
  bucket   = google_storage_bucket.state.name
  role     = "roles/storage.objectUser"
  member   = "serviceAccount:${each.value}"
}
