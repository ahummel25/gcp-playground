provider "google" {
  project = var.project_id
  region  = var.region
}

############################
##                        ##
##      Storage           ##
##                        ##
############################
module "storage" {
  source = "../../modules/storage"

  project_id = var.project_id
  region     = var.region
}

############################
##                        ##
##      Secrets           ##
##                        ##
############################
module "secrets" {
  source = "../../modules/secrets"

  github_token          = var.github_token
  service_account_email = var.service_account_email
}

############################
##                        ##
##   Cloud Functions      ##
##                        ##
############################
module "hello_world_function" {
  source = "../../modules/cloud-function"

  name                          = "hello-world"
  description                   = "A simple Hello World function"
  entry_point                   = "HelloHTTP"
  source_dir                    = "${path.root}/../../../functions/hello_world"
  runtime                       = "go125"
  region                        = var.region
  project_id                    = var.project_id
  bucket_name                   = module.storage.bucket_name
  service_account_email         = var.service_account_email
  invoker_service_account_email = var.service_account_email

  environment_variables = {
    PROJECT_ID = var.project_id
  }
}

module "git_commit_scheduler_function" {
  source = "../../modules/cloud-function"

  name                          = "git-commit"
  description                   = "Creates verified empty commits to GitHub repository on schedule"
  entry_point                   = "EmptyCommit"
  source_dir                    = "${path.root}/../../../functions/git_commit"
  runtime                       = "go125"
  region                        = var.region
  project_id                    = var.project_id
  bucket_name                   = module.storage.bucket_name
  service_account_email         = var.service_account_email
  invoker_service_account_email = var.service_account_email

  environment_variables = {
    GITHUB_REPO = "ahummel25/github-scheduler"
    PROJECT_ID  = var.project_id
  }

  schedule_config = {
    schedule    = "0 0 * * *"
    time_zone   = "UTC"
    description = "Triggers GitHub empty commit function once daily"
    job_name    = "git-commit-scheduler-job"
  }
}

