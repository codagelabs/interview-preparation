package log_utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/spitertech/recommender/configurations"
	"os"
	"runtime"

	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

const RequestId = "trace_id"

func NewLogger(loggerConfig configurations.LoggerConfig) {
	logger = log.New()
	logger.SetLevel(log.DebugLevel)
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetOutput(os.Stdout)
}

func init() {
	logger = log.New()
	logger.SetLevel(log.DebugLevel)
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetOutput(os.Stdout)
}

func GetLogger(ctx context.Context) *log.Entry {
	var depth = 1
	var requestId string

	//Tracking Request Using Context
	if ctxRqID, ok := ctx.Value(RequestId).(string); ok {
		requestId = ctxRqID
	}
	function, _, line, _ := runtime.Caller(depth)
	functionObject := runtime.FuncForPC(function)
	entry := logger.WithFields(log.Fields{
		"request_id": requestId,
		//"file":       file,
		"function": functionObject.Name(),
		"line":     line,
	})
	logger.SetOutput(os.Stdout)
	return entry

}

// WithRqID returns a context with request ID or creates new a requestId and assigns to context
func WithRqID(ctx *gin.Context) context.Context {
	return context.WithValue(ctx.Request.Context(), RequestId, generateRequestID())
}

func generateRequestID() string {
	requestID := uuid.New()
	return requestID.String()

}
