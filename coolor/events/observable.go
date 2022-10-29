package events

import (
	"fmt"
	// "unsafe"

	"github.com/gdamore/tcell/v2"
	"go.uber.org/zap/zapcore"

	. "github.com/digitallyserviced/coolors/coolor/sequenced"
)

func init() {
}


type EventHandlers []tcell.EventHandler

type (
	ObservableEvent struct {
		Ref interface{}
		Src interface{}
		*tcell.EventTime
		Note string
		Type ObservableEventType
		SequencedItem
	}
	Event struct {
		tcell.EventTime
	}

	EventObserver struct {
		observerName string
		SequencedItem
	}

	Observer interface {
		// tcell.Event
		Name() string
		HandleEvent(ObservableEvent) bool
	}
	Referenced interface {
		GetRef() interface{}
	}
)

func NewAnonymousHandlerFunc(f func(e ObservableEvent) bool) *AnonymousHandler {
	ahh := &AnonymousHandler{
		Callback: f,
	}
	return ahh
}


func NewAnonymousHandler(callbacks []Observer) *AnonymousObserver {
	ah := &AnonymousObserver{
		EventObserver: NewEventObserver("anon"),
	}
	return ah
}

func NewEventObserver(name string) *EventObserver {
	// unsafe.Pointer
	eo := &EventObserver{
		observerName: name,
		SequencedItem: SequencedItem{
			SeqNo: 0,
		},
	}
	return eo
}

func (psm ObservableEvent) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString("event.name", psm.Note)
	oe.AddString("event.type", psm.Type.String())
	return nil
}
func (o *ObservableEvent) String() string {
	return fmt.Sprintf("anon  %s %s %s %v %v %d", o.EventTime, o.Type.String(), o.Note, o.Ref, o.Src, o.SeqNo)
}

func (o *EventObserver) Name() string {
	return fmt.Sprintf("%s @ #%d", o.observerName, o.SeqNo)
}

func (o *EventObserver) HandleEvent(e ObservableEvent) bool {
	return true
}
type AnonymousObserver struct {
	*EventObserver
	Callbacks []Observer
}

type AnonymousHandler struct {
	Callback func(e ObservableEvent) bool
}

func (ah *AnonymousHandler) Name() string {
	return fmt.Sprintf("%s @ #%d", "Anon handler", 1)
}

func (ah *AnonymousHandler) HandleEvent(e ObservableEvent) bool {
	ah.Callback(e)
	return true
}

func (ah *AnonymousObserver) ObserverFunc(
	n Notifier,
	t ObservableEventType,
	f func(e ObservableEvent) bool,
) *AnonymousObserver {
	ahh := &AnonymousHandler{
		Callback: func(e ObservableEvent) bool {
			return f(e)
		},
	}
	n.Register(t, ahh)
	ah.Callbacks = append(ah.Callbacks, ahh)
	return ah
}
