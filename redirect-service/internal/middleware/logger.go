package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[97;41m"
	colorGreen   = "\033[97;42m"
	colorYellow  = "\033[90;43m"
	colorBlue    = "\033[97;44m"
	colorCyan    = "\033[97;46m"
	colorMagenta = "\033[97;45m"
	colorWhite   = "\033[90;47m"
)

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 300 && code < 400:
		return colorWhite
	case code >= 400 && code < 500:
		return colorYellow
	default:
		return colorRed
	}
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorCyan
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorGreen
	default:
		return colorMagenta
	}
}

// ColoredLoggerWithFields logs one line per request with colored status & method
// plus all structured fields.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Gather data
		ts, requestID, userID := time.Now().UTC().Format(time.RFC3339Nano), c.GetString("request_id"), c.GetString("user_uuid")
		method, path, query := c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery
		status, latency := c.Writer.Status(), time.Since(start)
		clientIP, ua := c.ClientIP(), c.Request.UserAgent()
		size := c.Writer.Size()
		if size < 0 {
			size = 0
		}
		errs := c.Errors.Errors()

		sc, mc := statusColor(status), methodColor(method)

		fmt.Printf(
			"[TIMESTAMP] %s |%s %3d %s| %7v | %s |%s %-7s %s| PATH=\"%s\" REQ_ID=\"%s\" USER_UUID=\"%s\" QUERY=\"%s\" BYTES=%d UA=\"%s\" ERRORS=%v\n",
			ts,
			sc, status, colorReset,
			latency,
			clientIP,
			mc, method, colorReset,
			path,
			requestID,
			userID,
			query,
			size,
			ua,
			errs,
		)
	}
}
