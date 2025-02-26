---
page_title: "Abrha: abrha_vms"
subcategory: "Vms"
---

# abrha_vms

Get information on Vms for use in other resources, with the ability to filter and sort the results.
If no filters are specified, all Vms will be returned.

This data source is useful if the Vms in question are not managed by Terraform or you need to
utilize any of the Vms' data.

Note: You can use the [`abrha_vm`](vm) data source to obtain metadata
about a single Vm if you already know the `id` or unique `name` to retrieve.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter images.

For example to find all Vms with size `deLinuxVPS4`:

```hcl
data "abrha_vms" "small" {
  filter {
    key    = "size"
    values = ["deLinuxVPS4"]
  }
}
```

You can filter on multiple fields and sort the results as well:

```hcl
data "abrha_vms" "small-with-backups" {
  filter {
    key    = "size"
    values = ["deLinuxVPS4"]
  }
  filter {
    key    = "backups"
    values = ["true"]
  }
  sort {
    key       = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the Vms by this key. This may be one of `backups`, `created_at`, `disk`, `id`,
  `image`, `ipv4_address`, `ipv4_address_private`, `ipv6`, `ipv6_address`, `ipv6_address_private`, `locked`,
  `memory`, `name`, `price_hourly`, `price_monthly`, `private_networking`, `region`, `size`,
  `status`, `vcpus`, `volume_ids`, or `vpc_uuid`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves Vms
  where the `key` field takes on one or more of the values provided here.
  
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.
 
`sort` supports the following arguments:

* `key` - (Required) Sort the Vms by this key. This may be one of `backups`, `created_at`, `disk`, `id`,
  `image`, `ipv4_address`, `ipv4_address_private`, `ipv6`, `ipv6_address`, `ipv6_address_private`, `locked`,
  `memory`, `name`, `price_hourly`, `price_monthly`, `private_networking`, `region`, `size`,
  `status`, `vcpus`, or `vpc_uuid`.

* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `vms` - A list of Vms satisfying any `filter` and `sort` criteria. Each Vm has the following attributes:  

  - `id` - The ID of the Vm.
  - `region` - The region the Vm is running in.
  - `image` - The Vm image ID or slug.
  - `size` - The unique slug that identifies the type of Vm.
  - `disk` - The size of the Vm's disk in GB.
  - `vcpus` - The number of the Vm's virtual CPUs.
  - `memory` - The amount of the Vm's memory in MB.
  - `price_hourly` - Vm hourly price.
  - `price_monthly` - Vm monthly price.
  - `status` - The status of the Vm.
  - `locked` - Whether the Vm is locked.
  - `ipv6_address` - The Vm's public IPv6 address
  - `ipv6_address_private` - The Vm's private IPv6 address
  - `ipv4_address` - The Vm's public IPv4 address
  - `ipv4_address_private` - The Vm's private IPv4 address
  - `backups` - Whether backups are enabled.
  - `ipv6` - Whether IPv6 is enabled.
  - `private_networking` - Whether private networks are enabled.
  - `vpc_uuid` - The ID of the VPC where the Vm is located.
