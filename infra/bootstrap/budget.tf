resource "google_billing_budget" "project" {
  billing_account = var.billing_account_id
  display_name    = "platform-ai-demo (${var.project_id})"

  budget_filter {
    projects = ["projects/${data.google_project.this.number}"]
  }

  amount {
    specified_amount {
      currency_code = var.budget_currency
      units         = tostring(var.budget_amount)
    }
  }

  dynamic "threshold_rules" {
    for_each = var.budget_thresholds
    content {
      threshold_percent = threshold_rules.value
      spend_basis       = "CURRENT_SPEND"
    }
  }

  # Les alertes partent par e-mail vers les administrateurs de facturation (comportement par défaut).

  depends_on = [google_project_service.this]
}
