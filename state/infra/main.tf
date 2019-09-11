provider "google" {
  credentials = "${var.credentials}"
  project = "selfhydro-197504"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_cloudfunctions_function" "TransferStateToDynamoDB" {
  name                  = "TransferStateToDynamoDB"
  region                = "us-central1"
  description           = "Send selfhydro state to dynamodb(aws)"
  runtime               = "go111"
  available_memory_mb   = 128
  entry_point           = "TransferStateToDynamoDB"
  labels = {
    project = "selfhydro"
  }

  event_trigger {
    event_type    = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource      = "telemetry-topic"
  }

  source_archive_bucket = "selfhydro-build-artifacts"
  source_archive_object = "selfhydro-state-release.zip"

  environment_variables = {
    AWS_ACCESS_KEY_ID = "${var.aws_access_key}"
    AWS_SECRET_ACCESS_KEY = "${var.aws_secret_key}"
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

output "gcp_credentials" {
  value       = google.credentials
  description = "The credentials for authenticating with gcp."
  sensitive   = true
}
