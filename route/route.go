package route

import (
	"net/http"

	"github.com/cyrusn/lineup-system/hub"
	"github.com/cyrusn/lineup-system/schedule"
)

// Route store information of route
type Route struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

// Routes is slice of route
func Routes(hub *hub.Hub, s schedule.MapSchedules) []Route {
	return []Route{
		Route{
			Path:    "/schedule",
			Methods: []string{"GET"},
			Handler: getScheduleHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"POST"},
			Handler: addScheduleHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"DELETE"},
			Handler: removeScheduleHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/order/{order}",
			Methods: []string{"PUT"},
			Handler: updateOrderHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-complete",
			Methods: []string{"PUT"},
			Handler: toggleIsCompleteHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-notified",
			Methods: []string{"PUT"},
			Handler: toggleIsNotifiedHandler(hub, s),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-meeting",
			Methods: []string{"PUT"},
			Handler: toggleIsMeetingHandler(hub, s),
		},
	}
}
