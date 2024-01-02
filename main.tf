terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.21.0"
    }
    rancher2 = {
      source  = "rancher/rancher2"
      version = "3.2.0"
    }
  }
}

provider "docker" {}

resource "docker_image" "rancher" {
  name         = "rancher/rancher:latest"
  keep_locally = false
}

resource "docker_container" "rancher" {
  image = docker_image.rancher.image_id
  name  = "rancher"
  ports {
    internal = 80
    external = 80
  }
  ports {
    internal = 443
    external = 443
  }


  privileged = true
  restart    = "unless-stopped"
  env        = ["CATTLE_BOOTSTRAP_PASSWORD=test"]
}

# Provider bootstrap config with alias
provider "rancher2" {
  alias = "bootstrap"

  api_url   = "https://localhost"
  bootstrap = true
  insecure  = true
}

# Create a new rancher2_bootstrap using bootstrap provider config
resource "rancher2_bootstrap" "admin" {
  depends_on = [docker_container.rancher]

  provider = rancher2.bootstrap

  initial_password = "test"
  password         = "rancherkindacringe"
  telemetry        = true
}

# Provider config for admin
provider "rancher2" {
  alias = "admin"

  api_url   = rancher2_bootstrap.admin.url
  token_key = rancher2_bootstrap.admin.token
  insecure  = true
}