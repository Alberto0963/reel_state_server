package controllers

import (
	"net/http"
	"reelState/services"

	"github.com/gin-gonic/gin"
)

func GenerateResponse(c *gin.Context) {
	message := c.Query("message")

	// Llamar a la función que genera la respuesta
	response, err := services.GenerateResponse(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver la respuesta en formato JSON
	c.JSON(http.StatusOK, gin.H{"response": response})
}
