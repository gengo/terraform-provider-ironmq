package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func newProvider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"ironmq_queue": queueSchema(),
		},
		ConfigureFunc: configure,
	}
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: newProvider,
	})
}
