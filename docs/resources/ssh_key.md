---
page_title: "Abrha: abrha_ssh_key"
subcategory: "Account"
---

# abrha\_ssh_key

Provides a Abrha SSH key resource to allow you to manage SSH
keys for Vm access. Keys created with this resource
can be referenced in your Vm configuration via their ID.

## Example Usage

```hcl
# Create a new SSH key
resource "abrha_ssh_key" "default" {
  name       = "Terraform Example"
  public_key = file("/Users/terraform/.ssh/id_rsa.pub")
}

# Create a new Vm using the SSH key
resource "abrha_vm" "web" {
  image    = "ubuntu-18-04-x64"
  name     = "web-1"
  region   = "frankfurt"
  size     = "deLinuxVPS4"
  ssh_keys = [abrha_ssh_key.default.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSH key for identification
* `public_key` - (Required) The public key. If this is a file, it
can be read using the file interpolation function

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the key
* `name` - The name of the SSH key
* `public_key` - The text of the public key
* `fingerprint` - The fingerprint of the SSH key

## Import

SSH Keys can be imported using the `ssh key id`, e.g.

```
terraform import abrha_ssh_key.mykey 263654
```
