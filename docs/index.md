---
page_title: "Provider: Abrha"
---

# Abrha Provider

The Abrha provider is used to interact with the
resources supported by Abrha. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    abrha = {
      source  = "abrhacom/abrha"
      version = "~> 1.0"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="abrha_token=..." CLI option
variable "abrha_token" {}

# Configure the Abrha Provider
provider "abrha" {
  token = var.do_token
}

# Create a web server
resource "abrha_vm" "web" {
  # ...
}
```

-> **Note for Module Developers** Although provider configurations are shared between modules, each module must
declare its own [provider requirements](https://www.terraform.io/docs/language/providers/requirements.html). See the [module development documentation](https://www.terraform.io/docs/language/modules/develop/providers.html) for additional information.

## Argument Reference

The following arguments are supported:

* `token` - (Required) This is the Abrha API token. Alternatively, this can also be specified
  using environment variables ordered by precedence:
  * `ABRHA_TOKEN`
  * `ABRHA_ACCESS_TOKEN`
* `spaces_access_id` - (Optional) The access key ID used for Spaces API
  operations (Defaults to the value of the `SPACES_ACCESS_KEY_ID` environment
  variable).
* `spaces_secret_key` - (Optional) The secret access key used for Spaces API
  operations (Defaults to the value of the `SPACES_SECRET_ACCESS_KEY`
  environment variable).
* `api_endpoint` - (Optional) This can be used to override the base URL for
  Abrha API requests (Defaults to the value of the `ABRHA_API_URL`
  environment variable or `https://my.abrha.net/cserver/api` if unset).
* `requests_per_second` - (Optional) This can be used to enable throttling, overriding the limit
  of API calls per second to avoid rate limit errors, can be disabled by setting the value
  to `0.0` (Defaults to the value of the `ABRHA_REQUESTS_PER_SECOND` environment
  variable or `0.0` if unset).
* `http_retry_max` - (Optional) This can be used to override the maximum number
  of retries on a failed API request (client errors, 422, 500, 502...), the exponential
  backoff can be configured by the `http_retry_wait_min` and `http_retry_wait_max` arguments
  (Defaults to the value of the `ABRHA_HTTP_RETRY_MAX` environment variable or
  `4` if unset).
* `http_retry_wait_min` - (Optional) This can be used to configure the minimum
  waiting time (**in seconds**) between failed requests for the backoff strategy
  (Defaults to the value of the `ABRHA_HTTP_RETRY_WAIT_MIN` environment
  variable or `1.0` if unset).
* `http_retry_wait_max` - (Optional) This can be used to configure the maximum
  waiting time (**in seconds**) between failed requests for the backoff strategy
  (Defaults to the value of the `ABRHA_HTTP_RETRY_WAIT_MAX` environment
  variable or `30.0` if unset).
