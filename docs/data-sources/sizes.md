---
page_title: "Abrha: abrha_sizes"
subcategory: "Vms"
---

# abrha_sizes

Retrieves information about the Vm sizes that Abrha supports, with
the ability to filter and sort the results. If no filters are specified, all sizes
will be returned.

## Example Usage

Most common usage will probably be to supply a size to Vm:

```hcl
data "abrha_sizes" "main" {
  filter {
    key    = "slug"
    values = ["deLinuxVPS4"]
  }
}

resource "abrha_vm" "web" {
  image  = "almalinux9-cloudinit-qcow2"
  name   = "web-1"
  region = "frankfurt"
  size   = element(data.abrha_sizes.main.sizes, 0).slug
}
```

The data source also supports multiple filters and sorts. For example, to fetch sizes with 1 or 2 virtual CPU that are available "sgp1" region, then pick the cheapest one:

```hcl
data "abrha_sizes" "main" {
  filter {
    key    = "vcpus"
    values = [1, 2]
  }

  filter {
    key    = "regions"
    values = ["frankfurt"]
  }

  sort {
    key       = "price_monthly"
    direction = "asc"
  }
}

resource "abrha_vm" "web" {
  image  = "almalinux9-cloudinit-qcow2"
  name   = "web-1"
  region = "frankfurt"
  size   = element(data.abrha_sizes.main.sizes, 0).slug
}
```

The data source can also handle multiple sorts. In which case, the sort will be applied in the order it is defined. For example, to sort by memory in ascending order, then sort by disk in descending order between sizes with same memory:

```hcl
data "abrha_sizes" "main" {
  sort {
    // Sort by memory ascendingly
    key       = "memory"
    direction = "asc"
  }

  sort {
    // Then sort by disk descendingly for sizes with same memory
    key       = "disk"
    direction = "desc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the sizes by this key. This may be one of `slug`,
  `regions`, `memory`, `vcpus`, `disk`, `transfer`, `price_monthly`,
  `price_hourly`, or `available`.
* `values` - (Required) Only retrieves sizes which keys has value that matches
  one of the values provided here.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the sizes by this key. This may be one of `slug`,
  `memory`, `vcpus`, `disk`, `transfer`, `price_monthly`, or `price_hourly`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.


## Attributes Reference

The following attributes are exported:

* `slug` - A human-readable string that is used to uniquely identify each size.
* `available` - This represents whether new Vms can be created with this size.
* `transfer` - The amount of transfer bandwidth that is available for Vms created in this size. This only counts traffic on the public interface. The value is given in terabytes. (For some sizes this value shows download bandwidth only, based on `transfer_type` field.)
* `price_monthly` - The monthly cost of Vms created in this size if they are kept for an entire month.
* `price_hourly` - The hourly cost of Vms created in this size as measured hourly.
* `memory` - The amount of RAM allocated to Vms created of this size. The value is measured in megabytes.
* `vcpus` - The number of CPUs allocated to Vms of this size.
* `disk` - The amount of disk space set aside for Vms of this size. The value is measured in gigabytes.
* `regions` - List of region slugs where Vms can be created in this size.
