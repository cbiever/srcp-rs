package main

import (
	"net/http"

	"srcp-rs/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CreateSession",
		"POST",
		"/sessions",
		handlers.CreateSession,
	},
	Route{
		"CreateGL",
		"POST",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls",
		handlers.CreateGL,
	},
	Route{
		"GetGL",
		"GET",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		handlers.GetGL,
	},
	Route{
		"UpdateGL",
		"PUT",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		handlers.UpdateGL,
	},
	Route{
		"UpdateGL",
		"PATCH",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		handlers.UpdateGL,
	},
	Route{
		"DeleteGL",
		"DELETE",
		"/sessions/{sessionId:[0-9]+}/buses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		handlers.DeleteGL,
	},
}
