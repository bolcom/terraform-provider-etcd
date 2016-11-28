package main

import (
	"github.com/bakins/terraform-provider-etcd/etcd"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: etcd.Provider,
	})
}
