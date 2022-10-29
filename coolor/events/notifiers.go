package events

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/digitallyserviced/coolors/coolor/sequenced"
)

type (
	Notifier interface {
		Register(ObservableEventType, Observer)
		Deregister(Observer)
		Notify(ObservableEvent)
	}
	EventNotifier struct {
		observers    map[Observer]ObservableEventType
		NotifierName string
		SequencedItem
	}
)

var Global = struct{ *EventNotifier }{
	EventNotifier: NewEventNotifier("global"),
}

func NewEventNotifier(name string) *EventNotifier {
	eo := &EventNotifier{
		observers:     make(map[Observer]ObservableEventType),
		NotifierName:  name,
		SequencedItem: SequencedItem{SeqNo: 0},
	}
	return eo
}

func (o *EventNotifier) Register(t ObservableEventType, l Observer) {
	o.observers[l] = t
}

func (o *EventNotifier) Deregister(l Observer) {
	delete(o.observers, l)
}

func (p *EventNotifier) Notify(e ObservableEvent) {
	for o, _ := range p.observers {
		if !o.HandleEvent(e) {
		} else {
    }
	}
}

func (n *EventNotifier) NewObservableEvent(
	t ObservableEventType,
	note string,
	ref Referenced,
	src Referenced,
) *ObservableEvent {
	oe := &ObservableEvent{
		EventTime:     &tcell.EventTime{},
		Type:          t,
		SequencedItem: SequencedItem{SeqNo: n.Next()},
		Note:          note,
	}
	if ref != nil {
		oe.Ref = ref.GetRef()
	}
	if src != nil {
		oe.Src = src.GetRef()
	}
	oe.SetEventNow()
	return oe
}
