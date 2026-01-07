output "database_instance_connection_name" {
  value = google_sql_database_instance.postgres_instance.connection_name
}

output "dsn_name" {
  value = google_sql_database_instance.postgres_instance.dns_name
}

output "database_user" {
  value = google_sql_user.default.name
}

output "self_link" {
  value = google_sql_database_instance.postgres_instance.self_link
}
