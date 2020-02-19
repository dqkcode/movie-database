package elasticsearch

import (
	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/elastic/go-elasticsearch/v7"
)

type (
	config struct {
		Addrs []string `envconfig:"ELASTICSEARCH_URL" default:"http://localhost:9200"`
	}
)

func LoadConfigFromEnv() *config {
	c := &config{}
	envconfig.Load(c)
	return c

}
func NewClient() *elasticsearch.Client {
	// c := LoadConfigFromEnv()
	// cf := elasticsearch.Config{
	// 	Addresses: c.Addrs,
	// 	Transport: &http.Transport{
	// 		MaxIdleConnsPerHost:   10,
	// 		ResponseHeaderTimeout: time.Millisecond,
	// 		DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
	// 		TLSClientConfig: &tls.Config{
	// 			MinVersion: tls.VersionTLS11,
	// 		},
	// 	},
	// }
	// client, err := elasticsearch.NewClient(cf)

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logrus.Panicf("can not create elasticsearch client")
	}
	logrus.Info("connected to elasticsearch server")
	return client
}
