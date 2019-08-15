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

resource "google_cloud_run_service" "selfhydro-state" {
  name     = "selfhydro-state"
  location = "${var.region}"

  metadata {
    namespace = "${var.project_id}"
  }

  spec {
    containers {
	     image = "${var.cloud_run_image}"
	   }
  }
}
