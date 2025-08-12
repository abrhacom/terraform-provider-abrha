---
page_title: "Abrha: abrha_reserved_ip"
subcategory: "Networking"
---

# abrha_reserved_ip

Get information on a reserved IP. This data source provides the region and VM id
as configured on your Abrha account. This is useful if the reserved IP
in question is not managed by Terraform or you need to find the VM the IP is
attached to.

An error is triggered if the provided reserved IP does not exist.

## Example Usage

Get the reserved IP:

```hcl
variable "public_ip" {}

data "abrha_reserved_ip" "example" {
  ip_address = var.public_ip
}

output "fip_output" {
  value = data.abrha_reserved_ip.example.vm_id
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The allocated IP address of the specific reserved IP to retrieve.

## Attributes Reference

The following attributes are exported:

* `region`: The region that the reserved IP is reserved to.
* `vm_id`: The VM id that the reserved IP has been assigned to.
