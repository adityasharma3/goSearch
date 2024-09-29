package elasticsearch

import (
	"log"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	EsClient *elasticsearch.Client
	once     sync.Once
)

func InitializeElasticSearch() {
	once.Do(func() {
		config := elasticsearch.Config{
			CloudID: os.Getenv("ELASTICSEARCH_CLOUD_ID"),
			APIKey:  os.Getenv("ELASTICSEARCH_API_KEY"),
		}

		client, err := elasticsearch.NewClient(config)
		if err != nil {
			log.Fatalf("Error occurred while connecting to Elasticsearch client: %s", err)
		}
		EsClient = client
	})
}

func GetElasticClient() *elasticsearch.Client {
	if EsClient == nil {
		InitializeElasticSearch()
	}
	return EsClient
}
