package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/j33pguy/Unifi_Provider/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.Provider})
}
