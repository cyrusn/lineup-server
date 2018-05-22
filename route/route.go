package route

import (
	"net/http"
)

// Route store information of route
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

// Routes is slice of route
func Routes(s ScheduleStore) []Route {
	return []Route{
		Route{
			Path:    "/schedule/{classcode}",
			Methods: []string{"GET"},
			Handler: getScheduleHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"POST"},
			Handler: addScheduleHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"DELETE"},
			Handler: removeScheduleHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/priority/{priority}",
			Methods: []string{"PUT"},
			Handler: updatePriorityHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-complete",
			Methods: []string{"PUT"},
			Handler: toggleIsCompleteHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-notified",
			Methods: []string{"PUT"},
			Handler: toggleIsNotifiedHandler(s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-meeting",
			Methods: []string{"PUT"},
			Handler: toggleIsMeetingHandler(s),
		},
	}
}
