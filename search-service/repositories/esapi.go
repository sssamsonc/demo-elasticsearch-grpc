package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/siongui/gojianfan"
	"search-service/configs"
	"search-service/models"
	pb "search-service/protos/esapi"
	"search-service/utils"
	"strconv"
	"strings"
	"time"
)

type EsApiRepository struct {
	EsClient *elastic.Client
}

func (r EsApiRepository) SearchNews(ctx context.Context, excludeCategoryIds []string, keywordStr string, current int32, size int32) (*models.StructEsApiSearchNews, error) {
	searchSource := elastic.NewSearchSource()
	mustQuery := elastic.NewBoolQuery()

	//some kinds of Search Logic in ElasticSearch
	if len(excludeCategoryIds) > 0 {
		boolQuery := elastic.NewBoolQuery()
		for _, categoryId := range excludeCategoryIds {
			boolQuery.MustNot(elastic.NewTermQuery("category_id", categoryId))
		}
		mustQuery.Must(boolQuery)
	}

	mustQuery.Must(elastic.NewRangeQuery("publish_time").Lt(time.Now().UnixNano() / int64(time.Millisecond)))

	values := strings.Split(keywordStr, " ")
	boolQueryMulti := elastic.NewBoolQuery()
	for i := 0; i < len(values); i++ {
		text := values[i]
		if len(text) == 0 {
			continue
		}
		traditionalText := gojianfan.S2T(text)
		simplifiedText := gojianfan.T2S(text)

		utils.Logger.Debug(fmt.Sprintf("Translation: ORG:%v, TC:%v, SC:%v", text, traditionalText, simplifiedText))

		if traditionalText != "" && text != traditionalText {
			boolQueryMulti.Should(elastic.NewMatchPhraseQuery("title", traditionalText).Slop(2).Boost(5.0)).
				Should(elastic.NewMatchPhraseQuery("digest", traditionalText).Slop(2).Boost(2.0))
		}

		if simplifiedText != "" && text != simplifiedText {
			boolQueryMulti.Should(elastic.NewMatchPhraseQuery("title", simplifiedText).Slop(2).Boost(5.0)).
				Should(elastic.NewMatchPhraseQuery("digest", simplifiedText).Slop(2).Boost(2.0))
		}

		boolQueryMulti.Should(elastic.NewMatchPhraseQuery("title", text).Slop(2).Boost(5.0)).
			Should(elastic.NewMatchPhraseQuery("digest", text).Slop(2).Boost(2.0))
	}

	// Combine queries into a single BoolQuery container
	combinedQuery := elastic.NewBoolQuery()
	combinedQuery.Must(mustQuery)
	combinedQuery.Filter(boolQueryMulti)

	searchSource.Query(combinedQuery).
		Sort("publish_time", false).
		Sort("_score", false).
		Highlight(elastic.NewHighlight().
			Field("title").
			Field("video_title").
			PreTags(configs.Get("HIGHLIGHT_START_TAG")).
			PostTags(configs.Get("HIGHLIGHT_END_TAG"))).
		From(int((current - 1) * size)).
		Size(int(size))

	searchService := r.EsClient.Search().Index(configs.Get("ES_INDEX_NAME")).SearchSource(searchSource) //.Pretty(true)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}
	if searchResult == nil {
		return nil, errors.New("searchResult is nil")
	}

	utils.Logger.Debug(fmt.Sprintf("Query took %d milliseconds", searchResult.TookInMillis))
	utils.Logger.Debug(fmt.Sprintf("Found a total of %d items", searchResult.TotalHits()))

	var newsList []*pb.News

	if searchResult.TotalHits() > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			var news *pb.News
			err := json.Unmarshal(hit.Source, &news)
			if err != nil {
				utils.Logger.Error("error when unmarshalling News: " + err.Error())
				continue
			}

			//Highlight the keywords
			//utils.Logger.Debug(fmt.Sprintf("Highlight: %v", hit.Highlight))

			if highlight, ok := hit.Highlight["title"]; ok {
				news.Title = highlight[0]
			}

			newsList = append(newsList, news)
		}
	} else {
		// No hits
		utils.Logger.Info("No result")
	}

	return &models.StructEsApiSearchNews{
		News:     newsList,
		EsResult: searchResult,
	}, nil
}

func (r EsApiRepository) UpsertNews(ctx context.Context, news *pb.News) (*elastic.IndexResponse, error) {
	if news == nil {
		return nil, errors.New("news is nil")
	}

	if len(strconv.FormatInt(news.GetPublishTime(), 10)) <= 10 {
		//convert seconds to milliseconds
		news.PublishTime = news.GetPublishTime() * 1000
	}
	if len(strconv.FormatInt(news.GetUpdateTime(), 10)) <= 10 {
		//convert seconds to milliseconds
		news.UpdateTime = news.GetUpdateTime() * 1000
	}

	// Add a document to the index
	result, err := r.EsClient.Index().
		Index(configs.Get("ES_INDEX_NAME")).
		BodyJson(news).
		//Refresh("wait_for").
		Id(news.Id).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r EsApiRepository) DeleteNews(ctx context.Context, id string) (*elastic.DeleteResponse, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	// Delete the document
	result, err := r.EsClient.Delete().
		Index(configs.Get("ES_INDEX_NAME")).
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, errors.New("failed to delete news:" + id + "::" + err.Error())
	}

	if result.Result != "deleted" {
		return nil, errors.New("failed to delete news:" + id)
	}

	return result, nil
}
