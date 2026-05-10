package es_query_excecutor

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EsExecutorTestSuite struct {
	suite.Suite
	recorder   *httptest.ResponseRecorder
	mockCtrl   *gomock.Controller
	context    *gin.Context
	esExecutor EsExecutor
}

func TestEsExecutorTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(EsExecutorTestSuite))
}

func (suite *EsExecutorTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("POST", "/index-doc", nil)

	//suite.esExecutor = NewEsExecutor(suite.context,)

}

func (suite *EsExecutorTestSuite) TearDownTest() {
}
