---
page_title: "Abrha: abrha_firewall"
subcategory: "Networking"
---

# abrha\_firewall

Provides a Abrha Cloud Firewall resource. This can be used to create,
modify, and delete Firewalls.

## Example Usage

```hcl
resource "abrha_vm" "web" {
  name   = "web-1"
  size   = "deLinuxVPS4"
  image  = "ubuntu24-cloudinit-qcow2"
  region = "frankfurt"
}

resource "abrha_firewall" "web" {
  name = "only-22-80-and-443"

  vm_ids = [abrha_vm.web.id]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["192.168.1.0/24", "10.20.0.0/24"]
  }
  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "443"
    destination_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "80"
    destination_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "53"
    destination_addresses = ["0.0.0.0/0"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Firewall name
* `vm_ids` (Optional) - The list of the IDs of the Vms assigned
  to the Firewall.
* `inbound_rule` - (Optional) The inbound access rule block for the Firewall.
  The `inbound_rule` block is documented below.
* `outbound_rule` - (Optional) The outbound access rule block for the Firewall.
  The `outbound_rule` block is documented below.

`inbound_rule` supports the following:

* `protocol` - (Required) The type of traffic to be allowed.
  This may be one of "tcp", "udp", or "icmp".
* `port_range` - (Optional) The ports on which traffic will be allowed
  specified as a string containing a single port, a range (e.g. "8000-9000"),
  or "1-65535" to open all ports for a protocol. Required for when protocol is
  `tcp` or `udp`.
* `source_addresses` - (Optional) An array of strings containing the IPv4
  addresses, IPv4 CIDRs from which the
  inbound traffic will be accepted.

`outbound_rule` supports the following:

* `protocol` - (Required) The type of traffic to be allowed.
  This may be one of "tcp", "udp", or "icmp".
* `port_range` - (Optional) The ports on which traffic will be allowed
  specified as a string containing a single port, a range (e.g. "8000-9000"),
  or "1-65535" to open all ports for a protocol. Required for when protocol is
  `tcp` or `udp`.
* `destination_addresses` - (Optional) An array of strings containing the IPv4
  addresses, IPv4 CIDRs to which the
  outbound traffic will be allowed.


## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Firewall.
* `status` - A status string indicating the current state of the Firewall.
  This can be "waiting", "succeeded", or "failed".
* `created_at` - A time value given in ISO8601 combined date and time format
  that represents when the Firewall was created.
* `pending_changes` - An list of object containing the fields, "vm_id",
  "removing", and "status".  It is provided to detail exactly which Vms
  are having their security policies updated.  When empty, all changes
  have been successfully applied.
* `name` - The name of the Firewall.
* `vm_ids` - The list of the IDs of the Vms assigned to
  the Firewall.
* `tags` - The names of the Tags assigned to the Firewall.
* `inbound_rule` - The inbound access rule block for the Firewall.
* `outbound_rule` - The outbound access rule block for the Firewall.

## Import

Firewalls can be imported using the firewall `id`, e.g.

```
terraform import abrha_firewall.myfirewall 12345
```
