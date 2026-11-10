variable "project_id" {
  description = "Identifiant du projet GCP qui héberge la plateforme."
  type        = string
}

variable "region" {
  description = "Région par défaut (cluster, bucket de state)."
  type        = string
  default     = "europe-west4"
}

variable "billing_account_id" {
  description = "Identifiant du compte de facturation (format XXXXXX-XXXXXX-XXXXXX) portant le budget."
  type        = string

  validation {
    condition     = can(regex("^[0-9A-F]{6}-[0-9A-F]{6}-[0-9A-F]{6}$", var.billing_account_id))
    error_message = "Format attendu : XXXXXX-XXXXXX-XXXXXX (hexadécimal majuscule)."
  }
}

variable "budget_amount" {
  description = "Montant mensuel du budget du projet, en unités entières de budget_currency."
  type        = number
  default     = 50

  validation {
    condition     = var.budget_amount > 0 && floor(var.budget_amount) == var.budget_amount
    error_message = "Le budget doit être un entier strictement positif."
  }
}

variable "budget_currency" {
  description = "Devise ISO 4217 du budget ; doit correspondre à la devise du compte de facturation."
  type        = string
  default     = "EUR"
}

variable "budget_thresholds" {
  description = "Seuils d'alerte du budget (fraction de la dépense réelle)."
  type        = list(number)
  default     = [0.5, 0.9, 1.0]
}
