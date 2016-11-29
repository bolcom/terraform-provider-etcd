package etcd

import (
	"log"
	"time"

	client "github.com/coreos/etcd/client"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
}

// Client() returns a new client for accessing etcd.
//
func (cfg *Config) Client() (client.KeysAPI, error) {

	log.Printf("[INFO] Etcd configured with endpoints: %s", cfg.Endpoint)

	config := client.Config{
		Endpoints: []string{cfg.Endpoint},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// should we "ping" here to ensure we can communicate?

	return client.NewKeysAPI(c), nil
}
