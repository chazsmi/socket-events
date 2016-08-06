package events

import (
	"encoding/json"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Connection struct {
	// To Reference an event to broadcast on
	Ref string
	// The Websocket Connection
	Ws *websocket.Conn
	// A channel to quit on to close the connection
	Done chan bool
}

func (c Connection) Close() {
	c.Done <- true
}

type Handler struct {
	// Channel to Receive the events on
	Receive    chan Event
	eventsLock sync.RWMutex

	// Store of events being listened for
	Store map[string]map[*websocket.Conn]*Connection
}

type Event struct {
	// The reference to be compared upon
	Ref  *string
	Data map[string]interface{}
}

func NewEventHandler() *Handler {
	return &Handler{
		Receive: make(chan Event),
		Store:   make(map[string]map[*websocket.Conn]*Connection),
	}
}

func (h *Handler) Init() {
	for {
		event := <-h.Receive
		if event.Ref == nil {
			log.Printf("Ref not provided or nil... %+v", event)
			continue
		}

		// Broadcast to the events with the event ref
		go h.BroadcastEvent(event, *event.Ref)
		// To all listeners that want all events
		go h.BroadcastEvent(event, "")
	}
}

func (h *Handler) BroadcastEvent(event Event, ref string) error {
	h.eventsLock.RLock()
	defer h.eventsLock.RUnlock()

	js, err := json.Marshal(event.Data)
	if err != nil {
		log.Println("Json marshal error after proto: ", err)
		return err
	}

	// Send the result back to every socket listening
	for _, conn := range h.Store[ref] {
		if err := websocket.JSON.Send(conn.Ws, string(js)); err != nil {
			conn.Ws.Close()
			return err
		}
	}
	return nil
}

func (h *Handler) RegisterEvent(c *Connection) {
	log.Println("Registering Handler from", c.Ws.LocalAddr().String())
	h.eventsLock.Lock()
	defer h.eventsLock.Unlock()
	if _, found := h.Store[c.Ref]; !found {
		h.Store[c.Ref] = make(map[*websocket.Conn]*Connection)
	}
	h.Store[c.Ref][c.Ws] = c
}

func (h *Handler) UnregisterEvent(c Connection) {
	h.eventsLock.Lock()
	defer h.eventsLock.Unlock()
	delete(h.Store[c.Ref], c.Ws)
}
