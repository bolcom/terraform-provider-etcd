package main

import (
	"github.com/bolcom/terraform-provider-etcd/etcd"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: etcd.Provider,
	})
}
