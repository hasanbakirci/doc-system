package elasticclient

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func ConnectElastic() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	res, err := es.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	return es, err
}
