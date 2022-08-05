package coolor

import (
	"fmt"

	"github.com/digitallyserviced/coolors/status"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/structs"
)

type CoolorColorActionFlag int
type CoolorColorActionSet CoolorColorActionFlag
type ActorFlag struct {
	actionFlag CoolorColorActionFlag
	name       string
}

const (
	NilFlag CoolorColorActionFlag = 1 << iota
	AddColorFlag
	RemoveColorFlag
	LockColorFlag
	DuplicateColor
	MixColorFlag
	InfoColorFlag
	ColorShadesFlag
	ColorContrastsFlag
	ColorGradientFlag

	MainPaletteActionsFlag = RemoveColorFlag | LockColorFlag | DuplicateColor | MixColorFlag | InfoColorFlag | ColorShadesFlag | ColorContrastsFlag | ColorGradientFlag
)

type CoolorColorActionFunctions struct {
	Activate func(cca *CoolorColorActor, cc *CoolorColor) bool
	Before   func(cca *CoolorColorActor, cc *CoolorColor) bool
	Every    func(cca *CoolorColorActor, cc *CoolorColor) bool
	Finalize func(cca *CoolorColorActor, cc *CoolorColor)
	Always   func(cca *CoolorColorActor)
	Actions  func(cca *CoolorColorActor) CoolorColorActionFlag
	Cancel   func(cca *CoolorColorActor) bool
}

type CoolorColorAction struct {
	flag       CoolorColorActionFlag
	icon, name string
	// options   *map[string]interface{}
	*CoolorColorActionFunctions
}
type ActorFunction func() *CoolorColorActionFunctions

type CoolorColorActionFunction func(name, icon string, actionFlag CoolorColorActionFlag) *CoolorColorAction

type CoolorColorActor struct {
	color     *CoolorColor
	activated bool
	menu      *CoolorToolMenu
	*CoolorColorAction
	actionSet CoolorColorActionFlag
	actor     *CoolorColorAct
}

type CoolorColorActors []*CoolorColorAct
type CoolorColorAct struct {
	name, icon string
	flag       CoolorColorActionFlag
	actor      ActorFunction
	actionSet  CoolorColorActionFlag
}

type ActorGroup []*CoolorColorActor
type ActorSet CoolorColorActionFlag

var (
	ActionOptions *structs.MapDataStore
	actors        CoolorColorActors
	MixColor      CoolorColorAct
	RemoveColor   CoolorColorAct
	AddColor      CoolorColorAct
	LockColor     CoolorColorAct

	acts   map[string]*CoolorColorActor
	groups map[string]ActorGroup
	sets   map[CoolorColorActionFlag]CoolorColorActionFlag
)

func init() {
	ActionOptions = structs.NewMapData()
	AddColor = CoolorColorAct{"add", "", AddColorFlag, addFunc, MixColorFlag}
	RemoveColor = CoolorColorAct{"remove", "", RemoveColorFlag, removeFunc, NilFlag}
	LockColor = CoolorColorAct{"lock", "", LockColorFlag, lockFunc, NilFlag}
	MixColor = CoolorColorAct{"mix", "ﭚ", MixColorFlag, mixFunc, NilFlag}
	actors = append(actors, &RemoveColor)
	actors = append(actors, &MixColor)
	actors = append(actors, &LockColor)
	// actors = append(actors, &MixColor)
	ActionsInit()
}

func ActionsInit() {
	acts = make(map[string]*CoolorColorActor)
	sets = make(map[CoolorColorActionFlag]CoolorColorActionFlag)
	// sets["main"] = MixColorFlag
	for _, v := range actors {
		if f, ok := sets[v.actionSet]; ok {
			sets[v.flag] = f | v.flag
		} else {
			sets[v.flag] = v.flag
		}
		acts[v.name] = SetupCoolorActor(v)
	}
}

func SetupCoolorActor(actor *CoolorColorAct) *CoolorColorActor {
	ncca := NewCoolorActor(actor.name, actor.icon, actor.flag, actor.actor, actor.actionSet)
	ncca.actor = actor
	return ncca
}
func NewCoolorActor(name, icon string, actionFlag CoolorColorActionFlag, f ActorFunction, actionSet CoolorColorActionFlag) *CoolorColorActor {
	cca := &CoolorColorAction{
		flag: actionFlag,
		icon: icon,
		name: name,
		// options:                    &map[string]interface{}{},
	}
	actor := &CoolorColorActor{
		color:             nil,
		activated:         false,
		actionSet:         actionSet,
		CoolorColorAction: cca,
	}
	cca.CoolorColorActionFunctions = f()

	return actor
}

