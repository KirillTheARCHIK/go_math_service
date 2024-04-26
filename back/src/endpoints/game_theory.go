package endpoints

import (
	"go_math_service/src/entities/game_theory"

	"github.com/gin-gonic/gin"
)

type GameTheoryEndpointBody struct {
	Matrix          game_theory.GameTheoryMatrix `json:"matrix" binding:"required"`
	FirstStrategies int                          `json:"firstStrategies" binding:"required"`
}

func GameTheoryEndpoint(c *gin.Context) {
	var json GameTheoryEndpointBody
	var err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	strategyDistribution, err := json.Matrix.Solve(json.FirstStrategies)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"strategyDistribution": strategyDistribution,
	})
}
