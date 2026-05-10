package error

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type ElasticSearchErrorInterceptorTestSuite struct {
	suite.Suite
	mockCtrl                      *gomock.Controller
	recorder                      *httptest.ResponseRecorder
	context                       context.Context
	elasticSearchErrorInterceptor ElasticSearchErrorInterceptor
}

func TestElasticSearchErrorInterceptorTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchErrorInterceptorTestSuite))
}

func (suite *ElasticSearchErrorInterceptorTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context =context.Background()
	suite.elasticSearchErrorInterceptor = NewElasticSearchErrorInterceptor()
}

func (suite ElasticSearchErrorInterceptorTestSuite) TearDownTest() {
	//suite.mockCtrl.Finish()
}

func (suite ElasticSearchErrorInterceptorTestSuite) TestErrorMapper_ShouldReturnInterNalServerError() {
	errMassage := "invalid request"
	actualError := suite.elasticSearchErrorInterceptor.ErrorMapper(suite.context, errors.New(errMassage))
	suite.Equal(InternalServerErrorFunc(errMassage),actualError)
}

func (suite ElasticSearchErrorInterceptorTestSuite) TestErrorMapper_ShouldReturnRecordNotFound() {
	actualError := suite.elasticSearchErrorInterceptor.ErrorMapper(suite.context, ErrElasticDocNotFound)
	suite.Equal(ErrRecordNotFound,actualError)

}
