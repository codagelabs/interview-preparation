package controller

import (
	"github.com/gin-gonic/gin"
	recomError "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/service"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
	"net/http"
)

type SearchController interface {
	Search(ctx *gin.Context)
	RecommendSearch(ctx *gin.Context)
}
type elasticSearchController struct {
	errResponseInterceptor recomError.ErrResponseInterceptor
	elasticSearchService   service.ElasticSearchService
}

func (elSearch elasticSearchController) Search(ctx *gin.Context) {
	logger := logUtils.GetLogger(ctx)
	logger.Info("ElasticSearchController.Search: Search request started")
	var searchRequest models.SearchRequest
	if err := ctx.ShouldBindQuery(&searchRequest); err != nil {
		logger.Errorf("ElasticSearchController.Search: Error in validating and binding request body. Error: %-v ", err)
		elSearch.errResponseInterceptor.HandleBadRequest(ctx, err)
		return
	}
	search, serviceError := elSearch.elasticSearchService.Search(ctx, searchRequest)
	if serviceError != nil {
		logger.Errorf("ElasticSearchController.Search: Error in elastic search service. Error: %+v", serviceError)
		elSearch.errResponseInterceptor.HandleServiceError(ctx, serviceError)
		return
	}
	logger.Info("ElasticSearchController.Search: Search Result return successfully")
	ctx.JSON(http.StatusOK, search)
}

func (elSearch elasticSearchController) RecommendSearch(ctx *gin.Context) {
	logger := logUtils.GetLogger(ctx)
	logger.Info("ElasticSearchController.RecommendSearch: Search search recommendation started")
	var searchRequest models.SearchSuggestionRequest
	if err := ctx.ShouldBindQuery(&searchRequest); err != nil {
		logger.Errorf("ElasticSearchController.RecommendSearch: Error in binding request search request body. Error: %+v ", err)
		elSearch.errResponseInterceptor.HandleBadRequest(ctx, err)
		return
	}
	search, serviceError := elSearch.elasticSearchService.GenerateSearchSuggestion(ctx, searchRequest)
	if serviceError != nil {
		logger.Errorf("ElasticSearchController.RecommendSearch: Error in elastic search service. Error: %-v", serviceError)
		elSearch.errResponseInterceptor.HandleServiceError(ctx, serviceError)
		return
	}
	logger.Info("ElasticSearchController.RecommendSearch: Search Result return successfully")
	ctx.JSON(http.StatusOK, search)
}

func NewElasticSearchController(errInterceptor recomError.ErrResponseInterceptor, elasticSearchService service.ElasticSearchService) SearchController {
	return &elasticSearchController{
		errResponseInterceptor: errInterceptor,
		elasticSearchService:   elasticSearchService,
	}
}
