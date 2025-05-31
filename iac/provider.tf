terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "google" {
  credentials = file("sre-final-460613-79de80926799.json")
  project     = var.project_id
  region      = var.region
  zone        = "us-central1-a"
}
