---
page_title: "Abrha: abrha_region"
subcategory: "Account"
---

# abrha_region

Get information on a single Abrha region. This is useful to find out 
what Vm sizes and features are supported within a region.

## Example Usage

```hcl
data "abrha_region" "frankfurt" {
  slug = "frankfurt"
}

output "region_name" {
  value = data.abrha_region.frankfurt.name
}
```

## Argument Reference

* `slug` - (Required) A human-readable string that is used as a unique identifier for each region.

## Attributes Reference

* `slug` - A human-readable string that is used as a unique identifier for each region.
* `name` - The display name of the region.
* `available` -	A boolean value that represents whether new Vms can be created in this region.
* `sizes` - A set of identifying slugs for the Vm sizes available in this region.
* `features` - A set of features available in this region.
