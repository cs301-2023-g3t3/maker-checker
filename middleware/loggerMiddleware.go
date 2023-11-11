package middleware

import (
	"net/http"
	"strings"
	"time"

	"makerchecker-api/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CustomData struct {
	Method  string
	URL     string
	Headers map[string][]string
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()

		// Process the request
		ctx.Next()

		// End Time request
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime).Milliseconds()

		// Request data
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		metadata, ok := ctx.Request.Context().Value("RequestMetadata").(models.RequestMetadata)
		var userAgent string
		var sourceIP string
		if ok {
			// Access the UserAgent and SourceIP
			userAgent = metadata.UserAgent
			sourceIP = metadata.SourceIP
		}

		// Logs for makerchecker entries
		if strings.Contains(reqUri, "/record") {
			data, _ := ctx.Get("makerchecker")
			makerchecker, _ := data.(models.Makerchecker)
			makecheckerFields := log.Fields{
				"id":           makerchecker.Id,
				"makerId":      makerchecker.MakerId,
				"makerEmail":   makerchecker.MakerEmail,
				"checkerId":    makerchecker.CheckerId,
				"checkerEmail": makerchecker.CheckerEmail,
				"endpoint":     makerchecker.Endpoint,
				"status":       makerchecker.Status,
				"data":         makerchecker.Data,
			}

			if reqMethod == http.MethodPost {
				log.WithFields(log.Fields{
					"METHOD":               reqMethod,
					"URI":                  reqUri,
					"STATUS":               statusCode,
					"LATENCY":              latencyTime,
					"MAKERCHECKER_DETAILS": makecheckerFields,
					"ACTION":              	"create",
					"USER_AGENT":           userAgent,
					"SOURCE_IP":            sourceIP,
				}).Info("CREATE MAKERCHECKER REQUEST")
			}

			if reqMethod == http.MethodPut {
				log.WithFields(log.Fields{
					"METHOD":               reqMethod,
					"URI":                  reqUri,
					"STATUS":               statusCode,
					"LATENCY":              latencyTime,
					"MAKERCHECKER_DETAILS": makecheckerFields,
					"ACTION":               makecheckerFields["status"],
					"USER_AGENT":           userAgent,
					"SOURCE_IP":            sourceIP,
				}).Info("UPDATE MAKERCHECKER REQUEST")
			}

		}
	}
}