func (cca *CoolorColorActor) Eq(s string, d string) bool {
	if v, has := cca.Has(s); !has || v != d {
		return false
	}
	return true
}

func (cca *CoolorColorActor) SetValue(s string, d interface{}) bool {
	ActionOptions.SetValue(cca.makeKey(s), d)
	return true
}

func (cca *CoolorColorActor) ClearColor() bool {
	return cca.SetValue("color", "")
}

func (cca *CoolorColorActor) SetColor(cc *CoolorColor) bool {
	return cca.SetValue("color", cc.Html())
}

func (cca *CoolorColorActor) GetColor() *CoolorColor {
	if c, has := cca.Has("color"); !has || c == "" {
		return nil
	} else {
		return NewCoolorColor(c.(string))
	}
}

func (cca *CoolorColorActor) TakeColor() *CoolorColor {
	if c, has := cca.Has("color"); !has || c == "" {
		return nil
	} else {
		cca.ClearColor()
		return NewCoolorColor(c.(string))
	}
}

func (cca *CoolorColorActor) Off(s string) bool {
	return cca.SetValue(s, false)
}

func (cca *CoolorColorActor) On(s string) bool {
	return cca.SetValue(s, true)
}

func (cca *CoolorColorActor) Has(s string) (interface{}, bool) {
	d, ok := ActionOptions.Value(cca.makeKey(s)) // .BoolVal(cca.makeKey(s))
	return d, ok
}

func (cca *CoolorColorActor) Is(s string) bool {
	if v, has := cca.Has(s); !has {
		return false
	} else {
		return v.(bool)
	}
}

func (cca *CoolorColorActor) MustOff(s string) bool {
	if v, has := cca.Has(s); !has || v.(bool) == false {
		return false
	}
	return cca.Off(s)
}
func (cca *CoolorColorActor) MustOn(s string) bool {
	if v, has := cca.Has(s); !has || v.(bool) == true {
		return false
	}
	return cca.On(s)
}

func (cca *CoolorColorActor) Dectivated() bool {
	return cca.Off("activated")
}

func (cca *CoolorColorActor) Activated() bool {
	return cca.On("activated")
}

func (cca *CoolorColorActor) IsActivated() bool {
	return cca.Is("activated")
}

func makeKey(typ, s string) string {
	return fmt.Sprintf("%s.%s", typ, s)
}
func (cca *CoolorColorActor) makeKey(s string) string {
	return makeKey(cca.actor.name, s)
}

