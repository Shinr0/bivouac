locals {
  labels = {
    env        = "platform-ai-demo"
    layer      = "bootstrap"
    managed-by = "terraform"
  }

  # APIs nécessaires à la couche éphémère. Elles restent activées entre deux démos
  # (disable_on_destroy = false) : l'activation est gratuite et lente.
  services = toset([
    "aiplatform.googleapis.com",
    "billingbudgets.googleapis.com",
    "cloudbilling.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "secretmanager.googleapis.com",
    "serviceusage.googleapis.com",
    "storage.googleapis.com",
  ])
}

data "google_project" "this" {
  project_id = var.project_id
}

resource "google_project_service" "this" {
  for_each = local.services

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}
