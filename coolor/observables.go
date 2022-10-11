package coolor

import (
	"fmt"
	// "unsafe"

	"github.com/gdamore/tcell/v2"
)

type ObservableEventType uint32

const (
	SelectedEvent = ObservableEventType(1 << iota)
	ColorSeentEvent
	ColorSelectedEvent
	ColorSelectionEvent
	ColorEvent

	ChangedEvent
	StatusEvent
	InputEvent
	DrawEvent
	EditEvent
	CancelableEvent
	ExclusiveEvent

	PaletteColorRemovedEvent
	PaletteColorSelectedEvent
	PaletteColorSelectionEvent
	PaletteColorModifiedEvent
	PaletteMetaUpdatedEvent
	PaletteCreatedEvent
	PaletteSavedEvent

	AnimationStarted
	AnimationIdle
	AnimationNext
	AnimationSet
	AnimationUpdate
	AnimationFinished
	AnimationCanceled
	AnimationPrevious

	AllEvents ObservableEventType = SelectedEvent | ColorSeentEvent | ColorEvent | ColorSelectedEvent | ColorSelectionEvent | ChangedEvent | StatusEvent | InputEvent | DrawEvent | EditEvent | CancelableEvent | ExclusiveEvent | PaletteColorModifiedEvent | PaletteColorRemovedEvent | PaletteMetaUpdatedEvent | PaletteCreatedEvent | PaletteSavedEvent | PaletteColorSelectedEvent | PaletteColorSelectionEvent | AnimationStarted | AnimationIdle | AnimationNext | AnimationSet | AnimationUpdate | AnimationFinished | AnimationCanceled | AnimationPrevious
)

var observableEventTypes = []enumName{
	{uint32(SelectedEvent), "SelectedEvent"},
	{uint32(ColorSeentEvent), "ColorSeentEvent"},
	{uint32(ChangedEvent), "ChangedEvent"},
	{uint32(StatusEvent), "StatusEvent"},
	{uint32(InputEvent), "InputEvent"},
	{uint32(DrawEvent), "DrawEvent"},
	{uint32(ColorEvent), "ColorEvent"},
	{uint32(EditEvent), "EditEvent"},
	{uint32(CancelableEvent), "CancelableEvent"},
	{uint32(ExclusiveEvent), "ExclusiveEvent"},
	{uint32(PaletteColorRemovedEvent), "PaletteColorRemovedEvent"},
	{uint32(PaletteColorSelectedEvent), "PaletteColorSelectedEvent"},
	{uint32(PaletteColorSelectedEvent), "PaletteColorSelectionEvent"},
	{uint32(PaletteColorModifiedEvent), "PaletteColorModifiedEvent"},
	{uint32(PaletteMetaUpdatedEvent), "PaletteMetaUpdatedEvent"},
	{uint32(PaletteCreatedEvent), "PaletteCreatedEvent"},
	{uint32(PaletteSavedEvent), "PaletteSavedEvent"},
	{uint32(AnimationStarted), "AnimationStarted"},
	{uint32(AnimationIdle), "AnimationIdle"},
	{uint32(AnimationNext), "AnimationNext"},
	{uint32(AnimationSet), "AnimationSet"},
	{uint32(AnimationUpdate), "AnimationUpdate"},
	{uint32(AnimationFinished), "AnimationFinished"},
	{uint32(AnimationCanceled), "AnimationCanceled"},
	{uint32(AnimationPrevious), "AnimationPrevious"},
	{uint32(AllEvents), "AllEvents"},
}

func init() {
}

func (v ObservableEventType) String() string {
	return enumString(uint32(v), observableEventTypes, false)
}
func (v ObservableEventType) GoString() string {
	return enumString(uint32(v), observableEventTypes, true)
}

type OnCoolorColorSelected interface {
	SelectedEvent(ev ObservableEvent) bool
}

type SelectedEventHandler interface {
	tcell.EventHandler
	OnCoolorColorSelected
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

	eventObserver struct {
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

func NewEventObserver(name string) *eventObserver {
	// unsafe.Pointer
	eo := &eventObserver{
		observerName: name,
		SequencedItem: SequencedItem{
			seqNo: 0,
		},
	}
	return eo
}

func (o *ObservableEvent) String() string {
	return fmt.Sprintf("anon  %s %s %s %v %v %d", o.EventTime, o.Type.String(), o.Note, o.Ref, o.Src, o.seqNo)
}

func (o *eventObserver) Name() string {
	return fmt.Sprintf("%s @ #%d", o.observerName, o.seqNo)
	// fmt.Printf("*** Observer %d received: \n", o.observerName)
	// return true
}

func (o *eventObserver) HandleEvent(e ObservableEvent) bool {
	// fmt.Printf("*** Observer %d received: \n", o.observerName)
	return true
}
type AnonymousObserver struct {
	*eventObserver
	Callbacks []Observer
}

type AnonymousHandler struct {
	Callback func(e ObservableEvent) bool
}

func (ah *AnonymousHandler) Name() string {
	return fmt.Sprintf("%s @ #%d", "Anon handler", 1)
}

func (ah *AnonymousHandler) HandleEvent(e ObservableEvent) bool {
	// fmt.Println("anon", e)
	ah.Callback(e)
	return true
}

func NewAnonymousHandlerFunc(f func(e ObservableEvent) bool) *AnonymousHandler {
	ahh := &AnonymousHandler{
		Callback: f,
	}
	return ahh
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

func NewAnonymousHandler(callbacks []Observer) *AnonymousObserver {
	ah := &AnonymousObserver{
		eventObserver: NewEventObserver("anon"),
	}
	return ah
}

