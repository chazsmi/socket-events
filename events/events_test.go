package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

var ref string = "test"

var testEvent = Event{
	Ref: &ref,
	Data: map[string]interface{}{
		"foo": "bar",
		"bar": "foo",
	},
}

type Response struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func mockClient(server *httptest.Server, ch chan Response) {
	ws, err := websocket.Dial(
		fmt.Sprintf("ws://%s/?test=test", strings.Replace(server.URL, "http://", "", 1)),
		"",
		"http://localhost/",
	)
	if err != nil {
		log.Fatal(err)
	}
	var msg string
	time.Sleep(1000 * time.Millisecond)
	jsonErr := websocket.JSON.Receive(ws, &msg)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var resp Response
	marshErr := json.Unmarshal([]byte(msg), &resp)
	if marshErr != nil {
		log.Fatal(marshErr)
	}
	ch <- resp
}

func mockServer(handle func(*websocket.Conn)) *httptest.Server {
	return httptest.NewServer(websocket.Handler(handle))
}

func (eh Handler) mockHandler(ws *websocket.Conn) {
	conn := &Connection{
		Ref: ref,
		Ws:  ws,
	}
	eh.RegisterEvent(conn)
	for {
		select {
		case <-conn.Done:
			// Kill the connection upon termination
			log.Println("Killing connection")
			break
		}
	}

}

func TestInit(t *testing.T) {
	eventsHandler := NewEventHandler()
	server := mockServer(eventsHandler.mockHandler)
	go eventsHandler.Init()

	respChan := make(chan Response)
	defer close(respChan)

	// Send a request to the mocked server
	go mockClient(server, respChan)
	// delay to allow cleint to make the request
	time.Sleep(500 * time.Millisecond)
	// Send a fake event back throught the chan
	eventsHandler.Receive <- testEvent

	for {
		resp := <-respChan
		if resp.Foo != "bar" {
			t.Fail()
		}
		break
	}
}
