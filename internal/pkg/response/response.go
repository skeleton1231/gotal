package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skeleton1231/gotal/internal/pkg/errors"
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
		logrus.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), APIErrorResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}
