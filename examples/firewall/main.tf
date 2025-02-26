terraform {
  required_providers {
    abrha = {
      source  = "registry.terraform.io/abrha/abrha"
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

resource "abrha_vm" "web" {
  count = 3
  name   = "abrha-Firewall-${count.index+1}"
  size   = "deLinuxVPS4"
  image  = "ubuntu20-cloudinit-qcow2"
  region = "frankfurt"
}

locals {
  vm_ids = abrha_vm.web[*].id
}

resource "abrha_firewall" "web" {
  name = "only-22-80-and-443"

  vm_ids = local.vm_ids


  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["192.168.10.0/24"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0"]
  }

  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0"]
  }

  inbound_rule {
    protocol         = "udp"
    port_range       = "2080"
    source_addresses = ["10.10.20.0/24", "192.168.30.0/24", "10.11.20.0/24"]
  }




  outbound_rule {
    protocol              = "tcp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0"]
  }

}