package databases

import (
	"errors"
	"github.com/olivere/elastic/v7"
	"os"
	"search-service/utils"
)

var elasticSearchClient *elastic.Client

var elasticSearchClientError error

var elasticSearchOnce utils.Once

func NewElasticSearchClient() (*elastic.Client, error) {
	//Perform connection creation operation only once.
	elasticSearchOnce.Do(func() {
		esClient, err := elastic.NewClient(
			elastic.SetURL(os.Getenv("ES_URL")),
			elastic.SetBasicAuth(os.Getenv("ES_NAME"), os.Getenv("ES_PASSWORD")),
			//elastic.SetSniff(false),
			//elastic.SetHealthcheck(false))
			elastic.SetSniff(false))

		if err != nil {
			elasticSearchClientError = err
			return
		}

		if esClient == nil {
			elasticSearchClientError = errors.New("esClient returned nil")
			return
		}

		elasticSearchClient = esClient
	})

	if elasticSearchClientError != nil {
		utils.Logger.Error("Failed to connect elastic search!")
		elasticSearchOnce.Reset()
	}

	return elasticSearchClient, elasticSearchClientError
}
