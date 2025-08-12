---
page_title: "Abrha: abrha_reserved_ip"
subcategory: "Networking"
---

# abrha\_reserved_ip

Provides a Abrha reserved IP to represent a publicly-accessible static IP addresses that can be mapped to one of your VMs.

~> **NOTE:** Reserved IPs can be assigned to a VM either directly on the `abrha_reserved_ip` resource by setting a `vm_id` or using the `abrha_reserved_ip_assignment` resource, but the two cannot be used together.

## Example Usage

```hcl
resource "abrha_vm" "example" {
  name               = "example"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "abrha_reserved_ip" "example" {
  vm_id  = abrha_vm.example.id
  region = abrha_vm.example.region
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region that the reserved IP is reserved to.
* `vm_id` - (Optional) The ID of VM that the reserved IP will be assigned to.

## Attributes Reference

The following attributes are exported:

* `ip_address` - The IP Address of the resource
* `urn` - The uniform resource name of the reserved ip

## Import

Reserved IPs can be imported using the `ip`, e.g.

```
terraform import abrha_reserved_ip.myip 192.168.0.1
```
