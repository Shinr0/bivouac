output "tfstate_bucket" {
  description = "Bucket GCS à utiliser comme backend de infra/platform."
  value       = google_storage_bucket.tfstate.name
}

output "github_token_secret" {
  description = "Nom complet du secret Secret Manager contenant le token GitHub."
  value       = google_secret_manager_secret.github_token.id
}

output "budget_id" {
  description = "Identifiant du budget GCP."
  value       = google_billing_budget.project.id
}
