# GCP Playground

A Terraform-managed Google Cloud Platform playground demonstrating Cloud Functions, Secret Manager, Cloud Scheduler integrations and more.

## Project Overview

This project deploys Google Cloud Functions (Gen 2) with automated scheduling, secret management, and infrastructure-as-code best practices. The infrastructure is organized using reusable Terraform modules and supports multiple environments.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    GCP Infrastructure                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────┐       ┌───────────────────┐           │
│  │ Cloud Scheduler  │──────▶│ Cloud Function    │           │
│  │ (Daily @ UTC)    │       │ (git-commit)      │           │
│  └──────────────────┘       └────────┬──────────┘           │
│                                      │                      │
│                                      ▼                      │
│                            ┌──────────────────┐             │
│                            │ Secret Manager   │             │
│                            │ (github_token)   │             │
│                            └────────┬─────────┘             │
│                                      │                      │
│                                      ▼                      │
│                            ┌────────────────────┐           │
│                            │ GitHub GraphQL API │           │
│                            │ (Verified Commits) │           │
│                            └────────────────────┘           │
│                                                             │
│  ┌──────────────────┐      ┌───────────────────┐            │
│  │ Cloud Function   │      │ Cloud Function    │            │
│  │ (hello_world)    │      │ (git-commit)      │            │
│  └──────────────────┘      └───────────────────┘            │
│         │                            │                      │
│         └────────────┬───────────────┘                      │
│                      ▼                                      │
│              ┌──────────────┐                               │
│              │ GCS Bucket   │                               │
│              │ (functions)  │                               │
│              └──────────────┘                               │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
gcp-playground/
├── README.md                    # This file
├── environments/
│   └── prod/                    # Production environment configuration
│       ├── main.tf              # Infrastructure definition
│       ├── variables.tf         # Variable definitions
│       ├── outputs.tf          # Output values
│       ├── terraform.tf         # Backend & provider config
│       ├── terraform.tfvars    # Non-sensitive variables
│       └── secrets.tfvars      # Sensitive variables (gitignored)
├── modules/                     # Reusable Terraform modules
│   ├── cloud-function/          # Cloud Function deployment module
│   ├── storage/                 # GCS bucket module
│   └── secrets/                 # Secret Manager module
└── functions/                   # Cloud Function source code
    ├── hello_world/             # Simple Hello World function
    └── git_commit/              # Git commit function
```

## Infrastructure Components

### Storage Module
- **GCS Bucket**: Stores function source code archives
- **Location**: Region-specific (us-central1)

### Secrets Module
- **Secret Manager Secret**: Stores GitHub personal access token
- **IAM**: Grants service account access to secrets

### Cloud Functions
- **hello_world_function**: Simple HTTP function demonstrating basic functionality
- **git-commit**: Creates verified empty commits to GitHub repositories on schedule

### Cloud Scheduler
- **git-commit-job**: Triggers the git-commit function daily at midnight UTC

## Getting Started

### Prerequisites

- Google Cloud SDK (`gcloud`) installed and authenticated
- Terraform >= 1.0 installed
- Access to GCP project with appropriate permissions
- GitHub personal access token with repository write permissions

### Initial Setup

1. **Configure variables**:
   ```bash
   cd environments/prod
   cp terraform.tfvars.example terraform.tfvars  # If exists
   # Edit terraform.tfvars with your project details
   ```

2. **Configure secrets**:
   ```bash
   # Create secrets.tfvars (gitignored)
   echo 'github_token = "your_github_token_here"' > secrets.tfvars
   ```

3. **Initialize Terraform**:
   ```bash
   terraform init
   ```

4. **Plan and apply**:
   ```bash
   terraform plan -var-file=secrets.tfvars
   terraform apply -var-file=secrets.tfvars
   ```

## Remote State

State is stored remotely in Google Cloud Storage:
- **Bucket**: `hummelgcp-terraform-state`
- **Prefix**: `environments/prod`

This enables:
- Team collaboration
- State versioning and backup
- Secure state management

## Adding New Functions

To add a new Cloud Function:

1. Create function source code in `functions/your_function/`
2. Add module block in `environments/prod/main.tf`:
   ```hcl
   module "your_function" {
     source = "../../modules/cloud-function"
     
     name                       = "your-function"
     entry_point               = "YourHandler"
     source_dir                = "${path.root}/../../functions/your_function"
     # ... other required variables
   }
   ```

See `modules/cloud-function/README.md` for complete module documentation.

## Environment Variables

Functions can be configured with environment variables in the module block:

```hcl
environment_variables = {
  MY_VAR = "value"
  ANOTHER_VAR = "another_value"
}
```

## Secrets Management

Sensitive values are stored in Secret Manager and accessed at runtime:

- **GitHub Token**: Stored in `github_token` secret
- Functions access secrets using the Secret Manager API
- IAM permissions are automatically configured

## Monitoring

View function logs:
```bash
gcloud functions logs read git-commit --region=us-central1 --gen2
```

View scheduler job status:
```bash
gcloud scheduler jobs describe git-commit-scheduler-job --location=us-central1
```

## Troubleshooting

### State Lock Issues
If you encounter state lock errors:
```bash
terraform force-unlock <LOCK_ID>
```

### Function Not Found
Ensure you're running Terraform from `environments/prod/` directory.

### Secret Access Denied
Verify the service account has `roles/secretmanager.secretAccessor` permission.

## Contributing

1. Make changes to function code in `functions/`
2. Update Terraform configuration in `environments/prod/`
3. Test with `terraform plan`
4. Apply with `terraform apply`

## Resources

- [Cloud Functions Documentation](https://cloud.google.com/functions/docs)
- [Terraform GCP Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [GitHub GraphQL API](https://docs.github.com/en/graphql)

