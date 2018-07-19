package route

import (
	"net/http"

	auth "github.com/cyrusn/goJWTAuthHelper"
	"github.com/cyrusn/lineup-system/route/handler"
)

// Route store information of route
type Route struct {
	Path    string
	Methods []string
	Auth    bool
	Handler func(http.ResponseWriter, *http.Request)
}

// Store collect function for auth and scheudle
type Store struct {
	handler.AuthStore
	handler.ScheduleStore
	auth.Secret
}

// Routes is slice of route
func Routes(s *Store) []Route {
	return []Route{
		Route{
			Path:    "/auth/login",
			Methods: []string{"POST"},
			Auth:    false,
			Handler: handler.LoginHandler(s.AuthStore, s.Secret),
		},
		Route{
			// refresh jwt token
			Path:    "/auth/refresh",
			Methods: []string{"GET"},
			Auth:    true,
			Handler: handler.RefreshHandler(s.AuthStore, s.Secret),
		},
		Route{
			Path:    "/schedule/{classcode}",
			Methods: []string{"GET"},
			Auth:    true,
			Handler: handler.GetScheduleHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"POST"},
			Auth:    true,
			Handler: handler.AddScheduleHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}",
			Methods: []string{"DELETE"},
			Auth:    true,
			Handler: handler.RemoveScheduleHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/priority/{priority}",
			Methods: []string{"PUT"},
			Auth:    true,
			Handler: handler.UpdatePriorityHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-complete",
			Methods: []string{"PUT"},
			Auth:    true,
			Handler: handler.ToggleIsCompleteHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-notified",
			Methods: []string{"PUT"},
			Auth:    true,
			Handler: handler.ToggleIsNotifiedHandler(s.ScheduleStore),
		},
		Route{
			Path:    "/schedule/{classcode}/{classno}/is-meeting",
			Methods: []string{"PUT"},
			Auth:    true,
			Handler: handler.ToggleIsMeetingHandler(s.ScheduleStore),
		},
	}
}
