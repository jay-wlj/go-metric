package middleware

import (
	"bytes"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	gometric "github.com/jay-wlj/go-metric"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type bodyWriter struct {
	context.ResponseWriter
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

func UseMiddleware(ctx iris.Context) {
	start := time.Now()
	blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.ResponseWriter()}
	ctx.ResetResponseWriter(blw)

	ctx.Next()

	retCode, err := jsonparser.GetInt(blw.body.Bytes(), "ret")
	var retStr string
	if err != nil {
		retStr = "-"
	} else {
		retStr = strconv.Itoa(int(retCode))
	}

	gometric.GetGlobalMeter().Components().NewHTTPServerTimer(
		ctx.Path(),
		retStr,
		ctx.GetStatusCode(),
	).UpdateSince(start)

}
