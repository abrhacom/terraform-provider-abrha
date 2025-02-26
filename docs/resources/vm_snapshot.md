---
page_title: "Abrha: abrha_vm_snapshot"
subcategory: "Backups & Snapshots"
---

# abrha\_vm\_snapshot

Provides a resource which can be used to create a snapshot from an existing Abrha VM.

## Example Usage

```hcl
resource "abrha_vm" "web" {
  name   = "web-01"
  size   = "deLinuxVPS4"
  image  = "ubuntu24-cloudinit-qcow2"
  region = "frankfurt"
}

resource "abrha_vm_snapshot" "web-snapshot" {
  vm_id = abrha_vm.web.id
  name  = "web-snapshot-01"
}


resource "abrha_vm" "from-snapshot" {
  image  = abrha_vm_snapshot.web-snapshot.id
  name   = "web-02"
  region = "frankfurt"
  size   = "deLinuxVPS6"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Vm snapshot.
* `vm_id` - (Required) The ID of the Vm from which the snapshot will be taken.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Vm snapshot.
* `created_at` - The date and time the Vm snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a Vm to be created based on this snapshot.
* `regions` - A list of Abrha region "slugs" indicating where the vm snapshot is available.
* `size` - The billable size of the Vm snapshot in gigabytes.


## Import

Vm Snapshots can be imported using the `snapshot id`, e.g.

```
terraform import abrha_vm_snapshot.mysnapshot 12345
```
