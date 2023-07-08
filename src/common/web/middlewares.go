package web

import (
	"btcRate/coin/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const defaultErrorMessage = "Internal server error. Please try again later."

func errorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e := e.Err.(type) {
				case *domain.EndpointInaccessibleError:
					c.String(http.StatusBadRequest, e.Error())
				case *domain.DataConsistencyError:
					c.String(http.StatusConflict, e.Error())
				case *domain.ArgumentError:
					c.String(http.StatusBadRequest, e.Error())
				default:
					c.String(http.StatusInternalServerError, defaultErrorMessage)
				}
			}
		}
	}
}
