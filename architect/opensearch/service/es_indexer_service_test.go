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

type EsIndexerServiceTestSuite struct {
	suite.Suite
	recorder                    *httptest.ResponseRecorder
	mockCtrl                    *gomock.Controller
	context                     context.Context
	esDocumentIndexerRepository *mocks.MockEsDocumentIndexerRepository
	elasticSearchErrorHandler   *mocks.MockElasticSearchErrorInterceptor
	esIndexerService        EsIndexerService
}

func TestEsIndexerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EsIndexerServiceTestSuite))
}

func (suite *EsIndexerServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context = context.Background()
	suite.esDocumentIndexerRepository = mocks.NewMockEsDocumentIndexerRepository(suite.mockCtrl)
	suite.elasticSearchErrorHandler = mocks.NewMockElasticSearchErrorInterceptor(suite.mockCtrl)
	suite.esIndexerService = NewEsIndexerService(suite.esDocumentIndexerRepository, suite.elasticSearchErrorHandler)

}

func (suite EsIndexerServiceTestSuite) TestIndexDocument_ShouldReturnAnErrorIfRepositoryReturnsAnError() {
	product := models.Product{
		ProductID:       "1",
		ProductName:     "test-product",
		ProductCategory: "test-category",
		ProductPrise:    0,
		IsAvailable:     false,
	}

	err := errors.New("something went wrong")

	suite.esDocumentIndexerRepository.EXPECT().IndexDocument(suite.context, product).Return(nil, err)
	suite.elasticSearchErrorHandler.EXPECT().ErrorMapper(suite.context, err).Return(error2.InternalServerErrorFunc("something-went-wrong"))
	indexedResonse, actualError := suite.esIndexerService.IndexDocuments(suite.context, product)
	suite.Equal(error2.InternalServerErrorFunc("something-went-wrong"), actualError)
	suite.Nil(indexedResonse)
}

func (suite EsIndexerServiceTestSuite) TestIndexDocument_ShouldReturnAnSuccessResponse() {
	product := models.Product{
		ProductID:       "1",
		ProductName:     "test-product",
		ProductCategory: "test-category",
		ProductPrise:    0,
		IsAvailable:     false,
	}

	resp:=map[string]interface{}{}

	suite.esDocumentIndexerRepository.EXPECT().IndexDocument(suite.context, product).Return(resp, nil)
	actualResponse, actualError := suite.esIndexerService.IndexDocuments(suite.context, product)
	suite.Equal(actualResponse,resp)
	suite.Nil(actualError)
}