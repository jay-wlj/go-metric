package otelgin

import (
	"bytes"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"

	"github.com/jay-wlj/go-metric/internal/global"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func HTTPServerTimerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// serve the request to the next middleware
		c.Next()

		var retStr string
		retCode, err := jsonparser.GetInt(blw.body.Bytes(), "ret")
		if err != nil {
			retStr = "-"
		} else {
			retStr = strconv.Itoa(int(retCode))
		}

		global.GetMeter().Components().NewHTTPServerTimer(
			c.Request.URL.Path,
			retStr,
			c.Writer.Status(),
		).UpdateSince(start)
	}
}
