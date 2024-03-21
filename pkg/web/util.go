package web

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebDiffPreviewer) log(loggerParams gin.LogFormatterParams) string {
	log.Info().
		Str("clientIP", loggerParams.ClientIP).
		Str("method", loggerParams.Method).
		Str("path", loggerParams.Path).
		Int("status", loggerParams.StatusCode).
		Str("latency", loggerParams.Latency.String()).
		Msg(loggerParams.ErrorMessage)
	return ""
}
