package coolor

import (
	// "log"

	"log"

	"github.com/gdamore/tcell/v2"
)

type (
	Notifier interface {
		Register(ObservableEventType, Observer)
		Deregister(Observer)
		Notify(ObservableEvent)
	}
	eventNotifier struct {
		observers    map[Observer]ObservableEventType
		notifierName string
		SequencedItem
	}
)

func NewEventNotifier(name string) *eventNotifier {
	eo := &eventNotifier{
		observers:     make(map[Observer]ObservableEventType),
		notifierName:  name,
		SequencedItem: SequencedItem{seqNo: 0},
	}
	return eo
}

func (o *eventNotifier) Register(t ObservableEventType, l Observer) {
	o.observers[l] = t
}

func (o *eventNotifier) Deregister(l Observer) {
	delete(o.observers, l)
}

func (p *eventNotifier) Notify(e ObservableEvent) {
	// log.Printf(
	// 	"*** Observer %s notified: \n %T  %T %s******",
	// 	e.Type.String(),
	// 	e.Ref,
	// 	e.Src,
	// 	e.Note,
	// )
	for o, f := range p.observers {
		log.Printf("notifier -- fns: %d  name: %s ev: %v flag: %d", len(p.observers), p.notifierName, e, e.Type&f)
		// if e.Type&f == 0 {
		// 	continue
		// }
		if !o.HandleEvent(e) {
		}
	}
}

func (n *eventNotifier) NewObservableEvent(
	t ObservableEventType,
	note string,
	ref Referenced,
	src Referenced,
) *ObservableEvent {
	oe := &ObservableEvent{
		// Ref:           ref.GetRef(),
		// Src:           src.GetRef(),
		EventTime:     &tcell.EventTime{},
		Type:          t,
		SequencedItem: SequencedItem{seqNo: n.Next()},
		Note:          note,
	}
	if ref != nil {
		oe.Ref = ref.GetRef()
	}
	if src != nil {
		oe.Src = src.GetRef()
	}
	oe.SetEventNow()
	// fmt.Printf("*** Observer %d received: \n %T  %T %s*****", oe.Type, oe.Ref, oe.Src, oe.Note)
	return oe
}
