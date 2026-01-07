provider "google" {
  project = var.project_name
  region  = var.region
  default_labels = merge(var.tags, {
    project     = lower(var.project_name)
    environment = lower(var.environment)
    managed_by  = "terraform"
  })
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
##   Git Commit State     ##
##   (random mode only)   ##
############################
module "git_commit_state" {
  count  = var.git_commit_mode == "random" ? 1 : 0
  source = "../../modules/git-commit-state"

  project_id          = var.project_id
  region              = var.region
  object_user_members = [var.service_account_email]
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

  timeout_seconds = 120

  environment_variables = merge(
    {
      GITHUB_REPO = "ahummel25/github-scheduler"
      PROJECT_ID  = var.project_id
      COMMIT_MODE = var.git_commit_mode
    },
    var.git_commit_mode == "random" ? {
      GIT_COMMIT_STATE_BUCKET = module.git_commit_state[0].bucket_name
    } : {}
  )

  schedule_config = {
    schedule    = var.git_commit_mode == "random" ? "0 0,6,12,18 * * *" : "0 0 * * *"
    time_zone   = "UTC"
    description = var.git_commit_mode == "random" ? "Triggers GitHub empty commit 4x daily with random commit/skip" : "Triggers GitHub empty commit function once daily"
    job_name    = "git-commit-scheduler-job"
  }
}

