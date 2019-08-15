provider "google" {
  credentials = "${var.credentials}"
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "us-central1-c"
}

output "gcp_credentials" {
  value       = google.credentials
  description = "The credentials for authenticating with gcp."
  sensitive   = true
}
