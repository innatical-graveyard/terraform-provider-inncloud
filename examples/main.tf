terraform {
  required_providers {
    inncloud = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/edu/inncloud"
    }
  }
}

provider "inncloud" {
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2Mjc2ODA2MDgsImF1ZCI6InNreWUuaW5uYXRpY2FsLmNvbSIsInN1YiI6IjQ4YmQ1M2IwLTYyZTQtNGJmYy1iNWQ5LTQxZGVjMTg1ODhkNiJ9.IMRYSzRSNZAd-1op6v_Vz505YNN656qjoLadSa4jeoQ"
  project_id = "57372f3d-25c3-4041-a2a2-bfdc353b90c9" 
}

resource "inncloud_server" "UwU" {
  name = "UwU"
  model = "starter"
  image = "ubuntu-20.04"
  region = "HEL1"
  cycle = "month"
}

resource "inncloud_server" "fuck" {
  name = "fuck"
  model = "pro"
  image = "ubuntu-20.04"
  region = "HEL1"
  cycle = "month"
}

resource "inncloud_server" "yes" {
  name = "yes"
  model = "starter"
  image = "ubuntu-20.04"
  region = "HEL1"
  cycle = "month"
}