provider "google" {
  credentials = "${var.credentials}"
  project = "${var.project_id}"
  region  = "${var.region}"
  zone    = "us-central1-c"
}

provider "google-beta" {
  credentials = "${var.credentials}"
  project     = "${var.project_id}"
  region      = "${var.region}"
}

output "gcp_credentials" {
  value       = google.credentials
  description = "The credentials for authenticating with gcp."
  sensitive   = true
}

output "gcp_beta_credentials" {
  value       = google-beta.credentials
  description = "The credentials for authenticating with gcp."
  sensitive   = true
}

resource "google_cloud_run_service" "selfhydro-state" {
  name     = "selfhydro-state"
  location = "${var.region}"
  provider = "google-beta"
  metadata {
    namespace = "${var.project_id}"
  }

  spec {
    containers {
	     image = "${var.cloud_run_image}:${file("../../version/version")}"
       env {
         AWS_ACCESS_KEY_ID = "${var.aws_access_key}"
         AWS_SECRET_ACCESS_KEY = "${var.aws_secret_key}"
       }
	   }
  }
}

data "google_iam_policy" "public-access" {
  binding {
    role = "roles/run.invoker"

    members = [
      "allUsers",
    ]
  }
}

output "aws_access_key" {
  value = "${var.aws_access_key}"
  sensitive = true
}

output "aws_secret_key" {
  value = "${var.aws_secret_key}"
  sensitive = true
}
