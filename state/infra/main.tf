provider "google" {
  credentials = "${var.credentials}"
  project = "selfhydro-197504"
  region  = "us-central1"
  zone    = "us-central1-c"
}

locals {
  function_name = "selfhydro-state-release"
  function_md5 = filemd5("../../../${var.function-local-directory}")
}

resource "google_storage_bucket" "bucket" {
  name = "cloud_function_source"
}

resource "google_storage_bucket_object" "archive" {
  //  append the app hash to the filename as a temporary workaround for https://github.com/terraform-providers/terraform-provider-google/issues/1938
  name   = "${local.function_name}-${lower(replace(base64encode(local.function_md5), "=", ""))}.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "../../../${var.function-local-directory}"
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

  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"

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
