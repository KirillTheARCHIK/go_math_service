package endpoints

import (
	"go_math_service/src/entities/math_function"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResolveEndpointBody struct {
	Input          string             `json:"input" binding:"required"`
	VariableValues map[string]float64 `json:"variableValues" binding:"required"`
}

func ResolveEndpoint(c *gin.Context) {
	var json ResolveEndpointBody
	var err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	json.Input = strings.ReplaceAll(json.Input, " ", "")
	var mathFunction = math_function.Constructor(json.Input, true)
	value, err := mathFunction.ResolveExpression(json.VariableValues)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"value": value,
	})
}
