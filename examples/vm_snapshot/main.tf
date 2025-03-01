terraform {
  required_providers {
    abrha = {
      source  = "abrhacom/abrha"
      version = "~> 1.0"
    }
  }
}

variable "token" {}
variable "api_endpoint" {}
variable "region" { default = "frankfurt" }

data "abrha_images" "selectedImage" {
  filter {
    key      = "distribution"
    values   = ["ubuntu"]
    match_by = "exact"
  }
  filter {
    key      = "slug"
    values   = ["24-cloudinit"]
    match_by = "substring"
  }
  sort {
    key       = "created"
    direction = "desc"
  }
}

data "abrha_sizes" "selectedSize" {
  filter {
    key      = "regions"
    values   = [var.region]
  }
  filter {
    key      = "memory"
    values   = ["2048"]
  }
  filter {
    key      = "vcpus"
    values   = ["1"]
  }
  filter {
    key      = "disk"
    values   = ["100"]
  }
  filter {
    key      = "slug"
    values   = ["Linux"]
    match_by = "substring"
  }
  sort {
    key       = "price_monthly"
    direction = "asc"
  }
}

provider "abrha" {
  api_endpoint = var.api_endpoint
  token        = var.token
  # You need to create a terraform.tfvars file and set required variables in it
}

resource "abrha_vm" "web" {
  name   = "test-vm-terraform"
  size   = data.abrha_sizes.selectedSize.sizes[0].slug
  image  = element(data.abrha_images.selectedImage.images, 0).slug
  region = var.region
  backups = true
}

resource "abrha_vm_snapshot" "web-snapshot" {
  vm_id = abrha_vm.web.id
  name  = "test-vm-snapshot"
}
