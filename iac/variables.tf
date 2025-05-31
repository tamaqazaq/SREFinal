variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP zone"
  type        = string
  default     = "us-central1-a"
}

variable "instance_name" {
  description = "Name for the VM instance"
  type        = string
  default     = "go-monolith-vm"
}

variable "db_password" {
  description = "Postgres DB password"
  type        = string
  sensitive   = true
}
