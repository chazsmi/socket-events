# Socket Events

This package handles async events that occur tracks whos listening to them and sends the event back to the listner via a websocket

Use case example would be an event comsumed on a queue that needed to be broadacsted to listeners who had registered an interest via a websocket.

## Example set up
```Go
    // Start listening for stock updates
    foo := "foo"
    bar := "bar"
    eventsHandler := events.NewEventHandler()
    go eventsHandler.Init()
    eventsHandler.EventRec <- events.Event{
        Ref: &foo,
        Data: map[string]interface{}{
            "foo": foo,
            "bar": bar,
        },
    }

// Web server & web socket set up...
```

