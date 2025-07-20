---
page_title: "Abrha: abrha_reserved_ip_assignment"
subcategory: "Networking"
---

# abrha\_reserved_ip_assignment

Provides a resource for assigning an existing Abrha reserved IP to a VM. This
makes it easy to provision reserved IP addresses that are not tied to the lifecycle of your
VM.

## Example Usage

```hcl
resource "abrha_reserved_ip" "example" {
  region = "nyc3"
}

resource "abrha_vm" "example" {
  name               = "baz"
  size               = "frankfurt"
  image              = "ubuntu24-cloudinit-qcow2"
  region             = "deLinuxVPS4"
}

resource "abrha_reserved_ip_assignment" "example" {
  ip_address = abrha_reserved_ip.example.ip_address
  vm_id = abrha_vm.example.id
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The reserved IP to assign to the VM.
* `vm_id` - (Optional) The ID of VM that the reserved IP will be assigned to.

## Import

Reserved IP assignments can be imported using the reserved IP itself and the `id` of
the VM joined with a comma. For example:

```
terraform import abrha_reserved_ip_assignment.foobar 192.0.2.1,2bff-7bf0-bdff-258a
```