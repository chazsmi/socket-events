# Socket Events

This package handles async events, tracks subscribers listening for updates to them and sends the event back to the listener via a websocket

Use case example would be an event comsumed on a queue that needed to be broadacsted to listeners who had registered an interest via a websocket.

## Example set up
```Go
// Set up
foo := "foo"
bar := "bar"
eventsHandler := events.NewEventHandler()
go eventsHandler.Init()

// Event comes in via transport eg message bus

// Pass your event to the handlers channel
eventsHandler.Receive <- events.Event{
    Ref: &foo,
    Data: map[string]interface{}{
        "foo": foo,
        "bar": bar,
    },
}

// Web server & web socket set up...
```

