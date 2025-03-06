package server

import (
	"fmt"
	"net/http"
	"time"
)

// requestWrapper is a middleware function that wraps another http.Handler and logs information
// about the request and response.
func (s *Server) requestWrapper(next http.Handler, routeName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := newWriterWrapper(w)

		next.ServeHTTP(ww, r)

		msg := fmt.Sprintf("%d %s %s %s %s",
			ww.status,
			r.Method,
			r.RequestURI,
			routeName,
			time.Since(start),
		)

		switch {
		case ww.status >= 400:
			var msgWithErrInfo string
			if ww.err == nil {
				msgWithErrInfo = fmt.Sprintf("%s %s", msg, ww.errorCode)
			} else {
				msgWithErrInfo = fmt.Sprintf("%s %s (%s)", msg, ww.errorCode, ww.err.Error())
			}

			if ww.status >= 500 {
				s.log.Error(msgWithErrInfo)
			} else { // status >= 400
				s.log.Warning(msgWithErrInfo)
			}

		default:
			s.log.Info(msg)
		}
	})
}

// writerWrapper is used to grab information passed to an http.ResponseWriter.  This information is
// typically used for logging or metrics perposes.  Instances should be created by using the
// newWriterWrapper() function in order to get the appropriate initial values and functionality.
type writerWrapper struct {
	http.ResponseWriter
	status    int
	statusSet bool
	errorCode errorCode
	err       error
}

// newWriterWrapper creates, initializes, and returns a new writerWrapper.
func newWriterWrapper(w http.ResponseWriter) *writerWrapper {
	return &writerWrapper{w, 200, false, "", nil}
}

func (w *writerWrapper) WriteHeader(httpStatus int) {
	if !w.statusSet {
		w.status = httpStatus
		w.ResponseWriter.WriteHeader((httpStatus))
		w.statusSet = true
	}
}

func (w *writerWrapper) Write(body []byte) (int, error) {
	w.statusSet = true // status can not be set after first write
	return w.ResponseWriter.Write(body)
}
