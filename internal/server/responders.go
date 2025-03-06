package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseMessage struct {
	ErrorCode  errorCode `json:"errorCode,omitempty"`
	Message    string    `json:"message"`
	StatusCode int       `json:"statusCode,omitempty"`
}

func (s *Server) newResponseMessage(msg string, statusCode int, errCode errorCode) responseMessage {
	if s.cfg.IncludeStatusCodeInMessages {
		return responseMessage{Message: msg, StatusCode: statusCode, ErrorCode: errCode}
	}
	return responseMessage{Message: msg, ErrorCode: errCode}
}

func (s *Server) RespondWithMessage(ctx context.Context, w http.ResponseWriter, httpStatus int, msg string) {
	s.RespondWithJSON(ctx, w, httpStatus, s.newResponseMessage(msg, httpStatus, ""))
}

func (s *Server) RespondWithJSON(ctx context.Context, w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if body == nil {
		// return an empty json object if there is no body
		w.WriteHeader(httpStatus)
		_, err := w.Write([]byte("{}"))
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		// this should be unreachable code, but if we get here we log an error and return a 500
		s.log.Error("unable to marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		// hardcoded json string since marshaling might be broken
		_, writeErr := w.Write([]byte(`{ "errorCode": "12345", "errorMessage": "unable to marshal response" }`))
		if writeErr != nil {
			s.log.Error(writeErr.Error())
		}
		return
	}

	w.WriteHeader(httpStatus)
	_, err = w.Write(bytes)
	if err != nil {
		s.log.Error(err.Error())
	}
}

func (s *Server) RespondWithError(ctx context.Context, w http.ResponseWriter, r *http.Request,
	httpStatus int, errCode errorCode, message string, err error) {

	if err != nil {
		if ww, ok := w.(*writerWrapper); ok {
			ww.errorCode = errCode
			ww.err = err
		}
	}

	s.RespondWithJSON(ctx, w, httpStatus, s.newResponseMessage(message, httpStatus, errCode))
}
