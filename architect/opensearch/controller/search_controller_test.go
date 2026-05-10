package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/mocks"
	"gitlab.com/spitertech/recommender/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ElasticSearchControllerTestSuite struct {
	suite.Suite
	mockCtrl                *gomock.Controller
	recorder                *httptest.ResponseRecorder
	context                 *gin.Context
	elasticsearchService    *mocks.MockElasticSearchService
	errResponseInterceptor  *mocks.MockErrResponseInterceptor
	elasticSearchController SearchController
}

func TestElasticSearchControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchControllerTestSuite))
}

func (suite *ElasticSearchControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.errResponseInterceptor = mocks.NewMockErrResponseInterceptor(suite.mockCtrl)
	suite.elasticsearchService = mocks.NewMockElasticSearchService(suite.mockCtrl)
	suite.elasticSearchController = NewElasticSearchController(suite.errResponseInterceptor, suite.elasticsearchService)
}

func (suite *ElasticSearchControllerTestSuite) TearDownTest() {
	//suite.mockCtrl.Finish()
}

func (suite ElasticSearchControllerTestSuite) TestSearch_ShouldReturnBadRequestErrorWhenRequestBodyInvalid() {
	suite.context.Request, _ = http.NewRequest("GET", "/search", nil)
	suite.errResponseInterceptor.EXPECT().HandleBadRequest(suite.context, gomock.Any())
	suite.elasticSearchController.Search(suite.context)
}

func (suite ElasticSearchControllerTestSuite) TestSearch_ShouldReturnSameErrorIfElasricSearchServiceReturnsAnError() {
	searchText := "search_text"
	var searchRequest = models.SearchRequest{
		SearchText: searchText,
	}
	suite.context.Request, _ = http.NewRequest("GET", "search?search_text="+searchText, nil)
	error := error2.InternalServerErrorFunc("something-went-wrong")
	suite.elasticsearchService.EXPECT().Search(suite.context, searchRequest).Return(nil, error)
	suite.errResponseInterceptor.EXPECT().HandleServiceError(suite.context, error)
	suite.elasticSearchController.Search(suite.context)
}

func (suite ElasticSearchControllerTestSuite) TestSearch_ShouldReturnSearchResultSuccessfully() {
	searchText := "search_text"
	var searchRequest = models.SearchRequest{
		SearchText: searchText,
	}
	resp := map[string]interface{}{
		"test": "test1",
	}
	suite.context.Request, _ = http.NewRequest("GET", "search?search_text="+searchText, nil)
	suite.elasticsearchService.EXPECT().Search(suite.context, searchRequest).Return(resp, nil)
	suite.elasticSearchController.Search(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite ElasticSearchControllerTestSuite) TestRecommendSearch_ShouldReturnBadRequestWhenMadndatoryFeildMissing() {
	suite.context.Request, _ = http.NewRequest("GET", "/search", nil)
	suite.errResponseInterceptor.EXPECT().HandleBadRequest(suite.context, gomock.Any())
	suite.elasticSearchController.RecommendSearch(suite.context)
}

func (suite ElasticSearchControllerTestSuite) TestRecommendSearch_ShouldReturnSameErrorIfElasricSearchServiceReturnsAnError() {
	searchText := "search_text"
	var searchRequest = models.SearchSuggestionRequest{
		SearchText: searchText,
	}
	suite.context.Request, _ = http.NewRequest("GET", "search?search_text="+searchText, nil)
	error := error2.InternalServerErrorFunc("something-went-wrong")
	suite.elasticsearchService.EXPECT().GenerateSearchSuggestion(suite.context, searchRequest).Return(models.SuggestionData{}, error)
	suite.errResponseInterceptor.EXPECT().HandleServiceError(suite.context, error)
	suite.elasticSearchController.RecommendSearch(suite.context)
}

func (suite ElasticSearchControllerTestSuite) TestRecommendSearch_ShouldReturnSearchResultSuccessfully() {
	searchText := "search_text"
	var searchRequest = models.SearchSuggestionRequest{
		SearchText: searchText,
	}
	resp := models.SuggestionData{
		Product: []models.Product{},
	}
	suite.context.Request, _ = http.NewRequest("GET", "search?search_text="+searchText, nil)
	suite.elasticsearchService.EXPECT().GenerateSearchSuggestion(suite.context, searchRequest).Return(resp, nil)
	suite.elasticSearchController.RecommendSearch(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func getProductRequest(value interface{}) io.Reader {
	byteData, _ := json.Marshal(value)
	return bytes.NewReader(byteData)
}
