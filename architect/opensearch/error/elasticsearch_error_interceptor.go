package error

import (
	"context"
	"errors"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
)

var ErrElasticDocNotFound = errors.New("no result found in elastic search")

type ElasticSearchErrorInterceptor interface {
	ErrorMapper(ctx context.Context, err error) *RECOMError
}

type elasticSearchErrorInterceptor struct {
}

func NewElasticSearchErrorInterceptor() ElasticSearchErrorInterceptor {
	return &elasticSearchErrorInterceptor{}
}

func (e elasticSearchErrorInterceptor) ErrorMapper(ctx context.Context, err error) *RECOMError {
	logger := logUtils.GetLogger(ctx)
	logger.Info("ElasticSearchService.ErrorMapper")

	switch err {
	case ErrElasticDocNotFound:
		logger.Infof("ElasticSearchService.ErrorMapper: Mapping elastic search record not found error to recom record not found error.")
		return ErrRecordNotFound
	default:
		logger.Infof("ElasticSearchService.ErrorMapper: Mapping elastic error to (default) internal server error.")
		return InternalServerErrorFunc(err.Error())
	}
}
