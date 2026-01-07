locals {
  functions_base_dir = dirname(var.source_dir)
}

# Archive that includes both the function directory (at root) and shared directory
# Cloud Functions Gen 2 expects function code at the root of the archive
# We include shared at the root level, and update the replace directive to use ./shared/logging
data "archive_file" "function_archive" {
  type        = "zip"
  output_path = "${var.source_dir}/function.zip"

  # Include all files from the function directory at the root of the archive
  # Exclude function.zip to avoid including the archive itself (binary file)
  dynamic "source" {
    for_each = [for f in fileset(var.source_dir, "**") : f if f != "function.zip"]
    content {
      content  = file("${var.source_dir}/${source.value}")
      filename = source.value
    }
  }

  # Include all files from the shared directory at root level (not in a subdirectory)
  # This way when extracted to /workspace/, shared/ will be at /workspace/shared/
  dynamic "source" {
    for_each = fileset("${local.functions_base_dir}/shared", "**")
    content {
      content  = file("${local.functions_base_dir}/shared/${source.value}")
      filename = "shared/${source.value}"
    }
  }

  # Override go.mod to use ./shared/logging instead of ../shared/logging for Cloud Build
  source {
    content = replace(
      file("${var.source_dir}/go.mod"),
      "replace github.com/hummelgcp/go/shared/logging => ../shared/logging",
      "replace github.com/hummelgcp/go/shared/logging => ./shared/logging"
    )
    filename = "go.mod"
  }
}

