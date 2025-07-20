---
page_title: "Abrha: abrha_vm"
subcategory: "Vms"
---

# abrha\_vm

Provides a Abrha Vm resource. This can be used to create,
modify, and delete Vms. Vms also support
[provisioning](https://www.terraform.io/docs/language/resources/provisioners/syntax.html).

## Example Usage

```hcl
# Create a new Web VM in the Frankfurt region with weekly backup enabled.
resource "abrha_vm" "web" {
  image  = "ubuntu24-cloudinit-qcow2"
  name   = "web-1"
  region = "frankfurt"
  size   = "deLinuxVPS4"
  backups = true
  backup_policy {
     plan    = "weekly"
     weekday = "TUE"
  }
}
```
Or:

```hcl
# Create a new Web VM in the Frankfurt region with monthly backup enabled.
resource "abrha_vm" "web" {
  image  = "ubuntu24-cloudinit-qcow2"
  name   = "web-1"
  region = "frankfurt"
  size   = "deLinuxVPS4"
  backups = true
  backup_policy {
     plan    = "monthly"
     monthday = 15
  }
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The Vm image ID or slug. This could be either image ID or vm snapshot ID.
* `name` - (Required) The Vm name.
* `region` - The region where the Vm will be created.
* `size` - (Required) The unique slug that identifies the type of Vm. You can find a list of available slugs on [Abrha API documentation](https://docs.parspack.com/reference/api/#tag/Sizes).
* `backups` - (Optional) Boolean controlling if backups are made. Defaults to
   false.
* `backup_policy` - (Optional) An object specifying the backup policy for the Droplet. If omitted and `backups` is `true`, the backup plan will default to daily.
   - `plan` - The backup plan used for the Droplet. The plan can be either `daily`, `weekly` or `monthly`.
  - `weekday` - Specifies the day of the week (`SUN`, `MON`, `TUE`, `WED`, `THU`, `FRI`, `SAT`) when the backup will run, applicable only if the backup plan is set to `weekly`.
  - `monthday` - Specifies the day of the month (1–28) on which the backup will run, applicable only if the backup plan is set to `monthly`.
* `ipv6` - (Optional) Boolean controlling if IPv6 is enabled. Defaults to false.
  Once enabled for a VM, IPv6 can not be disabled. When enabling IPv6 on
  an existing VM, [additional OS-level configuration](https://docs.digitalocean.com/products/networking/ipv6/how-to/enable/#on-existing-droplets)
  is required.
* `vpc_uuid` - (Optional) The ID of the VPC where the Vm will be located.If no `vpc_uuid` is provided, the Vm will be placed in your account's default VPC for the region.
* `ssh_keys` - (Optional) A list of SSH key IDs or fingerprints to enable in
   the format `[12345, 123456]`. To retrieve this info, use the
   [Abrha API](https://docs.parspack.com/reference/api/#tag/SSH-Keys). Once a Vm is created keys can not
   be added or removed via this provider. Modifying this field will prompt you
   to destroy and recreate the Vm.
* `user_data` (Optional) - A string of the desired User Data for the Vm.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Vm
* `urn` - The uniform resource name of the Vm
* `name`- The name of the Vm
* `region` - The region of the Vm
* `image` - The image of the Vm
* `ipv4_address` - The IPv4 address
* `ipv4_address_private` - The private networking IPv4 address
* `locked` - Is the Vm locked
* `private_networking` - Is private networking enabled
* `price_hourly` - Vm hourly price
* `price_monthly` - Vm monthly price
* `size` - The instance size
* `disk` - The size of the instance's disk in GB
* `vcpus` - The number of the instance's virtual CPUs
* `status` - The status of the Vm

## Import

Vms can be imported using the Vm `id`, e.g.

```
terraform import abrha_vm.myvm 100823
```
