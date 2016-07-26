package main

import "github.com/chazsmi/socket-events/events"

func main() {

	foo := "foo"
	bar := "bar"
	eventsHandler := events.NewEventHandler()
	eventsHandler.Init()
	eventsHandler.EventRec <- events.Event{
		Ref: &foo,
		Data: map[string]interface{}{
			"foo": foo,
			"bar": bar,
		},
	}

	// Web

}
