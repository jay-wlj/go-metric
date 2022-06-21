package middleware

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	gometric "github.com/jay-wlj/go-metric"
)

// customHandler implements the http.Handler interface and provides
// trace and metrics to beego web apps.
type customHandler struct {
	http.Handler
}

type respWriteWrapper struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *respWriteWrapper) Write(p []byte) (int, error) {
	w.body.Write(p)
	return w.ResponseWriter.Write(p)
}

func (w *respWriteWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// ServerHTTP calls the configured handler to serve HTTP for req to rr.
func (o *customHandler) ServeHTTP(rr http.ResponseWriter, req *http.Request) {
	start := time.Now()

	blw := &respWriteWrapper{body: bytes.NewBufferString(""), ResponseWriter: rr, statusCode: http.StatusOK}
	o.Handler.ServeHTTP(blw, req)

	var retStr string
	retCode, err := jsonparser.GetInt(blw.body.Bytes(), "ret")
	if err != nil {
		retStr = "-"
	} else {
		retStr = strconv.Itoa(int(retCode))
	}

	gometric.GetGlobalMeter().Components().NewHTTPServerTimer(
		req.URL.Path,
		retStr,
		blw.statusCode,
	).UpdateSince(start)
}

func NewOTelBeegoMiddleWare() func(h http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &customHandler{handler}
	}
}
