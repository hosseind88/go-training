package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func handleValidationError(c *gin.Context, err error) bool {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make([]string, 0)
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages,
					e.Field()+" is required")
			case "email":
				errorMessages = append(errorMessages,
					"Invalid email format")
			case "min":
				errorMessages = append(errorMessages,
					e.Field()+" must be at least "+e.Param()+" characters long")
			case "len":
				errorMessages = append(errorMessages,
					e.Field()+" must be exactly "+e.Param()+" characters long")
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return true
	}
	return false
}
