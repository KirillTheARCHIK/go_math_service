package endpoints

import (
	"go_math_service/src/entities/math_function"
	"strings"

	"github.com/gin-gonic/gin"
)

type DerivativeEndpointBody struct {
	Input          string             `json:"input" binding:"required"`
	VariableValues map[string]float64 `json:"variableValues" binding:"required"`
	Axis           string             `json:"axis" binding:"required"`
}

func DerivativeEndpoint(c *gin.Context) {
	var json DerivativeEndpointBody
	var err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	json.Input = strings.ReplaceAll(json.Input, " ", "")
	json.Input = " " + json.Input
	var mathFunction = math_function.Constructor(json.Input, true)
	value, err := mathFunction.GetDerivative(json.VariableValues, json.Axis)
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
