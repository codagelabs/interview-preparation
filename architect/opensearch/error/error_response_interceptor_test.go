package error

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrResponseInterceptorSuite struct {
	suite.Suite
	mockCtrl               *gomock.Controller
	recorder               *httptest.ResponseRecorder
	context                *gin.Context
	errResponseInterceptor ErrResponseInterceptor
}

func TestErrResponseInterceptorSuite(t *testing.T) {
	suite.Run(t, new(ErrResponseInterceptorSuite))
}

func (suite *ErrResponseInterceptorSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context,_= gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "", nil)
	suite.errResponseInterceptor = NewErrResponseInterceptor()
}

func (suite ErrResponseInterceptorSuite) TearDownTest() {
	//suite.mockCtrl.Finish()
}

func (suite ErrResponseInterceptorSuite) TestHandleBadRequest_ShouldReturnBadRequestError() {
	errMassage := "invalid request"
	suite.errResponseInterceptor.HandleBadRequest(suite.context, errors.New(errMassage))
	errResponse := ErrResponse{}
	bytesBody, _ := ioutil.ReadAll(suite.recorder.Body)
	json.Unmarshal(bytesBody, &errResponse)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(ErrResponse{ErrCode: BadRequestErrorCode, ErrMessage: "invalid request"}, errResponse)
}


func (suite ErrResponseInterceptorSuite) TestHandleServiceError_ShouldReturnInternalServerError() {
	errMassage := "invalid request"
	expectedErr:=InternalServerErrorFunc(errMassage)
	suite.errResponseInterceptor.HandleServiceError(suite.context,expectedErr)
	errResponse := ErrResponse{}
	bytesBody, _ := ioutil.ReadAll(suite.recorder.Body)
	json.Unmarshal(bytesBody, &errResponse)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(ErrResponse{ErrCode: InternalServerErrorCode, ErrMessage: errMassage}, errResponse)
}