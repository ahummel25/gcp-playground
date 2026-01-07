#!/bin/bash
# Terraform state migration script for refactored modules
# Run this script to move resources from old module structure to new one

set -e

echo "Starting Terraform state migration..."

# Storage module
echo "Migrating storage resources..."
terraform state mv -lock=false -lock=false 'module.gcp_function.google_storage_bucket.function_bucket' 'module.storage.google_storage_bucket.function_bucket' || echo "Storage bucket already moved or not found"

# Secrets module
echo "Migrating secrets resources..."
terraform state mv -lock=false 'module.gcp_function.google_secret_manager_secret.github_token' 'module.secrets.google_secret_manager_secret.github_token' || echo "Secret already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_secret_manager_secret_version.github_token' 'module.secrets.google_secret_manager_secret_version.github_token' || echo "Secret version already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_secret_manager_secret_iam_member.github_token_accessor' 'module.secrets.google_secret_manager_secret_iam_member.github_token_accessor' || echo "Secret IAM already moved or not found"

# Hello World Function
echo "Migrating hello_world function resources..."
terraform state mv -lock=false 'module.gcp_function.data.archive_file.hello_world_function_archive' 'module.hello_world_function.data.archive_file.function_archive' || echo "Hello world archive already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_storage_bucket_object.hello_world_function_zip' 'module.hello_world_function.google_storage_bucket_object.function_zip' || echo "Hello world zip already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloudfunctions2_function.hello_world' 'module.hello_world_function.google_cloudfunctions2_function.function' || echo "Hello world function already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloudfunctions2_function_iam_member.invoker' 'module.hello_world_function.google_cloudfunctions2_function_iam_member.invoker' || echo "Hello world IAM invoker already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloud_run_service_iam_member.cloud_run_invoker' 'module.hello_world_function.google_cloud_run_service_iam_member.cloud_run_invoker' || echo "Hello world Cloud Run IAM already moved or not found"
# Note: google_cloud_run_service_iam_binding.binding is not in new module (redundant with iam_member)
# This resource can be removed: terraform state rm 'module.gcp_function.google_cloud_run_service_iam_binding.binding'

# GitHub Scheduler Function
echo "Migrating github_scheduler function resources..."
terraform state mv -lock=false 'module.gcp_function.data.archive_file.github_scheduler_function_archive' 'module.github_scheduler_function.data.archive_file.function_archive' || echo "GitHub scheduler archive already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_storage_bucket_object.github_scheduler_function_zip' 'module.github_scheduler_function.google_storage_bucket_object.function_zip' || echo "GitHub scheduler zip already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloudfunctions2_function.github_scheduler' 'module.github_scheduler_function.google_cloudfunctions2_function.function' || echo "GitHub scheduler function already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloudfunctions2_function_iam_member.github_scheduler_invoker' 'module.github_scheduler_function.google_cloudfunctions2_function_iam_member.invoker' || echo "GitHub scheduler IAM invoker already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloud_run_service_iam_member.github_scheduler_cloud_run_invoker' 'module.github_scheduler_function.google_cloud_run_service_iam_member.cloud_run_invoker' || echo "GitHub scheduler Cloud Run IAM already moved or not found"
terraform state mv -lock=false 'module.gcp_function.google_cloud_scheduler_job.github_scheduler_job' 'module.github_scheduler_function.google_cloud_scheduler_job.scheduler[0]' || echo "GitHub scheduler job already moved or not found"

# Note: Some resources may not exist or may need different handling
echo ""
echo "Migration complete!"
echo ""
echo "Note: The following resources may need manual handling:"
echo "  - module.gcp_function.data.google_service_account.account (not in new structure - may need to remove)"
echo "  - module.gcp_function.google_project_iam_member.cloudfunctions_admin (not in new structure - may need to remove)"
echo "  - module.gcp_function.google_cloud_run_service_iam_binding.binding (redundant - can be removed)"
echo ""
echo "Run 'terraform plan' to verify the migration."

