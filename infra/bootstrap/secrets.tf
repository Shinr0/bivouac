# Le secret est créé vide : sa valeur est ajoutée à la main (README, prérequis).
# Terraform ne gère jamais de valeur de secret, rien ne fuit dans le state.
resource "google_secret_manager_secret" "github_token" {
  secret_id = "github-token"

  replication {
    auto {}
  }

  depends_on = [google_project_service.this]
}
