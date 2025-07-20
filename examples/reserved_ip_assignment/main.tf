terraform {
  required_providers {
    abrha = {
      source  = "registry.terraform.io/abrhacom/abrha"
      version = ">= 1.0.0"
    }
  }
}

variable "token" {}
variable "api_endpoint" {}

provider "abrha" {
  api_endpoint = var.api_endpoint
  token        = var.token
  # You need to create a terraform.tfvars file and set required variables in it
}

resource "abrha_reserved_ip" "example" {
  region = "frankfurt"
}

resource "abrha_vm" "example" {
  name               = "example"
  size               = "deVPS4"
  image              = "almalinux9-cloudinit-qcow2"
  region             = "frankfurt"
}

resource "abrha_reserved_ip_assignment" "example" {
  ip_address = abrha_reserved_ip.example.ip_address
  vm_id = abrha_vm.example.id
}