package main

import (
	"net/http"

	"srcp-rs/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	ContentType string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CreateSession",
		"POST",
		"/sessions",
		"application/json; charset=UTF-8",
		handlers.CreateSession,
	},
	Route{
		"GetBuses",
		"GET",
		"/sessions/{sessionId:[0-9]+}/buses",
		"application/json; charset=UTF-8",
		handlers.GetBuses,
	},
	Route{
		"DeleteBus",
		"DELETE",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}",
		"application/json; charset=UTF-8",
		handlers.DeleteBus,
	},
	Route{
		"CreateGL",
		"POST",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls",
		"application/json; charset=UTF-8",
		handlers.CreateGL,
	},
	Route{
		"GetGL",
		"GET",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		"application/json; charset=UTF-8",
		handlers.GetGL,
	},
	Route{
		"UpdateGL",
		"PUT",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		"application/json; charset=UTF-8",
		handlers.UpdateGL,
	},
	Route{
		"UpdateGL",
		"PATCH",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		"application/json; charset=UTF-8",
		handlers.UpdateGL,
	},
	Route{
		"DeleteGL",
		"DELETE",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		"application/json; charset=UTF-8",
		handlers.DeleteGL,
	},
	Route{
		"GetConfiguration",
		"GET",
		"/sessions/{sessionId:[0-9]+}/configuration",
		"text/plain; charset=UTF-8",
		handlers.GetConfiguration,
	},
	Route{
		"UpdateConfiguration",
		"POST",
		"/sessions/{sessionId:[0-9]+}/configuration",
		"text/plain; charset=UTF-8",
		handlers.UpdateConfiguration,
	},
}
