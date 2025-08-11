---
page_title: "Abrha: abrha_vm_snapshot"
subcategory: "Backups & Snapshots"
---

# abrha\_vm\_snapshot

Vm snapshots are saved instances of a Vm. Use this data
source to retrieve the ID of a Abrha Vm snapshot for use in other
resources.

## Example Usage

Get the Vm snapshot:

```hcl
data "abrha_vm_snapshot" "web-snapshot" {
  name_regex  = "^web"
  region      = "frankfurt"
  most_recent = true
}
```

Create vm from snapshot:

```hcl
data "abrha_vm_snapshot" "web-snapshot" {
  name_regex  = "^web"
  most_recent = true
}

resource "abrha_vm" "from-snapshot" {
  image  = data.abrha_vm_snapshot.web-snapshot.id
  name   = "web-02"
  region = "frankfurt"
  size   = "deLinuxVPS4"
}
```


## Argument Reference

* `name` - (Optional) The name of the Vm snapshot.

* `name_regex` - (Optional) A regex string to apply to the Vm snapshot list returned by Abrha. This allows more advanced filtering not supported from the Abrha API. This filtering is done locally on what Abrha returns.

* `most_recent` - (Optional) If more than one result is returned, use the most recent Vm snapshot.

~> **NOTE:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single Vm snapshot ID only, or use `most_recent` to choose the most recent one.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Vm snapshot.
* `created_at` - The date and time the Vm snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a Vm to be created based on this Vm snapshot.
* `regions` - A list of Abrha region "slugs" indicating where the Vm snapshot is available.
* `vm_id` - The ID of the Vm from which the Vm snapshot originated.
* `size` - The billable size of the Vm snapshot in gigabytes.
