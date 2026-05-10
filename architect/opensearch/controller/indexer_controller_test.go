package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/mocks"
	"gitlab.com/spitertech/recommender/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type IndexerControllerTestSuite struct {
	suite.Suite
	mockCtrl                 *gomock.Controller
	recorder                 *httptest.ResponseRecorder
	context                  *gin.Context
	esDocumentindexerService *mocks.MockEsIndexerService
	errResponseInterceptor   *mocks.MockErrResponseInterceptor
	elasticSearchController  IndexerController
}

func TestIndexerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IndexerControllerTestSuite))
}

func (suite *IndexerControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.errResponseInterceptor = mocks.NewMockErrResponseInterceptor(suite.mockCtrl)
	suite.esDocumentindexerService = mocks.NewMockEsIndexerService(suite.mockCtrl)
	suite.elasticSearchController = NewIndexerController(suite.errResponseInterceptor, suite.esDocumentindexerService)
}

func (suite *IndexerControllerTestSuite) TearDownTest() {
	//suite.mockCtrl.Finish()
}

func (suite IndexerControllerTestSuite) TestIndexDocument_ShouldReturnBadRequestErrorWhenRequestBodyInvalid() {
	suite.context.Request, _ = http.NewRequest("POST", "/index-doc", nil)
	suite.errResponseInterceptor.EXPECT().HandleBadRequest(suite.context, gomock.Any())
	suite.elasticSearchController.IndexDocument(suite.context)
}

func (suite IndexerControllerTestSuite) TestIndexDocument_ShouldReturnSameErrorWhenServiceReturnsError() {
	errorMessage:="something went wrong"
	request := models.Product{ProductID: "product-id", IsAvailable: true, ProductPrise: 1, ProductName: "product-name", ProductCategory: "category"}
	suite.context.Request, _ = http.NewRequest("POST", "/index-doc",getProductRequest(request))
	suite.esDocumentindexerService.EXPECT().IndexDocuments(suite.context,request ).Return(nil,error2.InternalServerErrorFunc(errorMessage))
	suite.errResponseInterceptor.EXPECT().HandleServiceError(suite.context, gomock.Any())
	suite.elasticSearchController.IndexDocument(suite.context)
}

func (suite IndexerControllerTestSuite) TestIndexDocument_ShouldReturnResponseSuccessfully() {
	request := models.Product{ProductID: "product-id", IsAvailable: true, ProductPrise: 1, ProductName: "product-name", ProductCategory: "category"}
	suite.context.Request, _ = http.NewRequest("POST", "/index-doc",getProductRequest(request))
	suite.esDocumentindexerService.EXPECT().IndexDocuments(suite.context,request ).Return(nil,nil)
	suite.elasticSearchController.IndexDocument(suite.context)
	suite.Equal(http.StatusOK,suite.recorder.Code)
}
