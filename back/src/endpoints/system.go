package endpoints

import (
	"go_math_service/src/entities/equality_system"

	"github.com/gin-gonic/gin"
)

type SystemEndpointBody struct {
	A [][]float64 `json:"a" binding:"required"`
	B []float64   `json:"b" binding:"required"`
}

func SystemEndpoint(c *gin.Context) {
	var json SystemEndpointBody
	var err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	////////////////////////////////////
	X, err := equality_system.Solve(json.A, json.B)
	////////////////////////////////////
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"x": X,
	})
}
