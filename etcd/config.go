package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	client "github.com/coreos/etcd/client"
)

type Config struct {
	Endpoints  string `mapstructure:"endpoints"`
	Scheme     string `mapstructure:"scheme"`
	Keyfile    string `mapstructure:"keyfile"`
	Certfile   string `mapstructure:"certfile"`
	CACertfile string `mapstructure:"cacertfile"`
}

// Client() returns a new client for accessing etcd.
func (cfg *Config) Client() (client.KeysAPI, error) {

	log.Printf("[INFO] Etcd configured with endpoints: %s", cfg.Endpoints)

	endpoints := strings.Split(cfg.Endpoints, ",")

	var transport client.CancelableTransport
	switch cfg.Scheme {
	case "http":
		transport = newHTTPTransport()
	case "https":
		transport = newHTTPSTransport(cfg.Certfile, cfg.Keyfile, cfg.CACertfile)
	}

	config := client.Config{
		Endpoints: endpoints,
		Transport: transport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(config)
	if err != nil {
		log.Fatal(err)
	}

	return client.NewKeysAPI(c), nil
}

func newHTTPTransport() client.CancelableTransport {
	var defaultTransport client.CancelableTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return defaultTransport
}

func newHTTPSTransport(tlsCertFile, tlsKeyFile, tlsCACertFile string) client.CancelableTransport {
	var cc *tls.Config = nil

	if tlsCertFile != "" && tlsKeyFile != "" {
		var rpool *x509.CertPool
		if tlsCACertFile != "" {
			if pemBytes, err := ioutil.ReadFile(tlsCACertFile); err == nil {
				rpool = x509.NewCertPool()
				rpool.AppendCertsFromPEM(pemBytes)
			}
		}

		if tlsCert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile); err == nil {
			cc = &tls.Config{
				RootCAs:            rpool,
				Certificates:       []tls.Certificate{tlsCert},
				InsecureSkipVerify: true,
			}
		}
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     cc,
	}

	return tr
}
