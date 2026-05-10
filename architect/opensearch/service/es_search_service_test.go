package service

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/mocks"
	"gitlab.com/spitertech/recommender/models"
	"net/http/httptest"
	"testing"
)

type ElasticSearchServiceTestSuite struct {
	suite.Suite
	recorder                  *httptest.ResponseRecorder
	mockCtrl                  *gomock.Controller
	context                   context.Context
	elasticSearchRepository   *mocks.MockElasticSearchRepository
	elasticSearchErrorHandler *mocks.MockElasticSearchErrorInterceptor
	elasticSearchService      ElasticSearchService
}

func TestElasticSearchServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchServiceTestSuite))
}

func (suite *ElasticSearchServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context = context.Background()
	suite.elasticSearchRepository = mocks.NewMockElasticSearchRepository(suite.mockCtrl)
	suite.elasticSearchErrorHandler = mocks.NewMockElasticSearchErrorInterceptor(suite.mockCtrl)
	suite.elasticSearchService = NewElasticSearchService(suite.elasticSearchRepository, suite.elasticSearchErrorHandler)

}

func (suite *ElasticSearchServiceTestSuite) TearDownTest() {
	//suite.mockCtrl.Finish()
}

func (suite ElasticSearchServiceTestSuite) TestSearch_ShouldReturnAnErrorIfRepositoryReturnsAnError() {
	var searchRequest = models.SearchRequest{
		SearchText: "search-text",
	}
	err := errors.New("something went wrong")

	suite.elasticSearchRepository.EXPECT().Search(suite.context, searchRequest.SearchText).Return(nil, err)
	suite.elasticSearchErrorHandler.EXPECT().ErrorMapper(suite.context, err).Return(error2.InternalServerErrorFunc("something-went-wrong"))

	indexedResonse, actualError := suite.elasticSearchService.Search(suite.context, searchRequest)
	suite.Equal(error2.InternalServerErrorFunc("something-went-wrong"), actualError)
	suite.Nil(indexedResonse)
}

func (suite ElasticSearchServiceTestSuite) TestSearch_ShouldReturnAnSuccessResponse() {
	var searchRequest = models.SearchRequest{
		SearchText: "search-text",
	}
	resp := map[string]interface{}{}
	suite.elasticSearchRepository.EXPECT().Search(suite.context, searchRequest.SearchText).Return(resp, nil)
	actualResponse, actualError := suite.elasticSearchService.Search(suite.context, searchRequest)
	suite.Equal(actualResponse, resp)
	suite.Nil(actualError)
}

func (suite ElasticSearchServiceTestSuite) TestGenerateSearchSuggestion_ShouldReturnAnErrorIfRepositoryReturnsAnError() {
	var searchRequest = models.SearchSuggestionRequest{
		SearchText: "mobile",
	}
	suite.elasticSearchRepository.EXPECT().SearchSuggestions(suite.context, searchRequest).Return(models.SuggestionData{}, errors.New("something-went-wrong"))
	suite.elasticSearchErrorHandler.EXPECT().ErrorMapper(suite.context, errors.New("something-went-wrong")).Return(error2.InternalServerErrorFunc("something-went-wrong"))
	actualResponse, actualError := suite.elasticSearchService.GenerateSearchSuggestion(suite.context, searchRequest)
	suite.Equal(error2.InternalServerErrorFunc("something-went-wrong"), actualError)
	suite.Empty(models.SuggestionData{}, actualResponse)
}

func (suite ElasticSearchServiceTestSuite) TestGenerateSearchSuggestion_ShouldReturnAnSuccessResponse() {
	var searchRequest = models.SearchSuggestionRequest{
		SearchText: "mobile",
	}
	resp := models.SuggestionData{
		Product: []models.Product{
			models.Product{ProductID: "p1", IsAvailable: true, ProductCategory: "electronics", ProductName: "mobile", ProductPrise: 20, Suggest: models.Suggest{Input: []string{"mobile", "mobile phone"}}},
		},
	}
	suite.elasticSearchRepository.EXPECT().SearchSuggestions(suite.context, searchRequest).Return(resp, nil)
	actualResponse, actualError := suite.elasticSearchService.GenerateSearchSuggestion(suite.context, searchRequest)
	suite.Equal(actualResponse, resp)
	suite.Nil(actualError)
}
