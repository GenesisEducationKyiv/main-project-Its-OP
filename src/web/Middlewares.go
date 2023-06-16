package web

import (
	"btcRate/domain"
	"github.com/gin-gonic/gin"
	"log"
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
					c.String(http.StatusBadRequest, e.Error())
				} else if _, ok := e.Err.(*domain.DataConsistencyError); ok {
					c.String(http.StatusConflict, e.Error())
				} else if _, ok := e.Err.(*domain.DatabaseError); ok {
					c.String(http.StatusInternalServerError, e.Error())
					nestedErr := e.Unwrap()
					log.Printf("ERROR: Database error, %v", nestedErr)
				}
			}
		}
	}
}
