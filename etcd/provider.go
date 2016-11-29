package etcd

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"scheme": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http",
				Description: "http or https",
			},
			"endpoints": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "127.0.0.1:2379",
				Description: "multiple etcd endpoints separated by comma",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Etcd username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Etcd password",
			},
			"keyfile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "File location keyfile",
			},
			"certfile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "File location certfile",
			},
			"cacertfile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "File location cacert",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"etcd_discovery": resourceEtcdDiscovery(),
			"etcd_keys":      resourceEtcdKeys(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	log.Printf("[INFO] Initializing etcd client")
	return config.Client()
}
