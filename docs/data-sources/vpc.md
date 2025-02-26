---
page_title: "Abrha: abrha_vpc"
subcategory: "Networking"
---

# abrha_vpc

Retrieve information about a VPC for use in other resources.

This data source provides all of the VPC's properties as configured on your
Abrha account. This is useful if the VPC in question is not managed by
Terraform or you need to utilize any of the VPC's data.

VPCs may be looked up by `id` or `name`. Specifying a `region` will
return that that region's default VPC.

## Example Usage

### VPC By Name

```hcl
data "abrha_vpc" "example" {
  name = "example-network"
}
```

Reuse the data about a VPC to assign a Vm to it:

```hcl
data "abrha_vpc" "example" {
  name = "example-network"
}

resource "abrha_vm" "example" {
  name     = "example-01"
  size     = "deLinuxVPS4"
  image    = "ubuntu24-cloudinit-qcow2"
  region   = "frankfurt"
  vpc_uuid = data.abrha_vpc.example.id
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - The unique identifier of an existing VPC.
* `name` - The name of an existing VPC.
* `region` - The Abrha region slug for the VPC's location.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the VPC.
* `name` - The name of the VPC.
* `region` - The Abrha region slug for the VPC's location.
* `description` - A free-form text field describing the VPC.
* `ip_range` - The range of IP addresses for the VPC in CIDR notation.
* `default` - A boolean indicating whether or not the VPC is the default one for the region.
* `created_at` - The date and time of when the VPC was created.
