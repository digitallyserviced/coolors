package events

import "github.com/digitallyserviced/coolors/coolor/util"

type ObservableEventType uint64

const (
	SelectedEvent = ObservableEventType(1 << iota)
	ColorSeentEvent
	ColorSelectedEvent
	ColorSelectionEvent
	ColorEvent
	ColorFavoriteEvent
	ColorUnfavoriteEvent

	PrimaryEvent
	SecondaryEvent
	CancelledEvent

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

  PromptedEvents ObservableEventType = PrimaryEvent | SecondaryEvent | CancelledEvent

	AllEvents ObservableEventType = SelectedEvent | ColorSeentEvent | ColorEvent | ColorSelectedEvent | ColorSelectionEvent | PrimaryEvent | SecondaryEvent | CancelledEvent | ChangedEvent | StatusEvent | InputEvent | DrawEvent |ColorFavoriteEvent | 	ColorUnfavoriteEvent| EditEvent | CancelableEvent | ExclusiveEvent | PaletteColorModifiedEvent | PaletteColorRemovedEvent | PaletteMetaUpdatedEvent | PaletteCreatedEvent | PaletteSavedEvent | PaletteColorSelectedEvent | PaletteColorSelectionEvent | AnimationInit | AnimationPlaying | AnimationPaused | AnimationDone | AnimationIdle | AnimationNext | AnimationSet | AnimationUpdate | AnimationFinished | AnimationLooped | AnimationCanceled | AnimationPrevious | PluginEvents
)

var observableEventTypes = []EnumName{
	{uint64(SelectedEvent), "SelectedEvent"},
	{uint64(ColorSeentEvent), "ColorSeentEvent"},
	{uint64(ChangedEvent), "ChangedEvent"},
	{uint64(StatusEvent), "StatusEvent"},
	{uint64(InputEvent), "InputEvent"},
	{uint64(DrawEvent), "DrawEvent"},
	{uint64(ColorEvent), "ColorEvent"},
{uint64(ColorFavoriteEvent), "ColorFavoriteEvent"},
{uint64(ColorUnfavoriteEvent), "ColorUnfavoriteEvent"},
	{uint64(EditEvent), "EditEvent"},
	{uint64(CancelableEvent), "CancelableEvent"},
	{uint64(ExclusiveEvent), "ExclusiveEvent"},
	{uint64(PrimaryEvent), "PrimaryEvent"},
	{uint64(SecondaryEvent), "SecondaryEvent"},
	{uint64(CancelledEvent), "CancelledEvent"},
	{uint64(PaletteColorRemovedEvent), "PaletteColorRemovedEvent"},
	{uint64(PaletteColorSelectedEvent), "PaletteColorSelectedEvent"},
	{uint64(PaletteColorSelectedEvent), "PaletteColorSelectionEvent"},
	{uint64(PaletteColorModifiedEvent), "PaletteColorModifiedEvent"},
	{uint64(PaletteMetaUpdatedEvent), "PaletteMetaUpdatedEvent"},
	{uint64(PaletteCreatedEvent), "PaletteCreatedEvent"},
	{uint64(PaletteSavedEvent), "PaletteSavedEvent"},
	{uint64(AnimationInit), "AnimationInit"},
	{uint64(AnimationPlaying), "AnimationPlaying"},
	{uint64(AnimationPaused), "AnimationPaused"},
	{uint64(AnimationIdle), "AnimationIdle"},
	{uint64(AnimationNext), "AnimationNext"},
	{uint64(AnimationSet), "AnimationSet"},
	{uint64(AnimationUpdate), "AnimationUpdate"},
	{uint64(AnimationFinished), "AnimationFinished"},
	{uint64(AnimationLooped), "AnimationLooped"},
	{uint64(AnimationDone), "AnimationDone"},
	{uint64(AnimationCanceled), "AnimationCanceled"},
	{uint64(AnimationPrevious), "AnimationPrevious"},
	{uint64(PluginEvents), "PluginEvents"},
	{uint64(AllEvents), "AllEvents"},
}

func init() {
}

func (a ObservableEventType) Is(b ObservableEventType) bool {
	return util.BitAnd(a, b)
}
func (v ObservableEventType) String() string {
	return EnumString(uint64(v), observableEventTypes, false)
}
func (v ObservableEventType) GoString() string {
	return EnumString(uint64(v), observableEventTypes, true)
}
