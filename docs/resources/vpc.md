---
page_title: "Abrha: abrha_vpc"
subcategory: "Networking"
---

# abrha_vpc

Provides a [Abrha VPC](https://docs.parspack.com/reference/api/#tag/VPCs) resource.

VPCs are virtual networks containing resources that can communicate with each
other in full isolation, using private IP addresses.

## Example Usage

```hcl
resource "abrha_vpc" "example" {
  name     = "exampleVpc"
  region   = "frankfurt"
  ip_range = "10.10.10.0"
}
```

### Resource Assignment

`abrha_vm` resources
may be assigned to a VPC by referencing its `id`. For example:

```hcl
resource "abrha_vpc" "example" {
  name   = "exampleVpc"
  region = "frankfurt"
}

resource "abrha_vm" "example" {
  name     = "example-01"
  size     = "deLinuxVPS4"
  image    = "ubuntu24-cloudinit-qcow2"
  region   = "frankfurt"
  vpc_uuid = abrha_vpc.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the VPC. Must be unique and contain alphanumeric characters only.
* `region` - (Required) The Abrha region slug for the VPC's location.
* `description` - (Optional) A free-form text field up to a limit of 255 characters to describe the VPC.
* `ip_range` - (Optional) The range of IP addresses for the VPC. Network ranges cannot overlap with other networks in the same account and must be in range of private addresses as defined in RFC1918. It may not be larger than `/16` or smaller than `/24`.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The unique identifier for the VPC.
* `default` - A boolean indicating whether or not the VPC is the default one for the region.
* `created_at` - The date and time of when the VPC was created.

## Import

A VPC can be imported using its `id`, e.g.

```
terraform import abrha_vpc.example 12345
```
