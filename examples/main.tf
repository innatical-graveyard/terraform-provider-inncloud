terraform {
  required_providers {
    inncloud = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/edu/inncloud"
    }
  }
}

variable "token" {
  type = string
}

variable "project_id" {
  type = string
}

provider "inncloud" {
  token = var.token
  project_id = var.project_id
}

resource "inncloud_server" "Example" {
  name = "Example"
  model = "starter"
  image = "ubuntu-20.04"
  region = "LA1"
  cycle = "month"
}