package coolor

import (
	"fmt"
	"strings"
	"time"

	"github.com/digitallyserviced/tview"

	"github.com/digitallyserviced/coolors/coolor/lister"
	"github.com/digitallyserviced/coolors/status"

	// "github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/structs"
)

type (
	CoolorColorActionFlag uint
	CoolorColorActionSet  CoolorColorActionFlag
	ActorFlag             struct {
		name       string
		actionFlag CoolorColorActionFlag
	}
)

const (
	NilFlag CoolorColorActionFlag = 1 << iota
	AddColorFlag
	RemoveColorFlag
	LockColorFlag
	DuplicateColorFlag
	SwapColorFlag
	MixColorFlag
	RandomizeColorFlag
	InfoColorFlag
	TagColorFlag
  FavoriteColorFlag
	ShadeColorFlag
	ColorContrastsFlag
	ColorGradientFlag

	MainPaletteActionsFlag = RemoveColorFlag | LockColorFlag | DuplicateColorFlag | MixColorFlag | InfoColorFlag | ShadeColorFlag | ColorContrastsFlag | ColorGradientFlag | RandomizeColorFlag | SwapColorFlag | TagColorFlag | FavoriteColorFlag
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
	*CoolorColorActionFunctions
	icon string
	name string
	flag CoolorColorActionFlag
}
type ActorFunction func() *CoolorColorActionFunctions

type CoolorColorActionFunction func(name, icon string, actionFlag CoolorColorActionFlag) *CoolorColorAction

type CoolorColorActor struct {
	color *CoolorColor
	menu  *CoolorToolMenu
	*CoolorColorAction
	actor     *CoolorColorAct
	actionSet CoolorColorActionFlag
	activated bool
}

type (
	CoolorColorActors []*CoolorColorAct
	CoolorColorAct    struct {
		actor     ActorFunction
		name      string
		icon      string
		flag      CoolorColorActionFlag
		actionSet CoolorColorActionFlag
	}
)

type (
	ActorGroup []*CoolorColorActor
	ActorSet   CoolorColorActionFlag
)

var (
	ActionOptions  *structs.Data
	actors         CoolorColorActors
	GradientColor  CoolorColorAct
	ShadeColor     CoolorColorAct
	SwapColor      CoolorColorAct
	RemoveColor    CoolorColorAct
	RandomizeColor CoolorColorAct
	FavoriteColor CoolorColorAct
	AddColor       CoolorColorAct
	TagColor       CoolorColorAct
	LockColor      CoolorColorAct

	actOrder []string
	acts     map[string]*CoolorColorActor
	groups   map[string]ActorGroup
	sets     map[CoolorColorActionFlag]CoolorColorActionFlag
)

func NewCoolorColorAct(
	name, icon string,
	flag CoolorColorActionFlag,
	ccaf ActorFunction,
	actionSet CoolorColorActionFlag,
	short rune,
) CoolorColorAct {
	cca := CoolorColorAct{
		name:      name,
		icon:      icon,
		flag:      flag,
		actor:     ccaf,
		actionSet: actionSet,
	}
	return cca
}

// ﱬ
// ₐₑₒₓₕₖₗₘₙₚₛₜ₀₁₂₃₄₅₆₇₈₉₊₋₌₍₎
// aeoxhklmnpst0123456789+-=()
func init() {
	ActionOptions = structs.NewData()
	//  
	AddColor = NewCoolorColorAct(
		"add",
		"",
		AddColorFlag,
		addFunc,
		AddColorFlag,
		'+',
	)
	TagColor = NewCoolorColorAct(
		"tag",
		"",
		TagColorFlag,
		tagFunc,
		NilFlag,
		't',
	)
	FavoriteColor = NewCoolorColorAct(
		"favorite",
		"",
		FavoriteColorFlag,
		favFunc,
		NilFlag,
		'h',
	)
	RemoveColor = NewCoolorColorAct(
		"remove",
		"",
		RemoveColorFlag,
		removeFunc,
		NilFlag,
		'-',
	)
	LockColor = NewCoolorColorAct(
		"lock",
		"",
		LockColorFlag,
		lockFunc,
		NilFlag,
		'l',
	)
	GradientColor = NewCoolorColorAct(
		"mix",
		"",
		MixColorFlag,
		mixFunc,
		NilFlag,
		'm',
	)
	//         ﭚ
	ShadeColor = NewCoolorColorAct(
		"shades",
		"",
		ShadeColorFlag,
		shadeFunc,
		NilFlag,
		'h',
	)
	SwapColor = NewCoolorColorAct(
		"swap",
		"",
		SwapColorFlag,
		swapFunc,
		NilFlag,
		'=',
	)
	RandomizeColor = NewCoolorColorAct(
		"randomize",
		"",
		RandomizeColorFlag,
		randomizeFunc,
		NilFlag,
		'o',
	)
	// SwapColor = CoolorColorAct{"swap", "", SwapColorFlag, swapFunc, NilFlag}
	// RandomizeColor = CoolorColorAct{"randomize", "", RandomizeColorFlag, randomizeFunc, NilFlag}
	actOrder = []string{
    "favorite",
		"tag",
		"lock",
		"swap",
		"randomize",
		"mix",
		"shades",
		"remove",
	}
	actors = append(actors, &FavoriteColor)
	actors = append(actors, &RemoveColor)
	actors = append(actors, &RandomizeColor)
	actors = append(actors, &GradientColor)
	actors = append(actors, &ShadeColor)
	actors = append(actors, &SwapColor)
	actors = append(actors, &LockColor)
	actors = append(actors, &TagColor)
	ActionsInit()
}

