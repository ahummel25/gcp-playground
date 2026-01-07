data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "${var.source_dir}/function.zip"
}

