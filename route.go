package main

import "net/http"

// Route ...
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

func (hub *Hub) routes() []Route {
	return []Route{
		Route{
			Path:    "/schedule",
			Methods: []string{"GET"},
			Handler: hub.getScheduleHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"POST"},
			Handler: hub.addScheduleHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"DELETE"},
			Handler: hub.removeScheduleHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/order/{order}",
			Methods: []string{"PUT"},
			Handler: hub.updateOrderHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-complete",
			Methods: []string{"PUT"},
			Handler: hub.toggleIsCompleteHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-notified",
			Methods: []string{"PUT"},
			Handler: hub.toggleIsNotifiedHandler,
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-meeting",
			Methods: []string{"PUT"},
			Handler: hub.toggleIsMeetingHandler,
		},
	}
}
