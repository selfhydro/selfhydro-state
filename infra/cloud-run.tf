resource "google_cloud_run_service" "selfhydro-state" {
  name     = "selfhydro-state"
  location = "${var.region}"
  provider = "google-beta"

  metadata {
    namespace = "${var.project_id}"
  }

  spec {
    containers {
	     image = "${var.cloud_run_image}"
	   }
  }
}
