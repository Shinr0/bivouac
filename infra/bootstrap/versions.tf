terraform {
  required_version = ">= 1.9"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "= 8.10.0"
    }
  }

  # State local, volontairement : voir docs/adr/0003-bootstrap-local-state.md
}
