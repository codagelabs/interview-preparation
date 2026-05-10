package router

import (
	"context"
	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"gitlab.com/spitertech/recommender/configurations"
	"gitlab.com/spitertech/recommender/controller"
	error2 "gitlab.com/spitertech/recommender/error"
	"gitlab.com/spitertech/recommender/models"
	"gitlab.com/spitertech/recommender/repository"
	"gitlab.com/spitertech/recommender/repository/helper"
	"gitlab.com/spitertech/recommender/service"
	"gitlab.com/spitertech/recommender/utils/log_utils"
	"log"
)

func InitRouters(r *gin.Engine) {

	r.Use(CORSMiddleware())
	r.Use(LoggingMiddleware())

	configurations := configurations.NewConfiguration("./configurations/config")
	log_utils.NewLogger(configurations.LoggerConfig)
	errInterceptor := error2.NewErrResponseInterceptor()
	elsErrorInterceptor := error2.NewElasticSearchErrorInterceptor()
	esClient := elasticSearchClient(configurations.ElasticSearchConfig)

	elsHelper := helper.NewElasticSearchHelper()
	elasticSearchRepo := repository.NewElasticSearchRepository(esClient, elsHelper)
	searchService := service.NewElasticSearchService(elasticSearchRepo, elsErrorInterceptor)
	searchController := controller.NewElasticSearchController(errInterceptor, searchService)

	elasticIndexerRepo := repository.NewEsDocumentIndexerRepository(esClient, elsHelper)
	indexerService := service.NewEsIndexerService(elasticIndexerRepo, elsErrorInterceptor)
	indexerController := controller.NewIndexerController(errInterceptor, indexerService)

	repository.NewEsMappingHandler(esClient).CreateMapping(context.Background())
	indexerHelper := repository.NewBulkDocumentIndexer(esClient, elsHelper)
	service.NewProductBulkIndexer(indexerHelper, configurations.BulkIndexerConfig, elsErrorInterceptor).ProductBulkIndexer(context.Background(), getProducts(),"product")

	r.GET("/health", func(c *gin.Context) {})
	apiVersion1 := r.Group("api/v1/")
	apiVersion1.GET("/search", searchController.Search)
	apiVersion1.GET("/search/suggest", searchController.RecommendSearch)
	apiVersion1.POST("/index-doc", indexerController.IndexDocument)
}

func getProducts() []models.Product {
	p1 := models.Product{ProductID: "p1", IsAvailable: true, ProductCategory: "electronics", ProductName: "MOBILE", ProductPrise: 20, Suggest: models.Suggest{Input: []string{"mobile", "mobile phone"}}}
	p2 := models.Product{ProductID: "p2", IsAvailable: true, ProductCategory: "electronics", ProductName: "DELL", ProductPrise: 100, Suggest: models.Suggest{Input: []string{"dell"}}}
	p3 := models.Product{ProductID: "p3", IsAvailable: true, ProductCategory: "furniture", ProductName: "SOFA", ProductPrise: 200.00, Suggest: models.Suggest{Input: []string{"mobile", "sofa"}}}
	p4 := models.Product{ProductID: "p4", IsAvailable: false, ProductCategory: "home appliances", ProductName: "SOFASET", ProductPrise: 300, Suggest: models.Suggest{Input: []string{"sofaset"}}}
	p5 := models.Product{ProductID: "p5", IsAvailable: true, ProductCategory: "beauty product", ProductName: "BODY SPRAY", ProductPrise: 180, Suggest: models.Suggest{Input: []string{"body spray"}}}
	p6 := models.Product{ProductID: "p6", IsAvailable: true, ProductCategory: "regular", ProductName: "SOYA OIL", ProductPrise: 30, Suggest: models.Suggest{Input: []string{"soya oil"}}}
	return []models.Product{p1, p2, p3, p4, p5, p6}

}

func elasticSearchClient(config configurations.ElasticSearchConfig) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://" + config.HostName + ":" + config.Port,
		},
		Username: "user",
		Password: "pass",
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return client
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.WithContext(log_utils.WithRqID(ctx))
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
