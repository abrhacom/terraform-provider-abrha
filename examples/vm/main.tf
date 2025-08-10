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
variable "ssh_public_key_path" {}
variable "ssh_private_key_path" {}
variable "region" { default = "frankfurt" }

provider "abrha" {
  api_endpoint = var.api_endpoint
  token        = var.token
  # You need to create a terraform.tfvars file and set required variables in it
}

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
    values   = ["40"]
  }
  sort {
    key       = "price_monthly"
    direction = "asc"
  }
}

resource "abrha_ssh_key" "default" {
  name       = "AbrahSshExample"
  public_key = file(var.ssh_public_key_path)
}

resource "abrha_vm" "testVm" {
  image     = element(data.abrha_images.selectedImage.images, 0).slug
  region    = var.region
  size      = data.abrha_sizes.selectedSize.sizes[0].slug
  name      = "TestVmTerraform"
  user_data = "#!/bin/bash\\necho \"Hello, World!\" > /root/hello.txt"
  ssh_keys = [abrha_ssh_key.default.id]
  backups  = true

  connection {
    host        = self.ipv4_address
    type        = "ssh"
    private_key = file(var.ssh_private_key_path)
    user        = "root"
    timeout     = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # Install Apache
      "apt update",
      "apt -y install apache2",
      "systemctl start apache2"
    ]
  }
}
