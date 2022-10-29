package events

import "github.com/digitallyserviced/coolors/coolor/util"

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

	AnimationInit
  AnimationPlaying
  AnimationPaused
  AnimationIdle
	AnimationFinished
	AnimationLooped
	AnimationNext
	AnimationPrevious
	AnimationSet
	AnimationUpdate
	AnimationCanceled
  AnimationDone

  PluginEvents

	AllEvents ObservableEventType = SelectedEvent | ColorSeentEvent | ColorEvent | ColorSelectedEvent | ColorSelectionEvent | ChangedEvent | StatusEvent | InputEvent | DrawEvent | EditEvent | CancelableEvent | ExclusiveEvent | PaletteColorModifiedEvent | PaletteColorRemovedEvent | PaletteMetaUpdatedEvent | PaletteCreatedEvent | PaletteSavedEvent | PaletteColorSelectedEvent | PaletteColorSelectionEvent | AnimationInit | AnimationPlaying | AnimationPaused | AnimationDone| AnimationIdle | AnimationNext | AnimationSet | AnimationUpdate | AnimationFinished | AnimationLooped | AnimationCanceled | AnimationPrevious | PluginEvents
)

var observableEventTypes = []EnumName{
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
	{uint32(AnimationInit), "AnimationInit"},
  {uint32(AnimationPlaying), "AnimationPlaying"},
  {uint32(AnimationPaused), "AnimationPaused"},
  {uint32(AnimationIdle), "AnimationIdle"},
	{uint32(AnimationNext), "AnimationNext"},
	{uint32(AnimationSet), "AnimationSet"},
	{uint32(AnimationUpdate), "AnimationUpdate"},
	{uint32(AnimationFinished), "AnimationFinished"},
	{uint32(AnimationLooped), "AnimationLooped"},
	{uint32(AnimationDone), "AnimationDone"},
	{uint32(AnimationCanceled), "AnimationCanceled"},
	{uint32(AnimationPrevious), "AnimationPrevious"},
	{uint32(PluginEvents), "PluginEvents"},
	{uint32(AllEvents), "AllEvents"},
}

func init() {
}

func (a ObservableEventType) Is(b ObservableEventType) bool {
	return util.BitAnd(a, b)
}
func (v ObservableEventType) String() string {
	return EnumString(uint32(v), observableEventTypes, false)
}
func (v ObservableEventType) GoString() string {
	return EnumString(uint32(v), observableEventTypes, true)
}
