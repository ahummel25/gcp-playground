output "hello_world_function_url" {
  description = "The URL of the deployed Hello World Cloud Function"
  value       = module.hello_world_function.function_url
}

output "git_commit_function_url" {
  description = "The URL of the deployed Git Commit Cloud Function"
  value       = module.git_commit_scheduler_function.function_url
}

