package main

import "net/http"

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
		CreateSession,
	},
	Route{
		"CreateGL",
		"POST",
		"/sessions/{sessionId:[0-9]+}/gls",
		CreateGL,
	},
	Route{
		"GetGL",
		"GET",
		"/sessions/{sessionId:[0-9]+}/busses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		GetGL,
	},
	Route{
		"UpdateGL",
		"PUT",
		"/sessions/{sessionId:[0-9]+}/busses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		UpdateGL,
	},
	Route{
		"DeleteGL",
		"DELETE",
		"/sessions/{sessionId:[0-9]+}/busses/{bus:[0-9]+}/gls/{address:[0-9]+}",
		DeleteGL,
	},
}
