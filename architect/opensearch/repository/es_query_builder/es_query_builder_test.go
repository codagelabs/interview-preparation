package es_query_builder

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ElasticSearchQueryBuilderTestSuite struct {
	suite.Suite
	recorder                  *httptest.ResponseRecorder
	mockCtrl                  *gomock.Controller
	context                   *gin.Context
	elasticSearchQueryBuilder ElasticSearchQueryBuilder
}

func TestElasticSearchQueryBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchQueryBuilderTestSuite))
}

func (suite *ElasticSearchQueryBuilderTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("POST", "/index-doc", nil)
	suite.elasticSearchQueryBuilder = NewElasticSearchQueryBuilder()

}

func (suite *ElasticSearchQueryBuilderTestSuite) TearDownTest() {
}

func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithSize() {
	query := map[string]interface{}{
		"size":  20,
		"query": map[string]interface{}{},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	NewElasticSearchQueryBuilder().WithSize(20).BindQuery(&actualBuffer)

	suite.Equal(expectedBuffer, actualBuffer)
}

func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithMatch() {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"product_name": "test",
			},
		},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	NewElasticSearchQueryBuilder().WithMatch("test").BindQuery(&actualBuffer)

	suite.Equal(expectedBuffer, actualBuffer)
}

func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithMatchAll() {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	NewElasticSearchQueryBuilder().WithMatchAll().BindQuery(&actualBuffer)

	suite.Equal(expectedBuffer, actualBuffer)
}

func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithMultiMatchAll() {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  "test",
				"fields": []string{"test1", "test2"},
			},
		},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	NewElasticSearchQueryBuilder().WithMultiMatch("test", []string{"test1", "test2"}).BindQuery(&actualBuffer)

	suite.Equal(expectedBuffer, actualBuffer)
}


func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithSuggester() {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  "test",
				"type":   "bool_prefix",
				"fields": []string{"test2", "test2.3gram","test2.3gram"},
			},
		},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	NewElasticSearchQueryBuilder().WithSuggester("test", []string{"test2", "test2.3gram","test2.3gram"}).BindQuery(&actualBuffer)

	suite.Equal(expectedBuffer, actualBuffer)
}

func (suite ElasticSearchQueryBuilderTestSuite) TestWithSize_ShouldReturnQueryWithError() {
	query := map[string]interface{}{
		"size":  20,
		"query": map[string]interface{}{},
	}
	var expectedBuffer bytes.Buffer
	json.NewEncoder(&expectedBuffer).Encode(query)
	var actualBuffer bytes.Buffer
	err := NewElasticSearchQueryBuilder().WithSize(20).BindQuery(&actualBuffer).Error()
	suite.Equal(expectedBuffer, actualBuffer)
	suite.Nil(err)
}
