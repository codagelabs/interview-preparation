package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
	"gitlab.com/spitertech/recommender/utils/log_utils"
)

type ElasticSearchHelper interface {
	DocumentMarshaller(ctx context.Context, doc interface{}) string
	GenerateDocumentId(ctx context.Context) string
	ResponseParser(ctx context.Context, res *esapi.Response) error
}

type elasticSearchHelper struct {
}

func NewElasticSearchHelper() ElasticSearchHelper {
	return &elasticSearchHelper{}
}

func (esHelper elasticSearchHelper) GenerateDocumentId(ctx context.Context) string {
	log_utils.GetLogger(ctx).Info("ElasticSearchHelper.GenerateDocumentId: Generating documentId for indexing.")
	uuId := uuid.New()
	log_utils.GetLogger(ctx).Info("ElasticSearchHelper.GenerateDocumentId: UUID generated for doc")
	return uuId.String()
}

func (esHelper elasticSearchHelper) DocumentMarshaller(ctx context.Context, doc interface{}) string {
	b, err := json.Marshal(doc)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(b)
}

func (esHelper elasticSearchHelper) responseErrorDecoder(ctx context.Context, resp *esapi.Response) error {
	type errorResponse struct {
		Error struct {
			Reason    string `json:"reason"`
			RootCause []struct {
				Reason string `json:"reason"`
				Type   string `json:"type"`
			} `json:"root_cause"`
			Type string `json:"type"`
		} `json:"error"`
		Status int `json:"status"`
	}
	logger := log_utils.GetLogger(ctx)
	errorResp := errorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		logger.Errorf("ElasticSearchHelper.responseErrorDecoder: Error in parsing response body. Error: %+v", err)
		return err
	} else {
		logger.Errorf("ElasticSearchHelper.responseErrorDecoder: Error return by bulk indexer api. Error: %+v", errorResp)
		err = errors.New(fmt.Sprintf("  Error Status: [%d], ErrorType: %s, ErrorReason: %s", errorResp.Status, errorResp.Error.Type, errorResp.Error.Reason))
		return err
	}
}

func (esHelper elasticSearchHelper) responseBulkDecoder(ctx context.Context, resp *esapi.Response) error {
	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}
	var blk *bulkResponse
	logger := log_utils.GetLogger(ctx)
	if err := json.NewDecoder(resp.Body).Decode(&blk); err != nil {
		logger.Errorf("ElasticSearchHelper.ResponseExtractor: Error in parsing response body. Error: %s", err)
		return err
	} else {
		for _, d := range blk.Items {
			// ... so for any HTTP status above 201 ...
			//
			logger.Infof("%+v : ",d)
			if d.Index.Status > 201 {
				// ... increment the error counter ...
				//

				// ... and print the response status and error information ...
				errString := errors.New(fmt.Sprintf("  Error: [%d]: %s: %s: %s: %s",
					d.Index.Status,
					d.Index.Error.Type,
					d.Index.Error.Reason,
					d.Index.Error.Cause.Type,
					d.Index.Error.Cause.Reason,
				))
				logger.Errorf("ElasticSearchHelper.ResponseBulkResponseDecoder: Error in indexing doc. Error: %+v ", errString)

			} else {
				// ... otherwise increase the success counter.
				//
			}
		}
	}
	return nil
}

func (esHelper elasticSearchHelper) ResponseParser(ctx context.Context, res *esapi.Response) error {
	logger := log_utils.GetLogger(ctx)
	logger.Info("ElasticSearchHelper.ResponseParser: ")
	if res.IsError() {
		err := esHelper.responseErrorDecoder(ctx, res)
		if err != nil {
			logger.Errorf("ElasticSearchHelper.ResponseParser: Error in decoding elastic search error response. Error: %+v ", err)
			return err
		}
	} else {
		err := esHelper.responseBulkDecoder(ctx, res)
		if err != nil {
			logger.Errorf("ElasticSearchHelper.ResponseParser: Error in decoding elastic search response. Error: %-v ", err)
			return err
		}

	}
	// Close the response body, to prevent reaching the limit for goroutines or file handles
	//
	res.Body.Close()
	return nil
}
