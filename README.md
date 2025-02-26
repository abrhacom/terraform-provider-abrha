Abrha Terraform Provider
==================

- Documentation: https://registry.terraform.io/providers/abrha/abrha/latest/docs

Requirements
------------

-	[Terraform](https://developer.hashicorp.com/terraform/install) 0.10+
-	[Go](https://go.dev/doc/install) 1.14+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/abrhacom/terraform-provider-abrha`

```sh
$ mkdir -p $GOPATH/src/github.com/abrhacom; cd $GOPATH/src/github.com/abrhacom
$ git clone git@github.com:abrha/terraform-provider-abrha
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/abrhacom/terraform-provider-abrha
$ make build
```

Using the provider
----------------------

See the [Abrha Provider documentation](https://registry.terraform.io/providers/abrha/abrha/latest/docs) to get started using the Abrha provider.

Developing the Provider
---------------------------

See [CONTRIBUTING.md](./CONTRIBUTING.md) for information about contributing to this project.

Generate the Provider Documentation
-----------------------------------

Run the following command to generate the provider documentation:

```sh
$ go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
$ tfplugindocs generate --provider-name abrha --rendered-provider-name abrha
```
