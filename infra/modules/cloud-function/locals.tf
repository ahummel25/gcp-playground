locals {
  // Append the app hash to the filename as a temporary workaround for https://github.com/terraform-providers/terraform-provider-google/issues/1938
  filename_on_gcs = "${var.name}-${lower(replace(base64encode(data.archive_file.function_archive.output_md5), "=", ""))}.zip"
}

