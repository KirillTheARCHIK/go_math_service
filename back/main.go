package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	// server.GET("/checkLogin", endpoints.CheckLogin)

	server.Run() // listen and serve on 0.0.0.0:8080
}