func ActionsInit() {
	acts = make(map[string]*CoolorColorActor)
	sets = make(map[CoolorColorActionFlag]CoolorColorActionFlag)
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
	ncca := NewCoolorActor(
		actor.name,
		actor.icon,
		actor.flag,
		actor.actor,
		actor.actionSet,
	)
	ncca.actor = actor
	return ncca
}

func NewCoolorActor(
	name, icon string,
	actionFlag CoolorColorActionFlag,
	f ActorFunction,
	actionSet CoolorColorActionFlag,
) *CoolorColorActor {
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
    return  d, ok
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
	// MainC.app.Sync()
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
	// cca.Always()
	// MainC.app.Sync()
}

func (cca *CoolorColorActor) Every(cc *CoolorColor) bool {
	return cca.CoolorColorActionFunctions.Every(cca, cc)
}

func (cca *CoolorColorActor) Before(
	cp *CoolorPaletteMainView,
	cc *CoolorColor,
) bool {
	return cca.CoolorColorActionFunctions.Before(cca, cc)
}

func (cca *CoolorColorActor) Activate(cc *CoolorColor) bool {
	// doLog(errorx.WithStack(errorx.Newf("%s %s", cc.Html(), cca.name)))
	// var st []uintptr
	// _ = runtime.Stack(st, true)
	// doLog(st)
	// dump.P(runtime.Caller(1))
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
				return false
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			if pc := cca.GetColor(); pc != nil {
				if pc.Html() != cc.Html() {
					cca.Cancel()
					return false
				}
			} else {
				cca.Cancel()
				return false
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
					status.NewStatusUpdateWithTimeout(
						"action_str",
						fmt.Sprintf("Removed %s", cc.TVPreview()),
						0*time.Second,
					)
				}
			} else {
				cca.Cancel()
			}
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				return RemoveColorFlag
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.Dectivated()
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
			return true
		},
	}
	return ccaf
}

func shadeFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cca.IsActivated() {
				if cca.Is("selection") {
					cca.Finalize(cc.Clone())
					return false
				}
			} else {
				cca.Activated()
				cca.SetColor(cc)
				cCol := cca.TakeColor()
				cca.On("selection")
				MainC.NewShades(cCol)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cca.IsActivated() && cca.Is("selection") {
				cca.icon = ""
				cca.name = "select"
				return true
			}
			return false
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			}
			if cca.Is("selection") {
				//         ﭚ
				cCol := cca.GetColor()
				if cCol == nil || cc == nil {
					return false
				}
				status.NewStatusUpdate("color", cc.TVPreview())
				status.NewStatusUpdateWithTimeout(
					"action_str",
					fmt.Sprintf(
						"[black:yellow:b] %s %sed [-:-:-]  %s <-> %s",
						cca.icon,
						cca.name,
						cCol.TVPreview(),
						cc.TVPreview(),
					),
					time.Microsecond*7000,
				)
			} else {
				cCol := cca.GetColor()
				if cCol == nil {
					return false
				}
				// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("= %s -> [black:yellow:b] ﭚ [-:-:-] -> %s", cCol.TVPreview(), cc.TVPreview()))
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if MainC.palette != nil {
				MainC.palette.AddCoolorColor(cc.Unstatic())
				SeentColor("selected_random_shade", cc, cc.pallette)
			}
			status.NewStatusUpdateWithTimeout(
				"action_str",
				"Finaliz'd Shade",
				0*time.Second,
			)
			cca.Always()
		},
		Always: func(cca *CoolorColorActor) {
			if MainC.shades != nil {
				MainC.shades.Clear()
			}
			cca.Dectivated()
			cca.TakeColor()
			cca.Off("selection")
			// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("%s", "End Shade"))
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
			MainC.pages.SwitchToPage("palette")
			MainC.pages.RemovePage("shades")
			MainC.shades = nil
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				if cca.Is("selection") {
					return ShadeColorFlag
				} else {
					return ShadeColorFlag
				}
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.TakeColor()
			cca.Always()
			status.NewStatusUpdateWithTimeout(
				"action_str",
				"Canceled Shade",
				0*time.Second,
			)
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
				} else {
					cCol := cca.TakeColor()
					cca.On("selection")
					MainC.NewMixer(cCol, cc)
				}
			} else {
				cca.Activated()
				cca.SetColor(cc)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if cca.IsActivated() && cca.Is("selection") {
				cca.icon = ""
				cca.name = "select"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			}
			//         ﭚ
			if cca.Is("selection") {
				status.NewStatusUpdate("color", cc.TVPreview())
				status.NewStatusUpdateWithTimeout(
					"action_str",
					fmt.Sprintf("[black:yellow:b]  [-:-:-] = %s", cc.TVPreview()),
					0*time.Second,
				)
			} else {
				cCol := cca.GetColor()
				if cCol == nil {
					return false
				}
				status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("%s -> [black:yellow:b]  [-:-:-] -> %s", cCol.TVPreview(), cc.TVPreview()), 0*time.Second)
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if MainC.palette != nil {
				MainC.palette.AddCoolorColor(cc.Unstatic())
				SeentColor("selected_blended_gradient_color", cc, cc.pallette)
			}
			cCol := cca.GetColor()
			status.NewStatusUpdateWithTimeout(
				"action_str",
				fmt.Sprintf(
					"= %s -> [black:yellow:b]  [-:-:-] -> %s",
					cCol.TVPreview(),
					cc.TVPreview(),
				),
				0*time.Second,
			)
			cca.Always()
		},
		Always: func(cca *CoolorColorActor) {
			if MainC.mixer != nil {
				MainC.mixer.Clear()
			}
			cca.Dectivated()
			cca.TakeColor()
			cca.Off("selection")
			// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("%s", "End Mix"))
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
			MainC.pages.SwitchToPage("palette")
			MainC.pages.RemovePage("mixer")
			MainC.mixer = nil
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				if cca.Is("selection") {
					return MixColorFlag
				} else {
					return MixColorFlag
				}
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.TakeColor()
			cca.Always()
			status.NewStatusUpdateWithTimeout(
				"action_str",
				fmt.Sprintf("%s", "Canceled Mix"),
				0*time.Second,
			)
			return true
		},
	}
	return ccaf
}

func addFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			MainC.palette.AddCoolorColor(cc.Unstatic())
			SeentColor("add_color", cc, cc.pallette)
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
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return 0
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
			// color, idx := MainC.palette.GetSelected()
			// dump.P(color.selected, color.locked, idx)
			color, _ := MainC.palette.ToggleLockSelected()
			// dump.P(color.selected, color.locked, idx)
			cca.Finalize(color)
			return false
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			ncc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if ncc == nil {
				return false
			}
			if ncc.GetLocked() {
				cca.icon = ""
				cca.name = "unlock"
			} else {
				cca.icon = ""
				cca.name = "lock"
			}
			return false
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			ncc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if ncc == nil {
				return false
			}
			if ncc.GetLocked() {
				cca.icon = ""
				cca.name = "unlock"
			} else {
				cca.icon = ""
				cca.name = "lock"
			}
			return false
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			ncc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if ncc == nil {
				return
			}
			icon := ""
			name := "lock"
			if !ncc.GetLocked() {
				icon = ""
				name = "unlock"
			}
			status.NewStatusUpdateWithTimeout(
				"action_str",
				fmt.Sprintf(
					"[black:yellow:b] %s %sed [-:-:-] -> %s",
					icon,
					name,
					cc.TVPreview(),
				),
				0*time.Second,
			)
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

func swapFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				cca.Activated()
				MainC.palette.Each(func(ccc *CoolorColor, i int) {
					if cc.Html() == ccc.Html() {
						cca.SetValue("orig_pos", i)
						cca.SetColor(ccc)
					}
				})
			} else {
				cca.Finalize(cc)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			pal, has := cca.Has("palette")
			pos, haspos := cca.Has("orig_pos")
			_, _, _, _, _ = cca.Activated(), has, haspos, pos, pal
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			if pc := cca.GetColor(); pc != nil {
				if pc.Html() != cc.Html() {
					from, to := -1, -1
					MainC.palette.Each(func(ccc *CoolorColor, i int) {
						if cc.Html() == ccc.Html() {
							to = i
						}
						if pc.Html() == ccc.Html() {
							from = i
						}
					})
					if to != -1 && from != -1 {
						MainC.palette.Swap(to, from)
						cca.Finalize(cc)
						MainC.palette.ResetViews()
					}
				}
			} else {
				cca.Cancel()
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if pc := cca.TakeColor(); pc != nil {
				if pc.Html() != cc.Html() {
					cca.Always()
					status.NewStatusUpdateWithTimeout(
						"action_str",
						fmt.Sprintf(
							"[black:yellow:b] %s %sed [-:-:-]  %s <-> %s",
							cca.icon,
							cca.name,
							pc.TVPreview(),
							cc.TVPreview(),
						),
						0*time.Second,
					)
					// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("Moved %s", pc.TVPreview()))
				}
			} else {
				cca.Cancel()
			}
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.SetValue("orig_pos", nil)
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				return SwapColorFlag
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			val, has := cca.Has("orig_pos")
			if has && val != nil {
				if pc := cca.GetColor(); pc != nil {
					from, to := -1, val
					MainC.palette.Each(func(ccc *CoolorColor, i int) {
						if pc.Html() == ccc.Html() {
							from = i
						}
					})
					if to != -1 && from != -1 {
						MainC.palette.Swap(from, to.(int))
						MainC.palette.ResetViews()
					}
				}
			}
			cca.Always()
			return true
		},
	}
	return ccaf
}

func randomizeFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				cca.Activated()
				cca.SetColor(cc)
				return true
			} else {
				cca.Finalize(cc)
				return false
			}
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			if !cca.IsActivated() {
				return false
			} else {
				cca.icon = ""
				cca.name = "confirm"
			}
			if pc := cca.GetColor(); pc != nil {
				if pc.Html() != cc.Html() {
					cca.Cancel()
					return false
				}
			} else {
				cca.Cancel()
				return false
			}
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if pc := cca.TakeColor(); pc != nil {
				oldCol := NewCoolorColor(pc.Html())
				if pc.Html() == cc.Html() {
					col, _ := MainC.palette.GetSelected()
					col.Random()
					MainC.palette.NavSelection(0)
					// MainC.palette.NavSelection(-1)
					cca.Always()
					status.NewStatusUpdateWithTimeout(
						"action_str",
						fmt.Sprintf(
							"[black:yellow:b] %s %sed [-:-:-] %s -> %s",
							cca.icon,
							cca.name,
							oldCol.TVPreview(),
							cc.TVPreview(),
						),
						0*time.Second,
					)
				}
			} else {
				cca.Cancel()
			}
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.IsActivated() {
				return RandomizeColorFlag
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.Dectivated()
			cca.icon = cca.actor.icon
			cca.name = cca.actor.name
			return true
		},
	}
	return ccaf
}

//  

func tagFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			cca.Activated()
			cca.SetColor(cc)
			tag := cc.GetTag(0)
			mc := MainC
			if mc.floater == nil {
				items := GetTerminalColorsAnsiTags()
				mc.floater = NewSelectionFloater(
					" Terminal ANSI Color Tag",
					items.GetListItems,
					func(lis lister.ListItem, hdr *tview.TextView, ftr *tview.TextView) {
						ccol, _ := MainC.palette.GetSelected()
						ti := lis.(*TagItem)
						MainC.palette.Each(func(cc *CoolorColor, i int) {
							tags := MainC.palette.Colors[i].GetTags()
							for _, v := range tags {
								// dump.P(ti, v)
								if ti.MainText() == v.MainText() {
									cc.ClearTags()
								}
							}
						})
						ccol.SetTag(ti)
						cca.SetValue("tag", lis.(*TagItem))
						cca.Before(ccol.pallette, ccol)
						cca.Every(ccol)
						cca.Finalize(cc)
					},
					func(lis lister.ListItem, hdr *tview.TextView, ftr *tview.TextView) {
						ftr.SetText(lis.MainText())
					},
				)
			}
			if !mc.pages.HasPage("floater") {
				mc.pages.AddPage("floater", mc.floater.GetRoot(), true, false)
			}
			mc.pages.ShowPage("floater")
			AppModel.helpbar.SetTable("floater")
			if tag != nil {
				mc.floater.(*ListFloater).Lister.SetCurrent(tag)
			}
			return true
		},
		Before: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			tags := cc.GetTags()
			tagStrs := make([]string, 0)
			for _, v := range tags {
				tagStrs = append(tagStrs, v.MainText())
			}
			status.NewStatusUpdate(
				"tag",
				fmt.Sprintf(
					"[black:purple:b] %s %s(s) [-:-:-] %s",
					cca.icon,
					cca.name,
					strings.Join(tagStrs, ", "),
				),
			)
			if pc := cca.GetColor(); pc != nil {
				if pc.Html() != cc.Html() {
					cca.Cancel()
					return false
				}
			} else {
				cca.Cancel()
				return false
			}
			return true
		},
		Every: func(cca *CoolorColorActor, cc *CoolorColor) bool {
			tags := cc.GetTags()
			tagStrs := make([]string, 0)
			for _, v := range tags {
				tagStrs = append(tagStrs, v.MainText())
			}
			status.NewStatusUpdate(
				"tag",
				fmt.Sprintf(
					"[black:purple:b] %s %s(s) [-:-:-] %s",
					cca.icon,
					cca.name,
					strings.Join(tagStrs, ", "),
				),
			)
			return true
		},
		Finalize: func(cca *CoolorColorActor, cc *CoolorColor) {
			if tag, has := cca.Has("tag"); has {
				tg, ok := tag.(*TagItem)
				if ok {
					var tagg TagItem = TagItem(*tg)
					status.NewStatusUpdateWithTimeout(
						"action_str",
						fmt.Sprintf(
							"[black:purple:b] %s %sged [-:-:-] %s",
							cca.icon,
							cca.name,
							tagg.MainText(),
						),
						0*time.Second,
					)

				}
			}
			cca.Always()
			// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("[black:yellow:b] %s %sed [-:-:-] %s -> %s", cca.icon, cca.name, oldCol.TVPreview(), cc.TVPreview()))
		},
		Always: func(cca *CoolorColorActor) {
			cca.Dectivated()
			cca.TakeColor()
			cca.SetValue("tag", nil)
			mc := MainC
			name, page := mc.pages.GetFrontPage()
			if name == "floater" {
				mc.pages.HidePage("floater")
				mc.pages.RemovePage("floater")
				page.Blur()
				mc.floater = nil
			}
			MainC.menu.UpdateVisibleActors(0)
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			if cca.Activated() {
				return TagColorFlag
			}
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			cca.Always()
			return true
		},
	}
	return ccaf
}

func favFunc() *CoolorColorActionFunctions {
	ccaf := &CoolorColorActionFunctions{
		Activate: func(cca *CoolorColorActor, ncc *CoolorColor) bool {
			cc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if cc == nil {
				return false
			}
			// GetStore().MetaService.ToggleFavorite(cc)
			cca.Finalize(cc)
			return false
		},
		Before: func(cca *CoolorColorActor, ncc *CoolorColor) bool {
			cc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if cc == nil {
				return false
			}
      // doCallers()
			_, idx := GetStore().MetaService.FavoriteColors.Contains(cc)
			fav := idx >= 0
			cca.icon = IfElseStr(fav, " ", " ")
			cca.name = IfElseStr(fav, "unfavorite", "favorite")
			return false
		},
		Every: func(cca *CoolorColorActor, ncc *CoolorColor) bool {
			cc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if cc == nil {
				return false
			}
			_, idx := GetStore().MetaService.FavoriteColors.Contains(cc)
			fav := idx >= 0
			cca.icon = IfElseStr(fav, " ", " ")
			cca.name = IfElseStr(fav, "unfavorite", "favorite")
			return false
		},
		Finalize: func(cca *CoolorColorActor, ncc *CoolorColor) {
			cc := MainC.palette.Colors[MainC.palette.selectedIdx]
			if cc == nil {
				return 
			}
			_, idx := GetStore().MetaService.FavoriteColors.Contains(cc)
			fav := idx >= 0
			icon := IfElseStr(fav, " ", " ")
			name := IfElseStr(fav, "unfavorite", "favorite")
			status.NewStatusUpdateWithTimeout(
				"action_str",
				fmt.Sprintf(
					"[black:yellow:b] %s %sed [-:-:-] -> %s",
					icon,
					name,
					cc.TVPreview(),
				),
				0*time.Second,
			)
		},
		Always: func(cca *CoolorColorActor) {
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return 0
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
			// status.NewStatusUpdateWithTimeout("action_str", fmt.Sprintf("[black:yellow:b] %s %sed [-:-:-] %s -> %s", cca.icon, cca.name, oldCol.TVPreview(), cc.TVPreview()))
		},
		Always: func(cca *CoolorColorActor) {
		},
		Actions: func(cca *CoolorColorActor) CoolorColorActionFlag {
			return 0
		},
		Cancel: func(cca *CoolorColorActor) bool {
			return true
		},
	}
	return ccaf
}

// vim: ts=2 sw=2 et ft=go
