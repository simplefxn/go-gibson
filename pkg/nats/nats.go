package nats

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

func CreateTlsConfiguration(conf *config.Nats) (t *tls.Config) {
	if conf.Cert != "" && conf.Key != "" || conf.CA != "" {
		cert, err := tls.LoadX509KeyPair(conf.Cert, conf.Key)
		if err != nil {
			logger.Log.Debug(err)
		}

		caCert, err := os.ReadFile(conf.CA)
		if err != nil {
			logger.Log.Debug(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: conf.VerifySSL,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
