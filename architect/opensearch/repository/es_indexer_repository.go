package repository

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"gitlab.com/spitertech/recommender/repository/helper"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
	"strings"
)

type EsDocumentIndexerRepository interface {
	IndexDocument(ctx context.Context, doc interface{}) (map[string]interface{}, error)
}

type esDocumentIndexerRepository struct {
	esClient *elasticsearch.Client
	esHelper helper.ElasticSearchHelper
}

func NewEsDocumentIndexerRepository(esClient *elasticsearch.Client, esHelper helper.ElasticSearchHelper) EsDocumentIndexerRepository {
	return &esDocumentIndexerRepository{
		esHelper: esHelper,
		esClient: esClient,
	}
}
func (repo esDocumentIndexerRepository) IndexDocument(ctx context.Context, doc interface{}) (map[string]interface{}, error) {
	logger := logUtils.GetLogger(ctx)
	logger.Info("ElasticSearchRepository.IndexDocument.")

	stringDoc := repo.esHelper.DocumentMarshaller(ctx, doc)
	req := esapi.IndexRequest{
		Index:      "product",
		DocumentID: repo.esHelper.GenerateDocumentId(ctx),
		Body:       strings.NewReader(stringDoc),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, repo.esClient)
	if err != nil {
		logger.Errorf("ElasticSearchRepository.IndexDocument: Error in IndexRequest api. Error: %+v ", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Errorf("ElasticSearchRepository.IndexDocument: Error response from IndexRequest api. http_status_code: %+v. Error:%+v  ", res.Status(), res.Body)
		var resMap map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			logger.Errorf("ElasticSearchRepository.IndexDocument: Error in decoding response body. Error: %+v ", err)
			return nil, err
		}
		logger.Errorf("ElasticSearchRepository.IndexDocument: Error response from IndexRequest api. http_status_code: %+v. Error: %+v ", res.Status(), resMap)

		return nil, err
	}

	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("ElasticSearchRepository.IndexDocument: Error in decoding response body. Error: %+v ", err)
		return nil, err
	}

	logger.Infof("\n IndexRequest() RESPONSE: \n Status: %s \n Result: %-v  \n Version : %v  \n Response : %+v ", res.Status(), resMap["result"], int(resMap["_version"].(float64)), res)
	return resMap, nil

}
