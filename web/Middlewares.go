package web

import (
	"btcRate/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func errorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute the next middleware or handler

		// Check if there was an error
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Check if the error is a CustomError
				if _, ok := e.Err.(*domain.EndpointInaccessibleError); ok {
					// If it is, respond with status code 400
					c.String(http.StatusBadRequest, e.Error())
				} else if _, ok := e.Err.(*domain.DataConsistencyError); ok {
					// If it is, respond with status code 400
					c.String(http.StatusConflict, e.Error())
				}
			}
		}
	}
}
