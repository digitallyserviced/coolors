package coolor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/digitallyserviced/coolors/status"
	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"
	// "github.com/lucasb-eyer/go-colorful"
)

type CoolColorMod struct {
	*CoolorColor
	// *colorful.Color
}
type ColorModifier struct {
	ChannelMod *ChannelMod
	*ChannelModOptions
	name string
}

type ColorModAction struct {
	Result   CoolColor
	Action   ColorModFunction
	Function string
	Argument float64
}

type ColorModActions []*ColorModAction

type Gradiater interface {
	Above() []CoolColor
	Below() []CoolColor
	At(value float64) CoolColor
	Set(value float64)
	Incr(value float64)
	Decr(value float64)
	GetChannelValue(cc CoolColor) float64
	GetCurrentChannelValue() float64
}

type ColorModder interface {
	Incr(float64)
	Decr(float64)
	Set(float64)
	Nop(bool)
}

type ColorMod struct {
	last *ColorModAction
	orig *CoolorColor
	*ColorModifier
	current CoolorColor
	history ColorModActions
	ring    CoolorColors
}

type ColorModFunction func(float64) *ColorModAction

// len(ColorModNames) 🮋🮒🮑🮐🮆🮔🮕🮖🮗🮟🮞🮝🮜🮘🮙🮚🮱🮴🮽🮾🮿🯄

var ColorModActionStrings = map[string]string{
	"set": "[blue:black:-][black:blue:-]=[blue:black:-][-:-:-]",
	"inc": "[green:black:-][black:green:-]+[green:black:-][-:-:-]",
	"dec": "[red:black:-][black:red:-]-[red:black:-][-:-:-]",
}

func (ccm *CoolColorMod) GetCC() *CoolorColor {
	return ccm.CoolorColor.GetCC()
}

func (cma *ColorModAction) Summary() string {
	return fmt.Sprintf("%s %0.2f %s", ColorModActionStrings[cma.Function], cma.Argument, cma.Result.GetCC().TVPreview())
}

func (cma *ColorModAction) String() string {
	return ColorModActionStrings[cma.Function]
}
func (cmlog ColorModActions) String() string {
	summActions := lo.Map(cmlog, func(x *ColorModAction, _ int) string {
		if x != nil {
			return x.Summary()
		}
		return ""
	})

	_ = summActions
	return strings.Join(summActions, "\n")
}

func NewColorModAction(name string, f ColorModFunction, arg float64, result CoolColor) *ColorModAction {
	return &ColorModAction{
		Function: name,
		Action:   f,
		Argument: arg,
		Result:   result,
	}
}

func (cm *ColorMod) Log(action *ColorModAction) **ColorModAction {
	cm.history = lo.Subset(cm.history, -19, 19)
	cm.history = append(cm.history, action)
	cm.last = action
  status.NewStatusUpdateWithTimeout("action_str", action.Summary(), 3 * time.Second)
	dump.P(cm.history.String())
	return &cm.last
}

func (cm *ColorMod) Set(value float64) *ColorModAction {
	value = cm.ChannelMod.Max() * value
	c, _ := cm.ChannelMod.SetValue(&cm.current, value) // .GetCC()
	cm.current = *c.GetCC()
	cm.updateState(false)
	return *cm.Log(NewColorModAction("set", cm.Set, value, cm.current.Clone()))
}

func (cm *ColorMod) Incr(value float64) *ColorModAction {
	if value == 0 {
		value = cm.increment
	}
	cnew, err := cm.ChannelMod.Mod(cm.current.GetCC(), math.Abs(cm.increment))
	cm.current = *cnew.GetCC()
	// dump.P(value,cm.current.TerminalPreview())
	_ = err
	cm.updateState(false)
	return *cm.Log(NewColorModAction("inc", cm.Incr, value, cm.current.Clone()))
}

func (cm *ColorMod) Decr(value float64) *ColorModAction {
	if value == 0 {
		value = cm.increment
	}
	cnew, err := cm.ChannelMod.Mod(cm.current.GetCC(), -math.Abs(cm.increment))
	cm.current = *cnew.GetCC()
	// dump.P(value,cm.current.TerminalPreview())
	_ = err
	cm.updateState(false)
	return *cm.Log(NewColorModAction("dec", cm.Decr, value, cm.current.Clone()))
}

func NewChannelMod(mod *ChannelModOptions, chm ChannelModifier) *ColorMod {
	cc := NewRandomCoolorColor()
	cm := &ColorMod{
		ColorModifier: &ColorModifier{},
	}
	cm.history = make(ColorModActions, 0)
	cm.last = NewColorModAction("nop", cm.Set, 0.0, nil)
	cm.orig = cc
	cm.current = *cc.Clone()
	cm.name = mod.name
	cm.ChannelMod = chm(mod)
	cm.ChannelModOptions = mod
	cm.updateState(true)
	return cm
}
func (cm *ColorModifier) SetSize(size float64) {
	cm.size = size
	cm.updateState(false)
}

func (cm *ColorModifier) updateState(noUpdateColor bool) {
}

func (cm *ColorMod) Pop(cc CoolColor) {
	cm.ring = append(cm.ring, cm.current.GetCC())
	cm.current = *cc.GetCC()
}

func (cm *ColorMod) Push(cc CoolColor) {
	cm.ring = append(cm.ring, cm.current.GetCC())
	cm.current = *cc.GetCC()
}

func (cm *ColorMod) Next() CoolColor {
	cnew := cm.ChannelMod.ModPct(cm.current.Clone(), cm.increment)
	if cnew.GetCC().Html() == cm.current.GetCC().Html() {
		return cnew.GetCC()
	}
	if len(cm.ring) < int(cm.size) {
		cm.Push(cnew.GetCC())
	}
	return cnew.GetCC()
}

func (cm *ColorMod) makeGrad(above bool) []CoolorColor {
	num := cm.size
	if above {
		num = (-math.Abs(cm.size))
	}
	fmt.Println(cm.size, cm.increment)
	colors := cm.ChannelMod.RangePct(&cm.current, cm.increment, num)
	return colors
}

func (cm *ColorMod) SetColor(cc CoolColor) {
	cm.orig = cc.GetCC()
	cm.current = *cm.orig.Clone()
	cm.updateState(false)
}

func (cm *ColorMod) GetCurrentChannelValue() float64 {
	return cm.GetChannelValue(&cm.current)
}

func (cm *ColorMod) GetStatus() (string, string) {
	last := ""
	status := ""
	if cm.last != nil {
		last = cm.last.String()
	}
	status = cm.ChannelMod.FormatValue(&cm.current)
	return status, last
}

func (cm *ColorModifier) String() string {
	return fmt.Sprintf(
		"%s +/- %0.2f ",
		cm.name,
		cm.increment,
	)
}

func (cm *ColorMod) GetChannelValue(cc CoolColor) float64 {
	value := cm.ChannelMod.GetValue(&cm.current)
	return float64(value)
}

func clamped(val, min, max float64) (float64, bool) {
	clampd := val > max
	
	if val < min {
		clampd = true
	}
	return clamp(val, min, max), clampd
}
func clamp(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

// vim: ts=2 sw=2 et ft=go
