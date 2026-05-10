package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/spitertech/recommender/configurations"
	recomError "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/repository"
	"gitlab.com/spitertech/recommender/utils/log_utils"
)

type ProductBulkIndexerService interface {
	ProductBulkIndexer(ctx context.Context, products []models.Product,indexName string) (interface{}, *recomError.RECOMError)
}
type productBulkIndexerService struct {
	bulkIndexer                   repository.BulkDocumentIndexer
	bulkIndexerConfig             configurations.BulkIndexerConfig
	elasticSearchErrorInterceptor recomError.ElasticSearchErrorInterceptor
}

func (b productBulkIndexerService) getBatchCount(totalProducts int) int {

	if totalProducts%b.bulkIndexerConfig.BatchSize == 0 {
		return totalProducts / b.bulkIndexerConfig.BatchSize
	} else {
		return (totalProducts / b.bulkIndexerConfig.BatchSize) + 1
	}
}
func (b productBulkIndexerService) ProductBulkIndexer(ctx context.Context, products []models.Product, indexName string) (interface{}, *recomError.RECOMError) {
	logger := log_utils.GetLogger(ctx)
	logger.Info("ProductBulkIndexerService.ProductBulkIndexer: Creating batches from data for indexing in elastic search")
	totalProducts := len(products)
	numBatches := b.getBatchCount(totalProducts)
	currBatch := 0
	var buf bytes.Buffer
	for i, product := range products {
		currBatch = i / b.bulkIndexerConfig.BatchSize
		if i == totalProducts-1 {
			currBatch++
		}

		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, product.ProductID, "\n"))

		var data, err = json.Marshal(product)
		if err != nil {
			logger.Errorf("ProductBulkIndexerService.ProductBulkIndexer: Error in encode product %s. Error %+v :", product.ProductID, err)
			return nil, b.elasticSearchErrorInterceptor.ErrorMapper(ctx, err)
		}

		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		if i > 0 && i%b.bulkIndexerConfig.BatchSize == 0 || i == totalProducts-1 {
			logger.Infof("ProductBulkIndexerService.ProductBulkIndexer: Batch processing started Current_Batch/Total_Batch  %v/%v :", currBatch, numBatches)
			b.bulkIndexer.BulkIndices(context.Background(), buf, indexName)
		}
	}
	return "success", nil
}

func NewProductBulkIndexer(bulkIndexerRepo repository.BulkDocumentIndexer, bulkIndexerConfig configurations.BulkIndexerConfig, elasticSearchErrorInterceptor recomError.ElasticSearchErrorInterceptor) ProductBulkIndexerService {
	return productBulkIndexerService{bulkIndexer: bulkIndexerRepo, bulkIndexerConfig: bulkIndexerConfig, elasticSearchErrorInterceptor: elasticSearchErrorInterceptor}

}
