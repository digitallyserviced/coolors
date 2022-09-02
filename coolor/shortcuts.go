package coolor

import (
	"fmt"

	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"
)

var SubShortcuts *ScriptShortcuts
var SuperShortcuts *ScriptShortcuts

type ScriptShortcuts struct {
	txtScriptMap               map[rune]ScriptShortcut
	used, origScript, origText []rune
	wrapBegin, wrapEnd         string
}

func (sss *ScriptShortcuts) Give(short ScriptShortcut) {
	sss.used = lo.Without[rune](sss.used, short.text)
}

func (sss *ScriptShortcuts) Take(txt rune) ScriptShortcut {
	dump.P(txt, sss.used, sss.origText)
	if lo.Contains[rune](sss.used, txt) {
		short := sss.TakeNext()
		return short
	}
	short := sss.txtScriptMap[txt]
	sss.used = append(sss.used, txt)
	dump.P(short, sss.txtScriptMap)
	return short
}

func (sss *ScriptShortcuts) Clear() {
	sss.used = make([]rune, 0)
}
func (sss *ScriptShortcuts) Avail() []rune {
	ava, avail := lo.Difference[rune](sss.used, sss.origText)
	dump.P(ava, avail)
	return avail
}
func (sss *ScriptShortcuts) SetWrapperStrings(b, e string) {
	sss.wrapBegin = b
	sss.wrapEnd = e
  sss.MapTxtShort()
  // for k, _ := range sss.txtScriptMap {
  //   sss.txtScriptMap[k].SetWrapperStrings(b, e)
  // }
}

func (sss *ScriptShortcuts) TakeNext() ScriptShortcut {
	// for _, v := range sss.origText {
	//   if lo.Contains[rune](sss.used, v) {
	//     continue
	//   }
	//   return sss.Take(v)
	// }
	avail := sss.Avail()
	if len(avail) == 0 {
		return NewScriptShortcut(rune(0), rune(0))
	} else {
		return sss.Take(avail[0])
	}
}

func NewScriptShortcuts(shortcutTexts, shortcutSubs []rune) *ScriptShortcuts {
	ssss := &ScriptShortcuts{
		txtScriptMap: make(map[rune]ScriptShortcut),
		used:         make([]rune, 0),
		origScript:   make([]rune, len(shortcutSubs)),
		origText:     make([]rune, len(shortcutTexts)),
	}

	copy(ssss.origScript, shortcutSubs)
	copy(ssss.origText, shortcutTexts)

  ssss.MapTxtShort()
	return ssss
}

func (ssss *ScriptShortcuts) MapTxtShort()  {
	for i, v := range ssss.origText {
		sub := ssss.origScript[i]
		sss := NewScriptShortcut(v, sub)
    sss.SetWrapperStrings(ssss.wrapBegin, ssss.wrapEnd)
		ssss.txtScriptMap[v] = sss
	}

}


func NewSuperScriptShortcuts() *ScriptShortcuts {
	sss := NewScriptShortcuts(superShortTexts, superShortSubs)
	sss.SetWrapperStrings("⁽", "⁾")
	return sss
}

func NewSubScriptShortcuts() *ScriptShortcuts {
	sss := NewScriptShortcuts(subShortTexts, subShortSubs)
	sss.SetWrapperStrings("₍", "₎")
	return sss
}

type ScriptShortcut struct {
	text               rune
	script             rune
	wrapBegin, wrapEnd string
}

func NewScriptShortcut(text, sub rune) ScriptShortcut {
	sss := &ScriptShortcut{
		text:      text,
		script:    sub,
		wrapBegin: "",
		wrapEnd:   "",
	}

	return *sss
}

func (sss ScriptShortcut) Script() rune {
	return sss.script
}

func (sss ScriptShortcut) Text() rune {
	return sss.text
}

func (sss *ScriptShortcut) SetWrapperStrings(b, e string) {
	sss.wrapBegin = b
	sss.wrapEnd = e
}

func (sss ScriptShortcut) String() string {
	return fmt.Sprintf("%s%c%s", sss.wrapBegin, sss.script, sss.wrapEnd)
}

var keyLabel []rune
var subLabel []rune

var superShortTexts []rune
var superShortSubs []rune

var subShortTexts []rune
var subShortSubs []rune

func init() {
	superShortTexts = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '=', 'n'}                                                      // '(', ')',
	superShortSubs = []rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹', '⁺', '⁻', '⁼', 'ⁿ'}                                                       // '⁽', '⁾',
	subShortSubs = []rune{'ₐ', 'ₑ', 'ₒ', 'ₓ', 'ₕ', 'ₖ', 'ₗ', 'ₘ', 'ₙ', 'ₚ', 'ₛ', 'ₜ', '₀', '₁', '₂', '₃', '₄', '₅', '₆', '₇', '₈', '₉', '₊', '₋', '₌'}  // '⁽', '⁾''₍', '₎'
	subShortTexts = []rune{'a', 'e', 'o', 'x', 'h', 'k', 'l', 'm', 'n', 'p', 's', 't', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '='} // '(', ')'

	SubShortcuts = NewSubScriptShortcuts()
	SuperShortcuts = NewSuperScriptShortcuts()

	keyLabel = []rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹'}
	subLabel = []rune{'ₐ', 'ₑ', 'ₒ', 'ₓ', 'ₕ', 'ₖ', 'ₗ', 'ₘ', 'ₙ', 'ₚ', 'ₛ', 'ₜ', '₀', '₁', '₂', '₃', '₄', '₅', '₆', '₇', '₈', '₉', '₊', '₋', '₌', '₍', '₎'}
	// ¹²³⁴⁵⁶⁷⁸⁹⁺⁻⁼⁽⁾ⁿ⁰₀₁₂₃₄₅₆₇₈₉ₔₐₑₒₓₕₖₗₘₙₚₛₜ₊₋₌₍₎ﴞ ﵆ ﵃ ﵄
}

// '⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹', '⁺', '⁻', '⁼', '⁽', '⁾', 'ⁿ',
// '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '=', '(', ')', 'n',
// ₐₑₒₓₕₖₗₘₙₚₛₜ₀₁₂₃₄₅₆₇₈₉₊₋₌₍₎
// aeoxhklmnpst0123456789+-=()
