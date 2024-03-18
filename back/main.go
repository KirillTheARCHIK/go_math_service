package main

import (
	"go_math_service/src/endpoints"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.POST("/resolve", endpoints.ResolveEndpoint)

	server.Run() // listen and serve on 0.0.0.0:8080
}
