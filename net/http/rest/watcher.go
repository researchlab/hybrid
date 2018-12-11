package rest

import (
	"sync"
	"time"
)

type class string

type classMap struct {
	sync.RWMutex
	m map[class]*watchEventChannels
}

func newClassMap() *classMap {
	return &classMap{m: map[class]*watchEventChannels{}}
}

func (p *classMap) find(name class) *watchEventChannels {
	p.RLock()
	defer p.RUnlock()
	return p.m[name]
}

func (p *classMap) incrementRef(name class) *watchEventChannels {
	p.Lock()
	defer p.Unlock()
	wec := p.m[name]
	if wec == nil {
		wec = newWatchEventChannels()
		wec.incrementRef()
		p.m[name] = wec
	}
	return wec
}

func (p *classMap) decrementRef(name class) *watchEventChannels {
	p.Lock()
	defer p.Unlock()
	wec := p.m[name]
	if wec != nil {
		wec.incrementRef()
		if wec.refCount == 0 {
			delete(p.m, name)
		}
	}
	return wec
}

type watchEventChannels struct {
	sync.RWMutex
	refCount int
	m        map[int64]chan *WatchEvent
}

func newWatchEventChannels() *watchEventChannels {

	return &watchEventChannels{m: map[int64]chan *WatchEvent{}}
}

func (p *watchEventChannels) append(event *OrmEvent) {
	p.RLock()
	defer p.RUnlock()
	for wk, ch := range p.m {
		go func(ch chan *WatchEvent, wk int64) {
			ch <- &WatchEvent{WaitKey: wk, OrmEvent: event}
		}(ch, wk)
	}
}

func (p *watchEventChannels) findOrCreateWatchEvents(waitKey int64) chan *WatchEvent {
	p.Lock()
	defer p.Unlock()
	watchEvents := p.m[waitKey]
	if watchEvents == nil {
		watchEvents = make(chan *WatchEvent)
		p.m[waitKey] = watchEvents
	}
	return watchEvents
}

func (p *watchEventChannels) removeWatchEvents(waitKey int64) chan *WatchEvent {
	p.Lock()
	defer p.Unlock()
	watchEvents := p.m[waitKey]
	if watchEvents != nil {
		close(watchEvents)
		delete(p.m, waitKey)
	}
	return watchEvents
}

func (p *watchEventChannels) incrementRef() {
	p.refCount++
}

func (p *watchEventChannels) decrementRef() {
	p.refCount--
}

type OrmEventType int

const (
	orm_read OrmEventType = iota
	orm_create
	orm_delete
	orm_update
)

type OrmEvent struct {
	EventType OrmEventType `json:"eventType"`
	Class     class        `json:"class"`
	Data      interface{}  `json:"data"`
	ID        interface{}  `json:"id"` //the id of object that has been deleted. it's nil if EventType isn't orm_delete
}

func NewOrmEvent(eventType OrmEventType, className class, data interface{}) *OrmEvent {
	return &OrmEvent{EventType: orm_read, Class: className, Data: data}
}

// WatchEvent
type WatchEvent struct {
	WaitKey  int64
	OrmEvent *OrmEvent
	Error    string
}

// watcher
type watcher struct {
	m *classMap
}

func newWatcher() *watcher {
	return &watcher{m: newClassMap()}
}

func (p *watcher) NotifyCreate(className class, data interface{}) {
	wec := p.m.find(className)
	if wec != nil {
		wec.append(&OrmEvent{EventType: orm_create, Class: className, Data: data})
	}
}

func (p *watcher) NotifyUpdate(class class, data interface{}) {
	wec := p.m.find(class)
	if wec != nil {
		wec.append(&OrmEvent{EventType: orm_update, Class: class, Data: data})
	}
}

func (p *watcher) NotifyDelete(class class, id interface{}) {
	wec := p.m.find(class)
	if wec != nil {
		wec.append(&OrmEvent{EventType: orm_delete, Class: class, ID: id})
	}
}

func (p *watcher) Wait(wk int64, timeout int, class class, id interface{}) (int64, chan *WatchEvent) {
	if timeout <= 0 {
		timeout = 30
	}
	wec := p.m.incrementRef(class)
	if wk == 0 {
		wk = time.Now().UnixNano()
		p.ttl(int64(timeout), class, wk)
	}
	return wk, wec.findOrCreateWatchEvents(wk)
}

func (p *watcher) ttl(timeout int64, class class, waitKey int64) {
	go func() {
		<-time.After(time.Duration(timeout) * time.Second)
		wec := p.m.decrementRef(class)
		if wec != nil {
			wec.removeWatchEvents(waitKey)
		}
	}()
}
