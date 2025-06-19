package services

import (
	"sync"

	"github.com/a-h/templ"
	datastar "github.com/starfederation/datastar/sdk/go"
)

type SSEHub struct {
	mutex sync.Mutex
	// client can have multiple connections
	clients map[string][]*datastar.ServerSentEventGenerator
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[string][]*datastar.ServerSentEventGenerator),
	}
}

var UserSSEHub = NewSSEHub()

func (hub *SSEHub) Add(userId string, sse *datastar.ServerSentEventGenerator) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	hub.clients[userId] = append(hub.clients[userId], sse)
}

func (hub *SSEHub) Remove(userId string, sse *datastar.ServerSentEventGenerator) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	for i, c := range hub.clients[userId] {
		if c == sse {
			hub.clients[userId] = append(hub.clients[userId][:i], hub.clients[userId][i+1:]...)
			break
		}
	}

}

func (hub *SSEHub) BroadcastFragments(userId string, component templ.Component) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for _, sse := range hub.clients[userId] {
		sse.MergeFragmentTempl(component)
	}

}

func (hub *SSEHub) BroadcastSignals(userId string, signals []byte) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for _, sse := range hub.clients[userId] {
		sse.MergeSignals(signals)
	}

}

func (hub *SSEHub) ExcuteScript(userId string, script string) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for _, sse := range hub.clients[userId] {
		sse.ExecuteScript(script)
	}
}
