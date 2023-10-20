package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
)

/*
	func: Timeout
	description: timeout middlerware
*/
func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap origin c.Writer with custom timeoutWriter
		tw := &timeoutWriter{
			ResponseWriter: c.Writer,
			h:              make(http.Header),
		}
		c.Writer = tw

		// wrap the request context with a timeout
		// and update the context of c.Request
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		// create channels
		finished := make(chan struct{})        // to indicate handler finished
		panicChan := make(chan interface{}, 1) // used to handle panic if we can't recover

		// call Next(), which is handler function
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			// invoke handler function
			c.Next()

			// send finished signal
			finished <- struct{}{}
		}()

		// Handling Either Finished, Panic, or Timeout/context.Done()
		select {
		case <-panicChan:
			// if we can't recover from painc
			// send server internal error
			e := apperrors.NewInternal()

			// build response header
			tw.ResponseWriter.WriteHeader(e.Status())

			// build response body
			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// if handler finished, set headers and write resp
			tw.mu.Lock()
			defer tw.mu.Unlock()

			// copy Header map from tw.h to
			// the one from tw.ResponseWriter
			// for response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}

			// copy Header code from tw.code to
			// the one from tw.ResponseWriter
			// for response
			tw.ResponseWriter.WriteHeader(tw.code)

			// copy buffer from tw.wbuf to
			// the one from tw.ResponseWriter
			// for response
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			// timeout has occured, send errTimeout and write headers
			tw.mu.Lock()
			defer tw.mu.Unlock()

			// config "Content-Type" in header
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")

			// config header status code
			tw.ResponseWriter.WriteHeader(errTimeout.Status())

			// set error reason in response body
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)

			// abort following handlers in handler chain
			c.Abort()

			// set timeout status to prevent the calls to tw.Write()
			// or tw.WriteHeader() to write
			tw.SetTimedout()
		}
	}
}

/*
	struct: timeoutWriter
	description:
		implements http.ResponseWriter, but tracks if Writer has timed out
		or has already written its header to prevent
		header and body overwrites
		also locks access to this writer to prevent race conditions
		holds the gin.ResponseWriter which we'll manually call Write()
		on in the middleware function to send response
*/
type timeoutWriter struct {
	gin.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

/*
	func: Write
	description: write to wbuf if context hasn't timeout
*/
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}

/*
	func: WriteHeader
	description: write header code if context hasn't timeout
*/
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)

	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut || tw.wroteHeader {
		return
	}

	tw.writeHeader(code)
}

/*
	func: Header
	description: get wrapped header map
*/
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

/*
	func: writeHeader
	description: write header code
			(internal invoke by WriteHeader)
*/
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

/*
	func: SetTimedout
	description: set timeout status
*/
func (tw *timeoutWriter) SetTimedout() {
	tw.timedOut = true
}

/*
	func: checkWriteHeaderCode
	description: check whether header code valid
			(internal invoke by WriteHeader)
*/
func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
