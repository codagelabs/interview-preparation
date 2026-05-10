package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/spitertech/recommender/router"
)

func main() {
	r := gin.Default()
	router.InitRouters(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
