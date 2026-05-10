package repository

import (
	"bytes"
	"context"
	"gitlab.com/spitertech/recommender/models"

	"github.com/elastic/go-elasticsearch"
	recomError "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/repository/es_query_builder"
	"gitlab.com/spitertech/recommender/repository/es_query_excecutor"
	"gitlab.com/spitertech/recommender/repository/helper"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
	"log"
)

type ElasticSearchRepository interface {
	Search(ctx context.Context, searchString string) (interface{}, error)
	SearchSuggestions(ctx context.Context, request models.SearchSuggestionRequest) (models.SuggestionData, error)
}

type elasticSearchRepository struct {
	esClient *elasticsearch.Client
	esHelper helper.ElasticSearchHelper
}

func NewElasticSearchRepository(elsClient *elasticsearch.Client, esHelper helper.ElasticSearchHelper) ElasticSearchRepository {
	return &elasticSearchRepository{esClient: elsClient, esHelper: esHelper}
}

func (repo elasticSearchRepository) Search(ctx context.Context, searchString string) (interface{}, error) {
	logger := logUtils.GetLogger(ctx)
	r := map[string]interface{}{}

	var buf bytes.Buffer
	err := es_query_builder.NewElasticSearchQueryBuilder().
		WithMultiMatch(searchString, []string{"product_category", "product_name"}).
		BindQuery(&buf).
		WithSize(2).
		Error()
	if err != nil {
		logger.Errorf("ElasticSearchRepository.Search: Error in building query response. Error: %-v", err)
		return nil, err
	}

	err = es_query_excecutor.NewEsExecutor(ctx, repo.esClient).
		WithContext().
		WithIndex("product").
		WithBody(&buf).
		WithTrackTotalHits(true).
		WithPretty().
		WithResponseAs(&r).
		ExecuteSearchQuery()
	if err != nil {
		logger.Errorf("ElasticSearchRepository.Search: Error in getting response. Error: %+v", err)
		logger.Errorf("ElasticSearchRepository.Search: Error in getting response. Error: %+v", r)

		return nil, err
	}
	logger.Infof("ElasticSearchRepository.Search: response %+v ", r)
	return r, nil

}

func (repo elasticSearchRepository) SearchSuggestions(ctx context.Context, request models.SearchSuggestionRequest) (models.SuggestionData, error) {
	logger := logUtils.GetLogger(ctx)
	logger.Info("ElasticSearchRepository.SearchSuggestions: querying elastic search for suggestions")
	var resp = []models.Product{}

	var buf bytes.Buffer
	err := es_query_builder.NewElasticSearchQueryBuilder().
		WithSuggester(request.SearchText, []string{"product_name", "product_name._2gram", "product_name._3gram"}).
		BindQuery(&buf).
		WithSize(10).
		Error()
	if err != nil {
		logger.Errorf("ElasticSearchRepository.SearchSuggestions: Error in building query response. Error: %+v", err)
		return models.SuggestionData{}, err
	}

	err = es_query_excecutor.NewEsExecutor(ctx, repo.esClient).
		WithContext().
		WithIndex("product").
		WithBody(&buf).
		WithTrackTotalHits(true).
		WithPretty().
		ResponseAs(&resp).
		ExecuteSuggestQuery()

	if err != nil {
		logger.Errorf("ElasticSearchRepository.SearchSuggestions: Error in Error getting response. Error: %+v", err)
		return models.SuggestionData{}, err
	}
	logger.Infof("ElasticSearchRepository.SearchSuggestions: response  Of %+v ", resp)
	return models.SuggestionData{Product: resp}, nil

}

func (repo elasticSearchRepository) responseMapper(ctx context.Context, r map[string]interface{}) (map[string]interface{}, error) {
	logger := logUtils.GetLogger(ctx)
	var resp = map[string]interface{}{}
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		logger.Infof("%s", hit.(map[string]interface{}))
		resp[hit.(map[string]interface{})["_id"].(string)] = hit.(map[string]interface{})
	}
	if len(resp) == 0 {
		return nil, recomError.ErrElasticDocNotFound
	}

	return resp, nil
}
