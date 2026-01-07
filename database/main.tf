resource "google_sql_database_instance" "postgres_instance" {
  name                = "postgres-instance"
  database_version    = "POSTGRES_16"
  region              = var.region
  deletion_protection = true

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = true
      ssl_mode     = "ENCRYPTED_ONLY"

      dynamic "authorized_networks" {
        for_each = local.allowed_onprem_ips
        iterator = onprem

        content {
          name  = "onprem-${onprem.key}"
          value = onprem.value
        }
      }
    }

    insights_config {
      query_insights_enabled  = true
      query_plans_per_minute  = 5
      query_string_length     = 1024
      record_application_tags = true
      record_client_address   = true
    }
  }
}

resource "google_sql_user" "default" {
  name     = "postgres"
  instance = google_sql_database_instance.postgres_instance.name
  password = var.db_password
}
