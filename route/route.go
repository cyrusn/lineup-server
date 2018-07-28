package route

import (
	"net/http"

	"github.com/cyrusn/lineup-system/route/handler"
)

// Route store information of route
type Route struct {
	Path    string
	Methods []string
	Scopes  []string
	Auth    bool
	Handler func(http.ResponseWriter, *http.Request)
}

// Store collect function for auth and scheudle
type Store struct {
	handler.AuthStore
	handler.ScheduleStore
}

// Routes is slice of route
func Routes(s *Store) []Route {
	return []Route{
		Route{
			Path:    "/auth/login",
			Methods: []string{"POST"},
			Scopes:  []string{},
			Auth:    false,
			Handler: handler.LoginHandler(s.AuthStore),
		},
		Route{
			// refresh jwt token, where key is the jwt key name in header
			Path:    "/auth/refresh/{key}",
			Methods: []string{"GET"},
			Scopes:  []string{},
			Auth:    true,
			Handler: handler.RefreshHandler(s.AuthStore),
		},
		Route{
			// use queries to get the informaton of classcode
			// e.g. ?classcode=3A&classcode=3D
			Path:    "/schedules",
			Methods: []string{"GET"},
			Scopes:  []string{},
			Auth:    true,
			Handler: handler.GetSchedulesHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"POST"},
			Scopes:  []string{},
			Auth:    true,
			Handler: handler.AddScheduleHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"DELETE"},
			Scopes:  []string{},
			Auth:    true,
			Handler: handler.RemoveScheduleHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/priority/{priority}",
			Methods: []string{"PUT"},
			Scopes:  []string{"teacher"},
			Auth:    true,
			Handler: handler.UpdatePriorityHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-complete",
			Methods: []string{"PUT"},
			Scopes:  []string{"teacher"},
			Auth:    true,
			Handler: handler.ToggleIsCompleteHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-notified",
			Methods: []string{"PUT"},
			Scopes:  []string{},
			Auth:    true,
			Handler: handler.ToggleIsNotifiedHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-meeting",
			Methods: []string{"PUT"},
			Scopes:  []string{"teacher"},
			Auth:    true,
			Handler: handler.ToggleIsMeetingHandler(s.ScheduleStore),
		},
	}
}
