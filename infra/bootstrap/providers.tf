provider "google" {
  project = var.project_id
  region  = var.region

  # Requis par l'API Billing Budgets avec des identifiants utilisateur (ADC),
  # sinon 403 : la consommation d'API est imputée à ce projet.
  billing_project       = var.project_id
  user_project_override = true

  default_labels = local.labels
}
