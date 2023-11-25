package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
	"github.com/skeleton1231/gotal/pkg/log"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type APIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteResponse writes the response (data or error) into the HTTP response body.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		errResponse := APIErrorResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		}
		// c.Set("response", errResponse)
		c.JSON(coder.HTTPStatus(), errResponse)
		return
	}

	response := APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}
	c.Set("response", response)
	c.JSON(http.StatusOK, response)
}
