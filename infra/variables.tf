variable "credentials" {
  default = ""
}

variable "region" {
    description = "region to operate resources in"
}

variable "project_id" {
    description = "ID of the project where the bucket will be created"
}

variable "cloud_run_image" {
    description = "url of cloud run docker registry"
}

variable "aws_access_key" {
  description = "aws access key for dynamo db"
  default = ""
}

variable "aws_secret_key" {
  description = "aws access key for dynamo db"
  default = ""
}
