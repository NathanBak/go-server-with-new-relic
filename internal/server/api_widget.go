package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/NathanBak/go-server-with-new-relic/pkg/widget"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type widgetIDHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (s *Server) widgetIDMiddleware(next widgetIDHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wdgtID := mux.Vars(r)["widgetId"]

		// It should not be possible to get an empty wdgtID here unless there's a problem with the
		// routing rules
		if wdgtID == "" {
			s.RespondWithError(r.Context(), w, r, http.StatusInternalServerError,
				WS1000005, "no widget id provided", nil)
			return
		}

		id, err := uuid.Parse(wdgtID)
		if err != nil {
			s.RespondWithError(r.Context(), w, r, http.StatusBadRequest,
				WS1000006, "invalid widget id", err)
			return
		}

		next(w, r, id)
	}
}

func (s *Server) addWidget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusBadRequest,
			WS1000001, "unable to read body", err)
		return
	}

	wdgtInfo := struct {
		Name  string       `json:"name"`
		Color widget.Color `json:"color"`
	}{}

	if err := json.Unmarshal(body, &wdgtInfo); err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusBadRequest,
			WS1000002, "invalid widget provided", err)
		return
	}

	wdgt := widget.New(wdgtInfo.Name, wdgtInfo.Color)

	err = s.storage.Set(wdgt.ID.String(), wdgt)
	if err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusInternalServerError,
			WS1000010, "unable to store widget", err)
		return

	}

	s.RespondWithJSON(ctx, w, http.StatusOK, &wdgt)
}

func (s *Server) deleteWidget(w http.ResponseWriter, r *http.Request, wdgtID uuid.UUID) {
	ctx := r.Context()
	wdgt, ok, err := s.storage.Delete(wdgtID.String())
	if err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusInternalServerError,
			WS1000007, "unable to delete widget", nil)
		return
	}

	if !ok {
		s.RespondWithError(ctx, w, r, http.StatusNotFound,
			WS1000003, "widget not found", nil)
		return
	}

	s.RespondWithJSON(ctx, w, http.StatusOK, &wdgt)
}

func (s *Server) getWidget(w http.ResponseWriter, r *http.Request, wdgtID uuid.UUID) {
	ctx := r.Context()

	wdgt, ok, err := s.storage.Get(wdgtID.String())
	if err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusInternalServerError,
			WS1000008, "unable to retrieve widget", nil)
		return
	}
	if !ok {
		s.RespondWithError(ctx, w, r, http.StatusNotFound,
			WS1000004, "widget not found", nil)
		return
	}

	s.RespondWithJSON(ctx, w, http.StatusOK, &wdgt)
}

func (s *Server) listWidgets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := s.storage.Keys()
	if err != nil {
		s.RespondWithError(r.Context(), w, r, http.StatusInternalServerError,
			WS1000009, "unable to retrieve widgets", nil)
		return
	}

	s.RespondWithJSON(ctx, w, http.StatusOK, &ids)
}
