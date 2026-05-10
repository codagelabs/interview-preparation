package service

import (
	"context"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/repository"
	"gitlab.com/spitertech/recommender/utils/log_utils"
)

type EsIndexerService interface {
	IndexDocuments(ctx context.Context, request models.Product) (resp interface{}, recomError *error2.RECOMError)
}

type esIndexerService struct {
	esDocumentIndexer  repository.EsDocumentIndexerRepository
	elErrorInterceptor error2.ElasticSearchErrorInterceptor
}

func (esService esIndexerService) IndexDocuments(ctx context.Context, product models.Product) (resp interface{}, recomError *error2.RECOMError) {
	logger := log_utils.GetLogger(ctx)
	logger.Info("EsIndexerService.IndexDocuments.")
	resp, repoError := esService.esDocumentIndexer.IndexDocument(ctx, product)
	if repoError != nil {
		logger.Errorf("EsIndexerService.IndexDocuments: Error in adding elastic search documents. Error: %+v", repoError)
		return nil, esService.elErrorInterceptor.ErrorMapper(ctx, repoError)
	}
	logger.Info("EsIndexerService.IndexDocuments: returning response.")
	return resp, nil
}

func NewEsIndexerService(esDocumentIndexer  repository.EsDocumentIndexerRepository, elErrorInterceptor error2.ElasticSearchErrorInterceptor) EsIndexerService {
	return &esIndexerService{
		esDocumentIndexer: esDocumentIndexer,
		elErrorInterceptor:      elErrorInterceptor,
	}
}
