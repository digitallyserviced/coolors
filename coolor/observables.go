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

  AllEvents ObservableEventType = SelectedEvent | ColorSeentEvent| ColorEvent| ColorSelectedEvent | ColorSelectionEvent | ChangedEvent | StatusEvent | InputEvent | DrawEvent | EditEvent | CancelableEvent | ExclusiveEvent | PaletteColorModifiedEvent | PaletteColorRemovedEvent | PaletteMetaUpdatedEvent | PaletteCreatedEvent | PaletteSavedEvent | PaletteColorSelectedEvent | PaletteColorSelectionEvent
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
	{uint32(AllEvents), "AllEvents"},
}

func init(){
}

func (v ObservableEventType) String() string   { return enumString(uint32(v), observableEventTypes, false) }
func (v ObservableEventType) GoString() string { return enumString(uint32(v), observableEventTypes, true) }

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
	return fmt.Sprintf("%d", o.Type)
}

func (o *eventObserver) HandleEvent(e ObservableEvent) bool {
	// fmt.Printf("*** Observer %d received: \n", o.observerName)
	return true
}
