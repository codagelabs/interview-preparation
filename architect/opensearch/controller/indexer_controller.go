package controller

import (
	"github.com/gin-gonic/gin"
	recomError "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/service"
	logUtils "gitlab.com/spitertech/recommender/utils/log_utils"
	"net/http"
)

type IndexerController interface {
	IndexDocument(ctx *gin.Context)
}
type indexerController struct {
	errResponseInterceptor recomError.ErrResponseInterceptor
	indexerService         service.EsIndexerService
}

func (elSearch indexerController) IndexDocument(ctx *gin.Context) {
	logger := logUtils.GetLogger(ctx)
	logger.Info("IndexerController.IndexDocument: Document Indexing started")
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		logger.Errorf("IndexerController.IndexDocument: Error in binding request body. Error: %+v ", err)
		elSearch.errResponseInterceptor.HandleBadRequest(ctx, err)
		return
	}
	documents, serviceError := elSearch.indexerService.IndexDocuments(ctx, product)
	if serviceError != nil {
		logger.Errorf("IndexerController.IndexDocument: Error from indexer service. Error: %+v", serviceError)
		elSearch.errResponseInterceptor.HandleServiceError(ctx, serviceError)
		return
	}
	logger.Info("IndexerController.IndexDocument: Document Indexed successfully.")
	ctx.JSON(http.StatusOK, documents)
}

func NewIndexerController(errInterceptor recomError.ErrResponseInterceptor, indexerService service.EsIndexerService) IndexerController {
	return &indexerController{
		errResponseInterceptor: errInterceptor,
		indexerService:         indexerService,
	}
}
