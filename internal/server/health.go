package server

import "net/http"

func (s *Server) livez(w http.ResponseWriter, r *http.Request) {
	msg := "It's alive!"
	s.RespondWithMessage(r.Context(), w, http.StatusOK, msg)
}

func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	msg := "Ready to rock!"
	s.RespondWithMessage(r.Context(), w, http.StatusOK, msg)
}
