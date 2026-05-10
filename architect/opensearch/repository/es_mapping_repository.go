package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"strings"
)

type ESMappingHandler interface {
	CreateMapping(ctx context.Context) error
}
type esMappingHandler struct {
	esClient *elasticsearch.Client
}

func NewEsMappingHandler(esClient *elasticsearch.Client) ESMappingHandler {
	return &esMappingHandler{esClient: esClient}
}

func (e esMappingHandler) CreateMapping(ctx context.Context) error {
	req := esapi.IndicesPutMappingRequest{
		Index: []string{"product"},
		Body: strings.NewReader(`{
		"mappings": {
			"_doc" : {
				"properties" : {
					"product_name" : {
						"type": "search_as_you_type",
					},
					"product_category" : {
						"type": "search_as_you_type"
					}
					"product_prise" : {
						"type": "number"
					}
					"is_available" :  {
						"type": "boolean"
					}

				}
			}
    	}
	}`),
	}

	res, err := req.Do(ctx, e.esClient)
	if err != nil {
		fmt.Errorf("ElasticSearchRepository.IndexDocument: Error in IndexRequest api. Error: %+v ", err)
		return err
	}
	defer res.Body.Close()
	resMap := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		fmt.Errorf("ElasticSearchRepository.IndexDocument: Error in decoding response body. Error: %+v ", err)
		return err
	}
	fmt.Println("Mapping created successfully")
	return nil
}
