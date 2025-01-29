package api

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func JSONStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

type errorResponse struct {
	Message string `json:message`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logger := zerolog.New(os.Stdout).Level(zerolog.TraceLevel)
	logger.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
