package brick

import (
	"log"
	"runtime/debug"
	"sync/atomic"
)

// EventHandler process the event
type EventHandler interface {
	// Handle process the event
	Handle(event string, data interface{})
}

// HandleEvent func
type HandleEvent func(event string, data interface{})

// EventHandlerFunc create EventHandler
func EventHandlerFunc(handle HandleEvent) EventHandler {
	return &eventHandlerWrap{handle}
}

type eventHandlerWrap struct {
	f HandleEvent
}

func (p *eventHandlerWrap) Handle(event string, data interface{}) {
	p.f(event, data)
}

// Observable event observable interface
type Observable interface {
	// Register a handler on the event
	On(event string, handler EventHandler)

	// Remove the handler on the event
	Off(event string, handler EventHandler)
}

// Notify events
type Notify interface {
	// Emmit an event
	Emmit(event string, data interface{})
}

type eventObject struct {
	event string
	data  interface{}
}

//Trigger events trigger struct
type Trigger struct {
	events map[string][]EventHandler
	ech    chan *eventObject
	init   uint64
}

//On Trigger on function
func (p *Trigger) On(event string, handler EventHandler) {
	if p.events == nil {
		p.events = map[string][]EventHandler{}
	}
	hs := p.events[event]
	if hs == nil {
		hs = []EventHandler{}
	}
	hs = append(hs, handler)
	p.events[event] = hs
}

//Off Trigger off function
func (p *Trigger) Off(event string, handler EventHandler) {
	if p.events == nil {
		return
	}
	hs := p.events[event]
	if hs == nil {
		return
	}
	for i, h := range hs {
		if h != handler {
			hs = append(hs[:i], hs[i+1:]...)
			break
		}
	}
	p.events[event] = hs
}

//Emmit Event trigger emmit
func (p *Trigger) Emmit(event string, data interface{}) {
	//fmt.Println("Trigger.Emmit", event)
	p.runDispatch()
	p.ech <- &eventObject{event, data}
}

func (p *Trigger) callHandle(h EventHandler, event string, data interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("emit event error.%s,%v.\nrecover:\n%v\n", event, data, r)
			debug.PrintStack()
		}
	}()
	h.Handle(event, data)
}

func (p *Trigger) runDispatch() {
	if atomic.CompareAndSwapUint64(&p.init, 0, 1) {
		p.ech = make(chan *eventObject, 1000)
		go func() {
			for e := range p.ech {
				if p.events == nil {
					continue
				}
				hs := p.events[e.event]
				if hs == nil {
					continue
				}
				for _, h := range hs {
					p.callHandle(h, e.event, e.data)
				}
			}
		}()
	}
}
