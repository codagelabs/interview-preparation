package repository

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch"
	"gitlab.com/spitertech/recommender/repository/helper"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
)

type BulkDocumentIndexer interface {
	BulkIndices(ctc context.Context, buf bytes.Buffer, indexName string) error
}

type bulkDocumentIndexer struct {
	esClient            *elasticsearch.Client
	elasticSearchHelper helper.ElasticSearchHelper
}

func (b bulkDocumentIndexer) BulkIndices(ctx context.Context, buf bytes.Buffer, indexName string) error {
	logger := logUtils.GetLogger(ctx)
	logger.Infof("BulkDocumentIndexer.BulkIndices: Document indexing in bulk for index name %s ", indexName)
	res, err := b.esClient.Bulk(bytes.NewReader(buf.Bytes()), b.esClient.Bulk.WithIndex(indexName))
	if err != nil {
		logger.Errorf("BulkDocumentIndexer.BulkIndices: Error in bulk indexing api. Error: %+v", err)
		return err
	}

	err = b.elasticSearchHelper.ResponseParser(ctx, res)
	if err != nil {
		logger.Errorf("BulkDocumentIndexer.BulkIndices: Error in parsing bulk indexing api response . Error: %+v", err)
		return err
	}
	buf.Reset()
	return nil
}

func NewBulkDocumentIndexer(esClient *elasticsearch.Client, elasticSearchHelper helper.ElasticSearchHelper) BulkDocumentIndexer {
	return &bulkDocumentIndexer{esClient: esClient, elasticSearchHelper: elasticSearchHelper}
}
