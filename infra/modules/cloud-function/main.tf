resource "google_storage_bucket_object" "function_zip" {
  name   = local.filename_on_gcs
  bucket = var.bucket_name
  source = data.archive_file.function_archive.output_path
}

resource "google_cloudfunctions2_function" "function" {
  name        = var.name
  description = var.description != "" ? var.description : "Cloud Function: ${var.name}"
  location    = var.region

  build_config {
    runtime     = var.runtime
    entry_point = var.entry_point
    source {
      storage_source {
        bucket = var.bucket_name
        object = google_storage_bucket_object.function_zip.name
      }
    }
  }

  service_config {
    max_instance_count               = var.max_instance_count
    max_instance_request_concurrency = var.max_instance_request_concurrency
    available_memory                 = var.available_memory
    timeout_seconds                  = var.timeout_seconds
    service_account_email            = var.service_account_email
    environment_variables            = var.environment_variables
  }
}

resource "google_cloudfunctions2_function_iam_member" "invoker" {
  project        = google_cloudfunctions2_function.function.project
  location       = google_cloudfunctions2_function.function.location
  cloud_function = google_cloudfunctions2_function.function.name
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:${var.invoker_service_account_email}"
}

resource "google_cloud_run_service_iam_member" "cloud_run_invoker" {
  project  = google_cloudfunctions2_function.function.project
  location = google_cloudfunctions2_function.function.location
  service  = replace(var.name, "_", "-")
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.invoker_service_account_email}"
}

resource "google_cloud_scheduler_job" "scheduler" {
  count = var.schedule_config != null ? 1 : 0

  name             = var.schedule_config.job_name
  description      = var.schedule_config.description
  schedule         = var.schedule_config.schedule
  time_zone        = var.schedule_config.time_zone
  region           = var.region
  attempt_deadline = "320s"

  http_target {
    http_method = "GET"
    uri         = google_cloudfunctions2_function.function.url

    oidc_token {
      service_account_email = var.invoker_service_account_email
    }
  }
}

