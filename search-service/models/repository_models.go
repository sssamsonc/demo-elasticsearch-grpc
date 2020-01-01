package models

import (
	"github.com/olivere/elastic/v7"
	"search-service/protos/esapi"
)

type StructEsApiSearchNews struct {
	News     []*esapi.News
	EsResult *elastic.SearchResult
}
