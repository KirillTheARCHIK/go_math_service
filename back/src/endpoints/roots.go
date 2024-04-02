package endpoints

import (
	"go_math_service/src/entities/math_function"
	"strings"

	"github.com/gin-gonic/gin"
)

type FindRootsEndpointBody struct {
	Input string  `json:"input" binding:"required"`
	Type  string  `json:"type"`
	Eps   float64 `json:"eps" binding:"required"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	// VariableValues map[string]float64 `json:"variableValues" binding:"required"`
}

func FindRootsEndpoint(c *gin.Context) {
	var json FindRootsEndpointBody
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

	var value []float64
	switch json.Type {
	case "dividing":
		{
			value, err = mathFunction.FindRootsDividing(json.A, json.B, json.Eps)
		}
	}
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
