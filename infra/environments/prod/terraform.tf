terraform {
  required_version = ">= 1.14.0"

  backend "gcs" {
    bucket = "hummelgcp-terraform-state"
    prefix = "environments/prod"
  }
}

module "shared" {
  source = "../../shared"

  region       = var.region
  environment  = var.environment
  project_name = var.project_name
  tags         = var.tags
}
