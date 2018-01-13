package main

import "net/http"

// Route ...
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

var routes = []Route{
	Route{
		Path:    "/schedule",
		Methods: []string{"GET"},
		Handler: getScheduleHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}",
		Methods: []string{"POST"},
		Handler: addScheduleHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}",
		Methods: []string{"DELETE"},
		Handler: removeScheduleHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/order/{order}",
		Methods: []string{"PUT"},
		Handler: updateOrderHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-complete",
		Methods: []string{"PUT"},
		Handler: toggleIsCompleteHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-notified",
		Methods: []string{"PUT"},
		Handler: toggleIsNotifiedHandler,
	},
	Route{
		Path:    "/schedule/{classcode}/{classno}/is-meeting",
		Methods: []string{"PUT"},
		Handler: toggleIsMeetingHandler,
	},
}
