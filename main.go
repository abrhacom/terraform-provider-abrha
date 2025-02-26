package main

import (
	"github.com/abrhacom/terraform-provider-abrha/abrha"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: abrha.Provider})
}
