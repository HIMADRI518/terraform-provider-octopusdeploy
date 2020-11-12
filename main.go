package main

import (
	"github.com/transactcampus/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: octopusdeploy.Provider})
}
