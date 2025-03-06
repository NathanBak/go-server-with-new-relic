package server

import "net/http"

// A Route describes a rest endpoint.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// routes returns an array of all Routes.
func (s *Server) routes() []Route {
	return []Route{

		{
			Name:        "Livez",
			Method:      http.MethodGet,
			Pattern:     "/livez",
			HandlerFunc: s.livez,
		},

		{
			Name:        "Readyz",
			Method:      http.MethodGet,
			Pattern:     "/readyz",
			HandlerFunc: s.readyz,
		},

		{
			Name:        "AddWidget",
			Method:      http.MethodPost,
			Pattern:     "/api/v1/widgets",
			HandlerFunc: s.addWidget,
		},

		{
			Name:        "DeleteWidget",
			Method:      http.MethodDelete,
			Pattern:     "/api/v1/widgets/{widgetId}",
			HandlerFunc: s.widgetIDMiddleware(s.deleteWidget),
		},

		{
			Name:        "GetWidget",
			Method:      http.MethodGet,
			Pattern:     "/api/v1/widgets/{widgetId}",
			HandlerFunc: s.widgetIDMiddleware(s.getWidget),
		},

		{
			Name:        "ListWidgets",
			Method:      http.MethodGet,
			Pattern:     "/api/v1/widgets",
			HandlerFunc: s.listWidgets,
		},
	}
}
