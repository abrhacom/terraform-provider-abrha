---
page_title: "Abrha: abrha_vm"
subcategory: "Vms"
---

# abrha_vm

Get information on a Vm for use in other resources. This data source provides
all of the Vm's properties as configured on your Abrha account. This
is useful if the Vm in question is not managed by Terraform or you need to
utilize any of the Vm's data.

## Example Usage

Get the Vm by name:

```hcl
data "abrha_vm" "example" {
  name = "web"
}

output "vm_output" {
  value = data.abrha_vm.example.ipv4_address
}
```

## Argument Reference

One of the following arguments must be provided:

* `id` - (Optional) The ID of the Vm
* `name` - (Optional) The name of the Vm.


## Attributes Reference

The following attributes are exported:

* `id`: The ID of the Vm.
* `region` - The region the Vm is running in.
* `image` - The Vm image ID or slug.
* `size` - The unique slug that identifies the type of Vm.
* `disk` - The size of the Vms disk in GB.
* `vcpus` - The number of the Vms virtual CPUs.
* `memory` - The amount of the Vms memory in MB.
* `price_hourly` - Vm hourly price.
* `price_monthly` - Vm monthly price.
* `status` - The status of the Vm.
* `locked` - Whether the Vm is locked.
* `ipv6_address` - The Vms public IPv6 address
* `ipv6_address_private` - The Vms private IPv6 address
* `ipv4_address` - The Vms public IPv4 address
* `ipv4_address_private` - The Vms private IPv4 address
* `backups` - Whether backups are enabled.
* `private_networking` - Whether private networks are enabled.
* `vpc_uuid` - The ID of the VPC where the Vm is located.
