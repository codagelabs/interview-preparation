package es_query_excecutor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	recomError "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
	"io"
	"reflect"
)

type EsExecutor interface {
	WithContext() EsExecutor
	WithIndex(index string) EsExecutor
	WithBody(v io.Reader) EsExecutor
	WithTrackTotalHits(true bool) EsExecutor
	WithResponseAs(respType *map[string]interface{}) EsExecutor
	WithPretty() EsExecutor
	ResponseAs(respType interface{}) EsExecutor
	ExecuteSearchQuery() error
	ExecuteSuggestQuery() error
}

type esExecutor struct {
	searchReq []func(*esapi.SearchRequest)
	esClient  *elasticsearch.Client
	ctx       context.Context
	resp      *map[string]interface{}
	Response  interface{}
	err       error
	apiResp   *esapi.Response
}

func (executor esExecutor) Error() error {
	if executor.err != nil {
		return executor.err
	}
	return nil
}

func (executor esExecutor) WithContext() EsExecutor {
	executor.searchReq = append(executor.searchReq, executor.esClient.Search.WithContext(executor.ctx))
	return executor
}

func (executor esExecutor) WithIndex(index string) EsExecutor {
	executor.searchReq = append(executor.searchReq, executor.esClient.Search.WithIndex(index))
	return executor
}

func (executor esExecutor) WithBody(v io.Reader) EsExecutor {
	executor.searchReq = append(executor.searchReq, executor.esClient.Search.WithBody(v))
	return executor
}

func (executor esExecutor) WithTrackTotalHits(value bool) EsExecutor {
	executor.searchReq = append(executor.searchReq, executor.esClient.Search.WithTrackTotalHits(value))
	return executor
}

func (executor esExecutor) WithPretty() EsExecutor {
	executor.searchReq = append(executor.searchReq, executor.esClient.Search.WithPretty())
	return executor
}

func (executor esExecutor) ExecuteSearchQuery() (err error) {

	logger := logUtils.GetLogger(executor.ctx)
	executor.apiResp, err = executor.esClient.Search(executor.searchReq...)

	if executor.apiResp.IsError() {
		var errorMap map[string]interface{}
		err = json.NewDecoder(executor.apiResp.Body).Decode(&errorMap)
		if err != nil {
			m := errorMap["error"].(map[string]interface{})
			err := fmt.Sprintf("ElasticSearchRepository.Search:  Api response : [%s]. Error Type: %s. Error: %s.",
				executor.apiResp.Status(),
				m["type"],
				m["reason"],
			)
			return errors.New(err)
		}
	}
	var resp = map[string]interface{}{}

	err = json.NewDecoder(executor.apiResp.Body).Decode(&resp)
	if err != nil {
		return err
	}

	fmt.Printf("Api Respponse :  %+v", resp)

	var r = map[string]interface{}{}

	for _, hit := range resp["hits"].(map[string]interface{})["hits"].([]interface{}) {
		logger.Infof(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		r[hit.(map[string]interface{})["_id"].(string)] = hit.(map[string]interface{})
	}
	if len(r) == 0 {
		return recomError.ErrElasticDocNotFound
	}

	*executor.resp = r
	fmt.Println(executor.resp)
	return nil
}

func (executor esExecutor) ExecuteSuggestQuery() (err error) {

	type SuggesterResponse struct {
		Shards struct {
			Failed     int `json:"failed"`
			Skipped    int `json:"skipped"`
			Successful int `json:"successful"`
			Total      int `json:"total"`
		} `json:"shards"`
		Hits struct {
			TimedOut bool `json:"timed_out"`
			Took     int  `json:"took"`
			Hits     []struct {
				ID     string  `json:"_id"`
				Result string  `json:"_index"`
				Score  float64 `json:"_score"`
				Type   string  `json:"_type"`
				Source struct {
					models.Product
					Suggest struct {
						Input  []string `json:"input"`
						Weight int      `json:"weight"`
					}
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	//logger := logUtils.GetLogger(executor.ctx)
	executor.apiResp, err = executor.esClient.Search(executor.searchReq...)

	if executor.apiResp.IsError() {
		var errorMap map[string]interface{}
		err = json.NewDecoder(executor.apiResp.Body).Decode(&errorMap)
		if err != nil {
			m := errorMap["error"].(map[string]interface{})
			err := fmt.Sprintf("ElasticSearchRepository.Search:  Api response : [%s]. Error Type: %s. Error: %s.",
				executor.apiResp.Status(),
				m["type"],
				m["reason"],
			)
			return errors.New(err)
		}
	}
	//var resp = map[string]interface{}{}
	var suggesterResponse = SuggesterResponse{}

	err = json.NewDecoder(executor.apiResp.Body).Decode(&suggesterResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Api Respponse Data :  %+v", suggesterResponse)

	var resp = []models.Product{}

	for _, hit := range suggesterResponse.Hits.Hits {
		resp = append(resp, hit.Source.Product)
	}

	if len(resp) == 0 {
		return recomError.ErrElasticDocNotFound
	}
	if executor.Response != nil {
		v := reflect.ValueOf(executor.Response).Elem()
		v.Set(reflect.ValueOf(resp))

	}
	return nil
}
func (executor esExecutor) ResponseAs(resp interface{}) EsExecutor {
	executor.Response = resp
	return executor
}
func (executor esExecutor) WithResponseAs(respAs *map[string]interface{}) EsExecutor {
	executor.resp = respAs
	return executor
}

func NewEsExecutor(ctx context.Context, esClient *elasticsearch.Client) EsExecutor {
	return newEsExecutor(ctx, esClient)
}

func newEsExecutor(ctx context.Context, client *elasticsearch.Client) EsExecutor {
	return &esExecutor{
		searchReq: make([]func(request *esapi.SearchRequest), 0),
		esClient:  client,
		ctx:       ctx,
	}
}
