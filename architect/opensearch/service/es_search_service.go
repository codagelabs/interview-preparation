package service

import (
	"context"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/repository"
	"gitlab.com/spitertech/recommender/utils/log_utils"
)

type ElasticSearchService interface {
	Search(ctx context.Context, request models.SearchRequest) (interface{}, *error2.RECOMError)
	GenerateSearchSuggestion(ctx context.Context, request models.SearchSuggestionRequest) (models.SuggestionData, *error2.RECOMError)
}

type elasticSearchService struct {
	elasticSearchRepository repository.ElasticSearchRepository
	elErrorInterceptor      error2.ElasticSearchErrorInterceptor
}

func (esService elasticSearchService) GenerateSearchSuggestion(ctx context.Context, request models.SearchSuggestionRequest) (models.SuggestionData, *error2.RECOMError) {
	logger := log_utils.GetLogger(ctx)
	logger.Info("ElasticSearchService.GenerateSearchSuggestion.")
	resp, repoError := esService.elasticSearchRepository.SearchSuggestions(ctx, request)
	if repoError != nil {
		logger.Errorf("ElasticSearchService.GenerateSearchSuggestion: Error in searching elastic search documents. Error: %+v", repoError)
		return models.SuggestionData{}, esService.elErrorInterceptor.ErrorMapper(ctx, repoError)
	}
	logger.Info("ElasticSearchService.GenerateSearchSuggestion: document search successfully.")
	return resp, nil
}

func (esService elasticSearchService) Search(ctx context.Context, request models.SearchRequest) (interface{}, *error2.RECOMError) {
	logger := log_utils.GetLogger(ctx)
	logger.Info("ElasticSearchService.Search.")
	search, repoError := esService.elasticSearchRepository.Search(ctx, request.SearchText)
	if repoError != nil {
		logger.Errorf("ElasticSearchService.search: Error in searching elastic search documents. Error: %+v", repoError)
		return search, esService.elErrorInterceptor.ErrorMapper(ctx, repoError)
	}
	logger.Info("ElasticSearchService.search: document search successfully.")
	return search, nil
}

func NewElasticSearchService(searchRepository repository.ElasticSearchRepository, elErrorInterceptor error2.ElasticSearchErrorInterceptor) ElasticSearchService {
	return &elasticSearchService{
		elasticSearchRepository: searchRepository,
		elErrorInterceptor:      elErrorInterceptor,
	}
}