func (cca *CoolorColorActor) Cancel() {
	cca.CoolorColorActionFunctions.Cancel(cca)
	cca.Always()
}
func (cca *CoolorColorActor) Actions() CoolorColorActionFlag {
	if cca == nil {
		return MainPaletteActionsFlag
	}
	return cca.CoolorColorActionFunctions.Actions(cca)
}
func (cca *CoolorColorActor) Always() {
	cca.CoolorColorActionFunctions.Always(cca)
}
func (cca *CoolorColorActor) Finalize(cc *CoolorColor) {
	cca.CoolorColorActionFunctions.Finalize(cca, cc)
	cca.Always()
}
func (cca *CoolorColorActor) Every(cc *CoolorColor) bool {
	return cca.CoolorColorActionFunctions.Every(cca, cc)
}
func (cca *CoolorColorActor) Before(cp *CoolorPalette, cc *CoolorColor) bool {
	return cca.CoolorColorActionFunctions.Before(cca, cc)
}
func (cca *CoolorColorActor) Activate(cc *CoolorColor) bool {
	return cca.CoolorColorActionFunctions.Activate(cca, cc)
}
func removeFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				cca.Activated()
				cca.SetColor(cc)
			} else {
				cca.Finalize(cc)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.CoolorColorAction.icon = ""
				cca.CoolorColorAction.name = "confirm"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.CoolorColorAction.icon = ""
				cca.CoolorColorAction.name = "confirm"
			}
			if pc := cca.GetColor(); pc != nil {
				if pc.Html() != cc.Html() {
					cca.Cancel()
				}
			} else {
				cca.Cancel()
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if pc := cca.TakeColor(); pc != nil {
				if pc.Html() == cc.Html() {
					col, _ := MainC.palette.GetSelected()
					col.Remove()
					MainC.palette.NavSelection(1)
					cca.Always()
					status.NewStatusUpdate("action_str", fmt.Sprintf("Removed %s", cc.TVPreview()))
				}
			} else {
				cca.Cancel()
			}
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.CoolorColorAction.icon = cca.actor.icon
			cca.CoolorColorAction.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				return RemoveColorFlag
			}
			return -1
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.Dectivated()
			cca.CoolorColorAction.icon = cca.actor.icon
			cca.CoolorColorAction.name = cca.actor.name
			return true
		},
	}
	return ccaf
}
func mixFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cca.IsActivated() {
				if cca.Is("selection") {
					cca.Finalize(cc.Clone())
					return true
				} else {
					cCol := cca.TakeColor()
					cca.On("selection")
					MainC.NewMixer(cCol, cc)
					return true
				}
			} else {
				cca.Activated()
				cca.SetColor(cc)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cca.IsActivated() && cca.Is("selection") {
				cca.CoolorColorAction.icon = ""
				cca.CoolorColorAction.name = "take"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			dump.P(ActionOptions.Data())
			if !cca.IsActivated() {
				return true
			}
			if cca.Is("selection") {
				status.NewStatusUpdate("color", cc.TVPreview())
				status.NewStatusUpdate("action_str", fmt.Sprintf("= %s", cc.TVPreview()))
			} else {
				cCol := cca.GetColor()
				if cCol == nil {
					return true
				}
				status.NewStatusUpdate("action_str", fmt.Sprintf("= %s -> [black:yellow:b] ﭚ [-:-:-] -> %s", cCol.TVPreview(), cc.TVPreview()))
			}
			// cca.menu.updateState()
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if MainC.palette != nil {
				MainC.palette.AddCoolorColor(cc.Unstatic())
				cca.Always()
			}
			if MainC.mixer != nil {
				MainC.mixer.Clear()
			}
			status.NewStatusUpdate("action_str", fmt.Sprintf("%s", "Finaliz'd Mix"))
			cca.Dectivated()
			cca.Off("selection")
			// status.NewStatusUpdate("action_str", fmt.Sprintf("%s", "End Mix"))
			cca.CoolorColorAction.icon = cca.actor.icon
			cca.CoolorColorAction.name = cca.actor.name
			MainC.pages.SwitchToPage("palette")
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.Off("selection")
			// status.NewStatusUpdate("action_str", fmt.Sprintf("%s", "End Mix"))
			cca.CoolorColorAction.icon = cca.actor.icon
			cca.CoolorColorAction.name = cca.actor.name
			MainC.pages.SwitchToPage("palette")
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				if cca.Is("selection") {
					return MixColorFlag
				} else {
					return MixColorFlag
				}
			}
			return MixColorFlag | RemoveColorFlag
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.TakeColor()
			cca.Dectivated()
			cca.Off("selection")
			status.NewStatusUpdate("action_str", fmt.Sprintf("%s", "Cancel Mix"))
			return true
		},
	}
	return ccaf
}

func addFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			MainC.palette.AddCoolorColor(cc.Unstatic())
			cca.Finalize(cc.Unstatic())
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			cca.Off("activated")
		},
		Always: func(cca *CoolorColorActor) {
			cca.CoolorColorAction.icon = cca.actor.icon
			cca.CoolorColorAction.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return -1
		},
		Cancel: func(cca *CoolorColorActor) bool {
			return true
		},
	}
	return ccaf
}

func lockFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cc == nil {
				return false
			}
			cc.SetLocked(!cc.GetLocked())
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cc == nil {
				return false
			}
			dump.P("before", cc.GetLocked())
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			ncc := MainC.palette.colors[MainC.palette.selectedIdx]
			// status.NewStatusUpdate("action_str", fmt.Sprintf("= %s -> [black:yellow:b] ﭚ [-:-:-] -> ", ncc))
			if ncc == nil {
				return false
			}
			if ncc.GetLocked() {
				cca.CoolorColorAction.icon = ""
				cca.CoolorColorAction.name = "unlock"
			} else {
				cca.CoolorColorAction.icon = ""
				cca.CoolorColorAction.name = "lock"
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
		},
		Always: func(cca *CoolorColorActor) {
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return MainPaletteActionsFlag
		},
		Cancel: func(cca *CoolorColorActor) bool {
			return true
		},
	}
	return ccaf
}

func templateFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
		},
		Always: func(cca *CoolorColorActor) {
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return -1
		},
		Cancel: func(cca *CoolorColorActor) bool {
			return true
		},
	}
	return ccaf
}

// vim: ts=2 sw=2 et ft=go
