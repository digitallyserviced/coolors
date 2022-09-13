package coolor

import (
	// "fmt"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"

	// "github.com/gookit/goutil/dump"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

type ListItemsVisibility uint16

const (
	ListItemDefault ListItemsVisibility = 1 << iota
	ListItemVisible
	ListItemNotVisible
	ListItemHidden
)

// listerItem represents one item in a List.
type listerItem struct {
	selected      func()
	mainText      string
	secondaryText string
	shortcut      rune
}

type ListItemSelected interface {
	Selected(idx int, i interface{}, lis []*ListItem)
	Changed(idx int, selected bool, i interface{}, lis []*ListItem)
	Cancelled(idx int, i interface{}, lis []*ListItem)
}

type ListItemText interface {
	MainText() string
	SecondaryText() string
	Shortcut() ScriptShortcut
}

type ListStyles struct {
	main, sec, short, sel string
}

type ListItemDrawable interface {
	GetPrimitive() tview.Primitive
}

type ListItemVisibility interface {
	Visibility() ListItemsVisibility
}

type ListItem interface {
	ListItemText
	ListItemSelected
	ListItemVisibility
}

func (li listerItem) Visibility() ListItemsVisibility {
	return ListItemVisible
}

func (li listerItem) MainText() string {
	return li.mainText
}

func (li listerItem) SecondaryText() string {
	return li.secondaryText
}

func (li listerItem) Shortcut() rune {
	return li.shortcut
}

func (li listerItem) Selected(idx int, i interface{}, lis []*ListItem) {
}

func (li listerItem) Changed(idx int, selected bool, i interface{}, lis []*ListItem) {
}

type ListerHighlightType uint

const (
	ListerHighlightDefault ListerHighlightType = 1 << iota
	ListerHighlightArrow
	ListerHighlightBars
	ListerHighlightThinArrow
)

// Lister displays rows of items, each of which can be selected.
//
// See https://github.com/rivo/tview/wiki/Lister for an example.
type Lister struct {
	selectedStyle      tcell.Style
	mainTextStyle      tcell.Style
	secondaryTextStyle tcell.Style
	shortcutStyle      tcell.Style
	selected           func(index int, mainText, secondaryText string, shortcut rune)
	changed            func(index int, mainText, secondaryText string, shortcut rune)
	*tview.Box
	done              func()
	canceled          func()
	itemsLister       func() []*ListItem
	items             []*ListItem
	currentItem       int
	itemOffset        int
	horizontalOffset  int
	wrapAround        bool
	overflowing       bool
	showSecondaryText bool
	highlightFullLine bool
	selectedFocusOnly bool
	highlightType     ListerHighlightType
}

// NewList returns a new form.
func NewLister() *Lister {
	list := &Lister{
		selectedStyle:      tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).Background(tview.Styles.PrimaryTextColor),
		mainTextStyle:      tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor),
		secondaryTextStyle: tcell.StyleDefault.Foreground(tview.Styles.TertiaryTextColor),
		shortcutStyle:      tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor),
		selected: func(index int, mainText string, secondaryText string, shortcut rune) {
		},
		changed: func(index int, mainText string, secondaryText string, shortcut rune) {
		},
		Box: tview.NewBox(),
		done: func() {
		},
		itemsLister: func() []*ListItem {
			return make([]*ListItem, 0)
		},
		items:             make([]*ListItem, 0),
		currentItem:       0,
		itemOffset:        0,
		horizontalOffset:  0,
		wrapAround:        true,
		overflowing:       false,
		showSecondaryText: true,
		highlightFullLine: false,
		selectedFocusOnly: false,
		highlightType:     ListerHighlightDefault,
	}
	list.SetHandlers(nil, nil)

	return list
}

// SetCurrentItem sets the currently selected item by its index, starting at 0
// for the first item. If a negative index is provided, items are referred to
// from the back (-1 = last item, -2 = second-to-last item, and so on). Out of
// range indices are clamped to the beginning/end.
//
// Calling this function triggers a "changed" event if the selection changes.
func (l *Lister) SetCurrentItem(index int) *Lister {
	if index < 0 {
		index = len(l.items) + index
	}
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	if index < 0 {
		index = 0
	}

	if index != l.currentItem && l.changed != nil {
		item := *l.items[index]
		l.changed(index, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
	}

	l.currentItem = index

	return l
}

func (l *Lister) Each(flag ListItemsVisibility, f func(i interface{}, idx int, visIdx int)) {
	index := 0
	for i, v := range l.items {
		if v == nil {
			continue
		}
		var li ListItem = *v
		if flag&li.Visibility() != 0 || flag == ListItemDefault {
			f(li, i, index)
			index += 1
		}
	}
}

// SetCurrent sets the item based on a reference to the actual item instead of index
// ranges through the items to find passed reference and updates the currentItem
// index value

func (l *Lister) SetCurrent(i interface{}) *Lister {
	index := -1
	for num, v := range l.items {
		if v == i {
			index = num
		}
	}
  l.SetCurrentItem(index)
	return l
}

// GetCurrentItem returns the index of the currently selected list item,
// starting at 0 for the first item.
func (l *Lister) GetCurrentItem() int {
	return l.currentItem
}

// SetOffset sets the number of items to be skipped (vertically) as well as the
// number of cells skipped horizontally when the list is drawn. Note that one
// item corresponds to two rows when there are secondary texts. Shortcut()s are
// always drawn.
//
// These values may change when the list is drawn to ensure the currently
// selected item is visible and item texts move out of view. Users can also
// modify these values by interacting with the list.
func (l *Lister) SetOffset(items, horizontal int) *Lister {
	l.itemOffset = items
	l.horizontalOffset = horizontal
	return l
}

// GetOffset returns the number of items skipped while drawing, as well as the
// number of cells item text is moved to the left. See also SetOffset() for more
// information on these values.
func (l *Lister) GetOffset() (int, int) {
	return l.itemOffset, l.horizontalOffset
}

// RemoveItem removes the item with the given index (starting at 0) from the
// list. If a negative index is provided, items are referred to from the back
// (-1 = last item, -2 = second-to-last item, and so on). Out of range indices
// are clamped to the beginning/end, i.e. unless the list is empty, an item is
// always removed.
//
// The currently selected item is shifted accordingly. If it is the one that is
// removed, a "changed" event is fired.
func (l *Lister) RemoveItem(index int) *Lister {
	if len(l.items) == 0 {
		return l
	}

	// Adjust index.
	if index < 0 {
		index = len(l.items) + index
	}
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	if index < 0 {
		index = 0
	}

	// Remove item.
	l.items = append(l.items[:index], l.items[index+1:]...)

	// If there is nothing left, we're done.
	if len(l.items) == 0 {
		return l
	}

	// Shift current item.
	previousCurrentItem := l.currentItem
	if l.currentItem >= index {
		l.currentItem--
	}

	// Fire "changed" event for removed items.
	if previousCurrentItem == index && l.changed != nil {
		item := *l.items[l.currentItem]
		l.changed(l.currentItem, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
	}

	return l
}

// SetMainText()Color sets the color of the items' main text.
func (l *Lister) SetMainTextColor(color tcell.Color) *Lister {
	l.mainTextStyle = l.mainTextStyle.Foreground(color)
	return l
}

// SetMainText()Style sets the style of the items' main text. Note that the
// background color is ignored in order not to override the background color of
// the list itself.
func (l *Lister) SetMainTextStyle(style tcell.Style) *Lister {
	l.mainTextStyle = style
	return l
}

// SetSecondaryText()Color sets the color of the items' secondary text.
func (l *Lister) SetSecondaryTextColor(color tcell.Color) *Lister {
	l.secondaryTextStyle = l.secondaryTextStyle.Foreground(color)
	return l
}

// SetSecondaryText()Style sets the style of the items' secondary text. Note that
// the background color is ignored in order not to override the background color
// of the list itself.
func (l *Lister) SetSecondaryTextStyle(style tcell.Style) *Lister {
	l.secondaryTextStyle = style
	return l
}

// SetShortcut()Color sets the color of the items' shortcut.
func (l *Lister) SetShortcutColor(color tcell.Color) *Lister {
	l.shortcutStyle = l.shortcutStyle.Foreground(color)
	return l
}

// SetShortcut()Style sets the style of the items' shortcut. Note that the
// background color is ignored in order not to override the background color of
// the list itself.
func (l *Lister) SetShortcutStyle(style tcell.Style) *Lister {
	l.shortcutStyle = style
	return l
}

// SetSelectedTextColor sets the text color of selected items. Note that the
// color of main text characters that are different from the main text color
// (e.g. color tags) is maintained.
func (l *Lister) SetSelectedTextColor(color tcell.Color) *Lister {
	l.selectedStyle = l.selectedStyle.Foreground(color)
	return l
}

// SetSelectedBackgroundColor sets the background color of selected items.
func (l *Lister) SetSelectedBackgroundColor(color tcell.Color) *Lister {
	l.selectedStyle = l.selectedStyle.Background(color)
	return l
}

// func (li listerItem) Selected(idx int, i interface{}, lis []*ListItem) {
// }
//
// func (li listerItem) Changed(idx int, selected bool, i interface{}, lis []*ListItem) {
// }

func (list *Lister) SetHandlers(sel func(idx int, i interface{}, lis []*ListItem), chg func(idx int, selected bool, i interface{}, lis []*ListItem)) {
	// items := &list.items
	list.SetChangedFunc(func(index int, s1, s2 string, r rune) {
		// tem := list.GetItem(index)
		// dump.P(index,tem, s1, s2, r)
		if list == nil || list.items == nil {
			return
		}
		if len(list.items) == 0 {
			return
		}
		if index < 0 {
			index = len(list.items) + index
		}
		if index >= len(list.items) {
			index = len(list.items) - 1
		}
		if index < 0 {
			index = 0
		}

		item := *list.items[index]
		selItem := list.GetCurrentItem() == index
		if chg != nil {
			chg(index, selItem, item, list.items)
		}
		item.Changed(index, selItem, item, list.items)
	})

	list.SetSelectedFunc(func(index int, s1, s2 string, r rune) {
		if list == nil || list.items == nil {
			return
		}
		if len(list.items) == 0 {
			return
		}
		if index < 0 {
			index = len(list.items) + index
		}
		if index >= len(list.items) {
			index = len(list.items) - 1
		}
		if index < 0 {
			index = 0
		}

		item := *list.items[index]
		if sel != nil {
			sel(index, item, list.items)
		}
		item.Selected(index, item, list.items)
	})
}

// SetSelectedStyle sets the style of the selected items. Note that the color of
// main text characters that are different from the main text color (e.g. color
// tags) is maintained.
func (l *Lister) SetSelectedStyle(style tcell.Style) *Lister {
	l.selectedStyle = style
	return l
}

// SetSelectedFocusOnly sets a flag which determines when the currently selected
// list item is highlighted. If set to true, selected items are only highlighted
// when the list has focus. If set to false, they are always highlighted.
func (l *Lister) SetSelectedFocusOnly(focusOnly bool) *Lister {
	l.selectedFocusOnly = focusOnly
	return l
}

// SetHighlightFullLine sets a flag which determines whether the colored
// background of selected items spans the entire width of the view. If set to
// true, the highlight spans the entire view. If set to false, only the text of
// the selected item from beginning to end is highlighted.
func (l *Lister) SetHighlightFullLine(highlight bool) *Lister {
	l.highlightFullLine = highlight
	return l
}

// ShowSecondaryText() determines whether or not to show secondary item texts.
func (l *Lister) ShowSecondaryText(show bool) *Lister {
	l.showSecondaryText = show
	return l
}

// SetWrapAround sets the flag that determines whether navigating the list will
// wrap around. That is, navigating downwards on the last item will move the
// selection to the first item (similarly in the other direction). If set to
// false, the selection won't change when navigating downwards on the last item
// or navigating upwards on the first item.
func (l *Lister) SetWrapAround(wrapAround bool) *Lister {
	l.wrapAround = wrapAround
	return l
}

// SetChangedFunc sets the function which is called when the user navigates to
// a list item. The function receives the item's index in the list of items
// (starting with 0), its main text, secondary text, and its shortcut rune.
//
// This function is also called when the first item is added or when
// SetCurrentItem() is called.
func (l *Lister) SetChangedFunc(handler func(index int, mainText string, secondaryText string, shortcut rune)) *Lister {
	l.changed = handler
	return l
}

// SetSelectedFunc sets the function which is called when the user selects a
// list item by pressing Enter on the current selection. The function receives
// the item's index in the list of items (starting with 0), its main text,
// secondary text, and its shortcut rune.
func (l *Lister) SetSelectedFunc(handler func(int, string, string, rune)) *Lister {
	l.selected = handler
	return l
}

// SetDoneFunc sets a function which is called when the user presses the Escape
// key.
func (l *Lister) SetDoneFunc(handler func()) *Lister {
	l.done = handler
	return l
}

// SetDoneFunc sets a function which is called when the user presses the Escape
// key.
func (l *Lister) SetCancelFunc(handler func()) *Lister {
	l.canceled = handler
	return l
}

// AddItem calls InsertItem() with an index of -1.
// func (l *Lister) AddItem(mainText, secondaryText string, shortcut rune, selected func()) *Lister {
// 	l.InsertItem(-1, mainText, secondaryText, shortcut, selected)
// 	return l
// }
func (l *Lister) AddItem(li *ListItem) *Lister {
	l.InsertItem(-1, li)
	return l
}

// InsertItem adds a new item to the list at the specified index. An index of 0
// will insert the item at the beginning, an index of 1 before the second item,
// and so on. An index of GetItemCount() or higher will insert the item at the
// end of the list. Negative indices are also allowed: An index of -1 will
// insert the item at the end of the list, an index of -2 before the last item,
// and so on. An index of -GetItemCount()-1 or lower will insert the item at the
// beginning.
//
// An item has a main text which will be highlighted when selected. It also has
// a secondary text which is shown underneath the main text (if it is set to
// visible) but which may remain empty.
//
// The shortcut is a key binding. If the specified rune is entered, the item
// is selected immediately. Set to 0 for no binding.
//
// The "selected" callback will be invoked when the user selects the item. You
// may provide nil if no such callback is needed or if all events are handled
// through the selected callback set with SetSelectedFunc().
//
// The currently selected item will shift its position accordingly. If the list
// was previously empty, a "changed" event is fired because the new item becomes
// selected.
func (l *Lister) InsertItem(index int, item *ListItem) *Lister {
	// item := &listerItem{
	// 	mainText:      mainText,
	// 	secondaryText: secondaryText,
	// 	shortcut:      shortcut,
	// 	selected:      selected,
	// }
	//
	// Shift index to range.
	if index < 0 {
		index = len(l.items) + index + 1
	}
	if index < 0 {
		index = 0
	} else if index > len(l.items) {
		index = len(l.items)
	}

	// Shift current item.
	if l.currentItem < len(l.items) && l.currentItem >= index {
		l.currentItem++
	}

	// Insert item (make space for the new item, then shift and insert).
	l.items = append(l.items, nil)
	if index < len(l.items)-1 { // -1 because l.items has already grown by one item.
		copy(l.items[index+1:], l.items[index:])
	}

	// var litem *ListItem = item.(*ListItem)

	litem := ListItem(*item)
	l.items[index] = &litem
	// l.items[index] = ListItem(*item)

	// Fire a "change" event for the first item in the list.
	if len(l.items) == 1 && l.changed != nil {
		item := *l.items[0]
		l.changed(0, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
	}

	return l
}

// GetItem returns the ListItem at the specified index in the list.
func (l *Lister) GetItem(index int) *ListItem {
	if index < 0 {
		index = len(l.items) - 1 + index
	}
	if index >= len(l.items) {
		index = index - len(l.items) - 1
		// return l.GetItem(index)
	}
	if index < 0 {
		index = 0
	}

	var item ListItem = (*l.items[index])
	// }

	// l.currentItem = index

	return &item
}

// GetItemCount returns the number of items in the list.
func (l *Lister) GetItemCount() int {
	return len(l.items)
}

// GetItemText returns an item's texts (main and secondary). Panics if the index
// is out of range.
func (l *Lister) GetItemText(index int) (main, secondary string) {
	return (*l.items[index]).MainText(), (*l.items[index]).SecondaryText()
}

// SetItemText sets an item's main and secondary text. Panics if the index is
// out of range.
// func (l *Lister) SetItemText(index int, main, secondary string) *Lister {
// 	item := l.items[index]
// 	item.MainText() = main
// 	item.SecondaryText() = secondary
// 	return l
// }

// FindItems searches the main and secondary texts for the given strings and
// returns a list of item indices in which those strings are found. One of the
// two search strings may be empty, it will then be ignored. Indices are always
// returned in ascending order.
//
// If mustContainBoth is set to true, mainSearch must be contained in the main
// text AND secondarySearch must be contained in the secondary text. If it is
// false, only one of the two search strings must be contained.
//
// Set ignoreCase to true for case-insensitive search.

func (l *Lister) ClampSelectionToVisible() {
	vis := l.FindVisibleItems(ListItemDefault | ListItemVisible)
	if lo.Contains[int](vis, l.currentItem) {
		return
	}
	if len(vis) == 0 {
		return
	}
	if len(vis) == 1 {
		l.currentItem = vis[0]
	}
	best, bidx := 0, 0
	if len(vis) > 1 {
		for i, v := range vis {
			if i+1 >= len(vis) {
				l.currentItem = v
				return
			}
			c := l.currentItem - v
			if c < best {
				best = c
				bidx = v
			}
		}
		l.currentItem = bidx
	}
	ind := l.currentItem
	for {
		if lo.Contains[int](vis, l.currentItem) {
			return
		}
		ind++
	}
}

func (l *Lister) FindVisibleItems(flag ListItemsVisibility) (indices []int) {
	indices = make([]int, 0)
	for index, itemP := range l.items {
		item := (*itemP)
		if item.Visibility()&(flag) == 0 {
			continue
		}
		indices = append(indices, index)
	}
	return indices
}

func (l *Lister) FindItems(mainSearch, secondarySearch string, mustContainBoth, ignoreCase bool) (indices []int) {
	if mainSearch == "" && secondarySearch == "" {
		return
	}

	if ignoreCase {
		mainSearch = strings.ToLower(mainSearch)
		secondarySearch = strings.ToLower(secondarySearch)
	}

	for index, itemP := range l.items {
		item := *itemP
		mainText := item.MainText()
		secondaryText := item.SecondaryText()
		if ignoreCase {
			mainText = strings.ToLower(mainText)
			secondaryText = strings.ToLower(secondaryText)
		}

		// strings.Contains() always returns true for a "" search.
		mainContained := strings.Contains(mainText, mainSearch)
		secondaryContained := strings.Contains(secondaryText, secondarySearch)
		if mustContainBoth && mainContained && secondaryContained ||
			!mustContainBoth && (mainText != "" && mainContained || secondaryText != "" && secondaryContained) {
			indices = append(indices, index)
		}
	}

	return
}

func (f *Lister) SetItemLister(il func() []*ListItem) {
	f.itemsLister = il
}

func (f *Lister) UpdateListItems() {
	if f.itemsLister != nil {
		f.items = f.itemsLister()
	}
}

func (f *Lister) SetListItems(li []*ListItem) {
	f.items = li
}

// Clear removes all items from the list.
func (l *Lister) Clear() *Lister {
	l.items = nil
	l.currentItem = 0
	return l
}

// Draw draws this primitive onto the screen.
func (l *Lister) Draw(screen tcell.Screen) {
	l.DrawForSubclass(screen, l)

	// Determine the dimensions.
	x, y, width, height := l.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	// Do we show any shortcuts?
	var showShortcuts bool
	for _, itemP := range l.items {
		item := *itemP
		if item.Shortcut().Text() != 0 {
			showShortcuts = true
			x += 5
			width -= 5
			break
		}
	}

	// Adjust offset to keep the current selection in view.
	if l.currentItem < l.itemOffset {
		l.itemOffset = l.currentItem
	} else if l.showSecondaryText {
		if 2*(l.currentItem-l.itemOffset) >= height-1 {
			l.itemOffset = (2*l.currentItem + 3 - height) / 2
		}
	} else {
		if l.currentItem-l.itemOffset >= height {
			l.itemOffset = l.currentItem + 1 - height
		}
	}
	if l.horizontalOffset < 0 {
		l.horizontalOffset = 0
	}

	// Draw the list items.
	var (
		maxWidth int // The maximum printed item width.
		// overflowing bool // Whether a text's end exceeds the right border.
	)
	// k, v := range l.items

	si := 0
	for index, itemP := range l.items {
		item := (*itemP)
		if index < l.itemOffset || item.Visibility()&(ListItemDefault|ListItemVisible) == 0 {
			continue
		}

		if y >= bottomLimit {
			break
		}

		// Shortcut()s.
		if showShortcuts && item.Shortcut().Script() != 0 {
			shortStr := fmt.Sprintf(" %s", item.Shortcut().String())
			dump.P(shortStr)
			printWithStyle(screen, shortStr, x-5, y, 0, len(shortStr), AlignLeft, l.shortcutStyle, true)
		}

		// Main text.
		_, printedWidth, _, end := printWithStyle(screen, item.MainText(), x, y, l.horizontalOffset, width, tview.AlignLeft, l.mainTextStyle, true)
		if printedWidth > maxWidth {
			maxWidth = printedWidth
		}
		if end < len(item.MainText()) {
			// overflowing = true
		}

		// Background color of selected text.
		if index == l.currentItem && (!l.selectedFocusOnly || l.HasFocus()) {
			textWidth := width
			if l.highlightType == ListerHighlightDefault {
				if !l.highlightFullLine {
					if w := tview.TaggedStringWidth(item.MainText()); w < textWidth {
						textWidth = w
					}
				}
				mainTextColor, _, _ := l.mainTextStyle.Decompose()
				for bx := 0; bx < textWidth; bx++ {
					m, c, style, _ := screen.GetContent(x+bx, y)
					fg, _, _ := style.Decompose()
					style = l.selectedStyle
					if fg != mainTextColor {
						style = style.Foreground(fg)
					}
					screen.SetContent(x+bx, y, m, c, style)
				}
				// ❱ ➤ ┃
			}
			highlightSymbol := "SHIT"
			switch l.highlightType {
			case ListerHighlightDefault:
			case ListerHighlightBars:
				highlightSymbol = fmt.Sprintf("%s", "❱")
			case ListerHighlightThinArrow:
				highlightSymbol = fmt.Sprintf("%s", "➤")
				// case ListerHighlightBars:
				// 	highlightSymbol = fmt.Sprintf("%s", "┃")
			}
			if l.highlightType != ListerHighlightDefault {
				printWithStyle(screen, highlightSymbol, x-5, y, 0, len(highlightSymbol), AlignLeft, tcell.StyleDefault.Foreground(tcell.ColorRed), true)
			}
		}

		y++

		if y >= bottomLimit {
			break
		}

		// Secondary text.
		if l.showSecondaryText {
			_, printedWidth, _, end := printWithStyle(screen, item.SecondaryText(), x, y, l.horizontalOffset, width, tview.AlignLeft, l.secondaryTextStyle, true)
			if printedWidth > maxWidth {
				maxWidth = printedWidth
			}
			if end < len(item.SecondaryText()) {
				// overflowing = true
			}
			y++
		}
		si++
	}

	// We don't want the item text to get out of view. If the horizontal offset
	// is too high, we reset it and redraw. (That should be about as efficient
	// as calculating everything up front.)
	// if l.horizontalOffset > 0 && maxWidth < width {
	// 	l.horizontalOffset -= width - maxWidth
	// 	l.Draw(screen)
	// }
	// l.overflowing = overflowing
}

// InputHandler returns the handler for this primitive.
func (l *Lister) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return l.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyF40 {
			if l.done != nil {
				l.done()
			}
			return
		}
		if event.Key() == tcell.KeyEscape {
			if l.canceled != nil {
				if l.currentItem >= 0 && l.currentItem < len(l.items) {
					item := *l.items[l.currentItem]
					item.Cancelled(l.currentItem, item, l.items)
					if l.selected != nil {
						l.canceled()
					}
				}
			}
		} else if len(l.items) == 0 {
			return
		}

		previousItem := l.currentItem
		key := event.Key()
		ch := event.Rune()
		switch {
		case key == tcell.KeyTab || key == tcell.KeyDown || ch == 'j':
			l.currentItem++
		case key == tcell.KeyBacktab || key == tcell.KeyUp || ch == 'k':
			l.currentItem--
		// case tcell.KeyRight:
		// 	if l.overflowing {
		// 		l.horizontalOffset += 2 // We shift by 2 to account for two-cell characters.
		// 	} else {
		// 		l.currentItem++
		// 	}
		// case tcell.KeyLeft:
		// 	if l.horizontalOffset > 0 {
		// 		l.horizontalOffset -= 2
		// 	} else {
		// 		l.currentItem--
		// 	}
		case key == tcell.KeyHome:
			l.currentItem = 0
		case key == tcell.KeyEnd:
			l.currentItem = len(l.items) - 1
		case key == tcell.KeyPgDn:
			_, _, _, height := l.GetInnerRect()
			l.currentItem += height
			if l.currentItem >= len(l.items) {
				l.currentItem = len(l.items) - 1
			}
		case key == tcell.KeyPgUp:
			_, _, _, height := l.GetInnerRect()
			l.currentItem -= height
			if l.currentItem < 0 {
				l.currentItem = 0
			}
		case key == tcell.KeyEnter:
			if l.currentItem >= 0 && l.currentItem < len(l.items) {
				item := *l.items[l.currentItem]
				// item.Selected(l.currentItem, item, l.items)
				if l.selected != nil {
					l.selected(l.currentItem, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
				}
			}
		case key == tcell.KeyRune:
			ch := event.Rune()
			if ch != ' ' {
				// It's not a space bar. Is it a shortcut?
				var found bool
				for index, itemP := range l.items {
					item := (*itemP)
					if item.Shortcut().Text() == ch {
						// We have a shortcut.
						item.Selected(l.currentItem, item, l.items)
						found = true
						l.currentItem = index
						break
					}
				}
				if !found {
					break
				}
			}
			item := *l.items[l.currentItem]
			// item.Selected(l.currentItem, item, l.items)
			// if item.Selected != nil {
			// item.Selected()
			// }
			if l.selected != nil {
				l.selected(l.currentItem, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
			}
		}

		if l.currentItem < 0 {
			if l.wrapAround {
				l.currentItem = len(l.items) - 1
			} else {
				l.currentItem = 0
			}
		} else if l.currentItem >= len(l.items) {
			if l.wrapAround {
				l.currentItem = 0
			} else {
				l.currentItem = len(l.items) - 1
			}
		}

		l.ClampSelectionToVisible()

		if l.currentItem != previousItem && l.currentItem < len(l.items) && l.changed != nil {
			item := *l.items[l.currentItem]
			l.changed(l.currentItem, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
		}
	})
}

// indexAtPoint returns the index of the list item found at the given position
// or a negative value if there is no such list item.
func (l *Lister) indexAtPoint(x, y int) int {
	rectX, rectY, width, height := l.GetInnerRect()
	if rectX < 0 || rectX >= rectX+width || y < rectY || y >= rectY+height {
		return -1
	}

	index := y - rectY
	if l.showSecondaryText {
		index /= 2
	}
	index += l.itemOffset

	if index >= len(l.items) {
		return -1
	}
	return index
}

// MouseHandler returns the mouse handler for this primitive.
func (l *Lister) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return l.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !l.InRect(event.Position()) {
			return false, nil
		}

		// Process mouse event.
		switch action {
		case tview.MouseLeftClick:
			setFocus(l)
			index := l.indexAtPoint(event.Position())
			if index != -1 {
				item := *l.items[index]
				item.Selected(index, item, l.items)
				// item.Selected
				if l.selected != nil {
					l.selected(index, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
				}
				if index != l.currentItem && l.changed != nil {
					l.changed(index, item.MainText(), item.SecondaryText(), item.Shortcut().Text())
				}
				l.currentItem = index
			}
			consumed = true
		case tview.MouseScrollUp:
			if l.itemOffset > 0 {
				l.itemOffset--
			}
			consumed = true
		case tview.MouseScrollDown:
			lines := len(l.items) - l.itemOffset
			if l.showSecondaryText {
				lines *= 2
			}
			if _, _, _, height := l.GetInnerRect(); lines > height {
				l.itemOffset++
			}
			consumed = true
		}

		return
	})
}

func printWithStyle(screen tcell.Screen, text string, x, y, skipWidth, maxWidth, align int, style tcell.Style, maintainBackground bool) (int, int, int, int) {
	totalWidth, totalHeight := screen.Size()
	if maxWidth <= 0 || len(text) == 0 || y < 0 || y >= totalHeight {
		return 0, 0, 0, 0
	}

	// Decompose the text.
	colorIndices, colors, _, _, escapeIndices, strippedText, strippedWidth := decomposeString(text, true, false)

	// We want to reduce all alignments to AlignLeft.
	if align == AlignRight {
		if strippedWidth-skipWidth <= maxWidth {
			// There's enough space for the entire text.
			return printWithStyle(screen, text, x+maxWidth-strippedWidth+skipWidth, y, skipWidth, maxWidth, AlignLeft, style, maintainBackground)
		}
		// Trim characters off the beginning.
		var (
			bytes, width, colorPos, escapePos, tagOffset, from, to int
			foregroundColor, backgroundColor, attributes           string
		)
		originalStyle := style
		iterateString(strippedText, func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Update color/escape tag offset and style.
			if colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
				foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
				style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
				tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
				colorPos++
			}
			if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
				tagOffset++
				escapePos++
			}
			if strippedWidth-screenPos <= maxWidth {
				// We chopped off enough.
				if escapePos > 0 && textPos+tagOffset-1 >= escapeIndices[escapePos-1][0] && textPos+tagOffset-1 < escapeIndices[escapePos-1][1] {
					// Unescape open escape sequences.
					escapeCharPos := escapeIndices[escapePos-1][1] - 2
					text = text[:escapeCharPos] + text[escapeCharPos+1:]
				}
				// Print and return.
				bytes, width, from, to = printWithStyle(screen, text[textPos+tagOffset:], x, y, 0, maxWidth, AlignLeft, style, maintainBackground)
				from += textPos + tagOffset
				to += textPos + tagOffset
				return true
			}
			return false
		})
		return bytes, width, from, to
	} else if align == AlignCenter {
		if strippedWidth-skipWidth == maxWidth {
			// Use the exact space.
			return printWithStyle(screen, text, x, y, skipWidth, maxWidth, AlignLeft, style, maintainBackground)
		} else if strippedWidth-skipWidth < maxWidth {
			// We have more space than we need.
			half := (maxWidth - strippedWidth + skipWidth) / 2
			return printWithStyle(screen, text, x+half, y, skipWidth, maxWidth-half, AlignLeft, style, maintainBackground)
		} else {
			// Chop off runes until we have a perfect fit.
			var choppedLeft, choppedRight, leftIndex, rightIndex int
			rightIndex = len(strippedText)
			for rightIndex-1 > leftIndex && strippedWidth-skipWidth-choppedLeft-choppedRight > maxWidth {
				if skipWidth > 0 || choppedLeft < choppedRight {
					// Iterate on the left by one character.
					iterateString(strippedText[leftIndex:], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						if skipWidth > 0 {
							skipWidth -= screenWidth
							strippedWidth -= screenWidth
						} else {
							choppedLeft += screenWidth
						}
						leftIndex += textWidth
						return true
					})
				} else {
					// Iterate on the right by one character.
					iterateStringReverse(strippedText[leftIndex:rightIndex], func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedRight += screenWidth
						rightIndex -= textWidth
						return true
					})
				}
			}

			// Add tag offsets and determine start style.
			var (
				colorPos, escapePos, tagOffset               int
				foregroundColor, backgroundColor, attributes string
			)
			originalStyle := style
			for index := range strippedText {
				// We only need the offset of the left index.
				if index > leftIndex {
					// We're done.
					if escapePos > 0 && leftIndex+tagOffset-1 >= escapeIndices[escapePos-1][0] && leftIndex+tagOffset-1 < escapeIndices[escapePos-1][1] {
						// Unescape open escape sequences.
						escapeCharPos := escapeIndices[escapePos-1][1] - 2
						text = text[:escapeCharPos] + text[escapeCharPos+1:]
					}
					break
				}

				// Update color/escape tag offset.
				if colorPos < len(colorIndices) && index+tagOffset >= colorIndices[colorPos][0] && index+tagOffset < colorIndices[colorPos][1] {
					if index <= leftIndex {
						foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
						style = overlayStyle(originalStyle, foregroundColor, backgroundColor, attributes)
					}
					tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
					colorPos++
				}
				if escapePos < len(escapeIndices) && index+tagOffset >= escapeIndices[escapePos][0] && index+tagOffset < escapeIndices[escapePos][1] {
					tagOffset++
					escapePos++
				}
			}
			bytes, width, from, to := printWithStyle(screen, text[leftIndex+tagOffset:], x, y, 0, maxWidth, AlignLeft, style, maintainBackground)
			from += leftIndex + tagOffset
			to += leftIndex + tagOffset
			return bytes, width, from, to
		}
	}

	// Draw text.
	var (
		drawn, drawnWidth, colorPos, escapePos, tagOffset, from, to int
		foregroundColor, backgroundColor, attributes                string
	)
	iterateString(strippedText, func(main rune, comb []rune, textPos, length, screenPos, screenWidth int) bool {
		// Skip character if necessary.
		if skipWidth > 0 {
			skipWidth -= screenWidth
			from = textPos + length
			to = from
			return false
		}

		// Only continue if there is still space.
		if drawnWidth+screenWidth > maxWidth || x+drawnWidth >= totalWidth {
			return true
		}

		// Handle color tags.
		for colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
			foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
			tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
			colorPos++
		}

		// Handle escape tags.
		if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
			if textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}
		}

		// Memorize positions.
		to = textPos + length

		// Print the rune sequence.
		finalX := x + drawnWidth
		finalStyle := style
		if maintainBackground {
			_, _, existingStyle, _ := screen.GetContent(finalX, y)
			_, background, _ := existingStyle.Decompose()
			finalStyle = finalStyle.Background(background)
		}
		finalStyle = overlayStyle(finalStyle, foregroundColor, backgroundColor, attributes)
		for offset := screenWidth - 1; offset >= 0; offset-- {
			// To avoid undesired effects, we populate all cells.
			if offset == 0 {
				screen.SetContent(finalX+offset, y, main, comb, finalStyle)
			} else {
				screen.SetContent(finalX+offset, y, ' ', nil, finalStyle)
			}
		}

		// Advance.
		drawn += length
		drawnWidth += screenWidth

		return false
	})

	return drawn + tagOffset + len(escapeIndices), drawnWidth, from, to
}

func overlayStyle(style tcell.Style, fgColor, bgColor, attributes string) tcell.Style {
	_, _, defAttr := style.Decompose()

	if fgColor != "" && fgColor != "-" {
		style = style.Foreground(tcell.GetColor(fgColor))
	}

	if bgColor != "" && bgColor != "-" {
		style = style.Background(tcell.GetColor(bgColor))
	}

	if attributes == "-" {
		style = style.Bold(defAttr&tcell.AttrBold > 0).
			Italic(defAttr&tcell.AttrItalic > 0).
			Blink(defAttr&tcell.AttrBlink > 0).
			Reverse(defAttr&tcell.AttrReverse > 0).
			Underline(defAttr&tcell.AttrUnderline > 0).
			Dim(defAttr&tcell.AttrDim > 0)
	} else if attributes != "" {
		style = style.Normal()
		for _, flag := range attributes {
			switch flag {
			case 'l':
				style = style.Blink(true)
			case 'b':
				style = style.Bold(true)
			case 'i':
				style = style.Italic(true)
			case 'd':
				style = style.Dim(true)
			case 'r':
				style = style.Reverse(true)
			case 'u':
				style = style.Underline(true)
			case 's':
				style = style.StrikeThrough(true)
			}
		}
	}

	return style
}

// decomposeString returns information about a string which may contain color
// tags or region tags, depending on which ones are requested to be found. It
// returns the indices of the color tags (as returned by
// re.FindAllStringIndex()), the color tags themselves (as returned by
// re.FindAllStringSubmatch()), the indices of region tags and the region tags
// themselves, the indices of an escaped tags (only if at least color tags or
// region tags are requested), the string stripped by any tags and escaped, and
// the screen width of the stripped string.
func decomposeString(text string, findColors, findRegions bool) (colorIndices [][]int, colors [][]string, regionIndices [][]int, regions [][]string, escapeIndices [][]int, stripped string, width int) {
	// Shortcut for the trivial case.
	if !findColors && !findRegions {
		return nil, nil, nil, nil, nil, text, stringWidth(text)
	}

	// Get positions of any tags.
	if findColors {
		colorIndices = colorPattern.FindAllStringIndex(text, -1)
		colors = colorPattern.FindAllStringSubmatch(text, -1)
	}
	if findRegions {
		regionIndices = regionPattern.FindAllStringIndex(text, -1)
		regions = regionPattern.FindAllStringSubmatch(text, -1)
	}
	escapeIndices = escapePattern.FindAllStringIndex(text, -1)

	// Because the color pattern detects empty tags, we need to filter them out.
	for i := len(colorIndices) - 1; i >= 0; i-- {
		if colorIndices[i][1]-colorIndices[i][0] == 2 {
			colorIndices = append(colorIndices[:i], colorIndices[i+1:]...)
			colors = append(colors[:i], colors[i+1:]...)
		}
	}

	// Make a (sorted) list of all tags.
	allIndices := make([][3]int, 0, len(colorIndices)+len(regionIndices)+len(escapeIndices))
	for indexType, index := range [][][]int{colorIndices, regionIndices, escapeIndices} {
		for _, tag := range index {
			allIndices = append(allIndices, [3]int{tag[0], tag[1], indexType})
		}
	}
	sort.Slice(allIndices, func(i int, j int) bool {
		return allIndices[i][0] < allIndices[j][0]
	})

	// Remove the tags from the original string.
	var from int
	buf := make([]byte, 0, len(text))
	for _, indices := range allIndices {
		if indices[2] == 2 { // Escape sequences are not simply removed.
			buf = append(buf, []byte(text[from:indices[1]-2])...)
			buf = append(buf, ']')
			from = indices[1]
		} else {
			buf = append(buf, []byte(text[from:indices[0]])...)
			from = indices[1]
		}
	}
	buf = append(buf, text[from:]...)
	stripped = string(buf)

	// Get the width of the stripped string.
	width = stringWidth(stripped)

	return
}

func styleFromTag(fgColor, bgColor, attributes string, tagSubstrings []string) (newFgColor, newBgColor, newAttributes string) {
	if tagSubstrings[colorForegroundPos] != "" {
		color := tagSubstrings[colorForegroundPos]
		if color == "-" {
			fgColor = "-"
		} else if color != "" {
			fgColor = color
		}
	}

	if tagSubstrings[colorBackgroundPos-1] != "" {
		color := tagSubstrings[colorBackgroundPos]
		if color == "-" {
			bgColor = "-"
		} else if color != "" {
			bgColor = color
		}
	}

	if tagSubstrings[colorFlagPos-1] != "" {
		flags := tagSubstrings[colorFlagPos]
		if flags == "-" {
			attributes = "-"
		} else if flags != "" {
			attributes = flags
		}
	}

	return fgColor, bgColor, attributes
}

// TaggedStringWidth returns the width of the given string needed to print it on
// screen. The text may contain color tags which are not counted.
func TaggedStringWidth(text string) int {
	_, _, _, _, _, _, width := decomposeString(text, true, false)
	return width
}

// stringWidth returns the number of horizontal cells needed to print the given
// text. It splits the text into its grapheme clusters, calculates each
// cluster's width, and adds them up to a total.
func stringWidth(text string) (width int) {
	g := uniseg.NewGraphemes(text)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = runewidth.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return
}

// Text alignment within a box.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// Common regular expressions.
var (
	colorPattern     = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	regionPattern    = regexp.MustCompile(`\["([a-zA-Z0-9_,;: \-\.]*)"\]`)
	escapePattern    = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
	nonEscapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\]`)
	boundaryPattern  = regexp.MustCompile(`(([,\.\-:;!\?&#+]|\n)[ \t\f\r]*|([ \t\f\r]+))`)
	spacePattern     = regexp.MustCompile(`\s+`)
)

// Transformation describes a widget state modification.
type Transformation int

// Widget transformations.
const (
	TransformFirstItem    Transformation = 1
	TransformLastItem     Transformation = 2
	TransformPreviousItem Transformation = 3
	TransformNextItem     Transformation = 4
	TransformPreviousPage Transformation = 5
	TransformNextPage     Transformation = 6
)

// Positions of substrings in regular expressions.
const (
	colorForegroundPos = 1
	colorBackgroundPos = 3
	colorFlagPos       = 5
)

// Predefined InputField acceptance functions.
var (
	// InputFieldInteger accepts integers.
	InputFieldInteger func(text string, ch rune) bool

	// InputFieldFloat accepts floating-point numbers.
	InputFieldFloat func(text string, ch rune) bool

	// InputFieldMaxLength returns an input field accept handler which accepts
	// input strings up to a given length. Use it like this:
	//
	//   inputField.SetAcceptanceFunc(InputFieldMaxLength(10)) // Accept up to 10 characters.
	InputFieldMaxLength func(maxLength int) func(text string, ch rune) bool
)

// Package initialization.
func init() {
	// Initialize the predefined input field handlers.
	InputFieldInteger = func(text string, ch rune) bool {
		if text == "-" {
			return true
		}
		_, err := strconv.Atoi(text)
		return err == nil
	}
	InputFieldFloat = func(text string, ch rune) bool {
		if text == "-" || text == "." || text == "-." {
			return true
		}
		_, err := strconv.ParseFloat(text, 64)
		return err == nil
	}
	InputFieldMaxLength = func(maxLength int) func(text string, ch rune) bool {
		return func(text string, ch rune) bool {
			return len([]rune(text)) <= maxLength
		}
	}
}

// Escape escapes the given text such that color and/or region tags are not
// recognized and substituted by the print functions of this package. For
// example, to include a tag-like string in a box title or in a TextView:
//
//   box.SetTitle(tview.Escape("[squarebrackets]"))
//   fmt.Fprint(textView, tview.Escape(`["quoted"]`))
func Escape(text string) string {
	return nonEscapePattern.ReplaceAllString(text, "$1[]")
}

// iterateString iterates through the given string one printed character at a
// time. For each such character, the callback function is called with the
// Unicode code points of the character (the first rune and any combining runes
// which may be nil if there aren't any), the starting position (in bytes)
// within the original string, its length in bytes, the screen position of the
// character, and the screen width of it. The iteration stops if the callback
// returns true. This function returns true if the iteration was stopped before
// the last character.
func iterateString(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	var screenPos int

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		r := gr.Runes()
		from, to := gr.Positions()
		width := stringWidth(gr.Str())
		var comb []rune
		if len(r) > 1 {
			comb = r[1:]
		}

		if callback(r[0], comb, from, to-from, screenPos, width) {
			return true
		}

		screenPos += width
	}

	return false
}

// iterateStringReverse iterates through the given string in reverse, starting
// from the end of the string, one printed character at a time. For each such
// character, the callback function is called with the Unicode code points of
// the character (the first rune and any combining runes which may be nil if
// there aren't any), the starting position (in bytes) within the original
// string, its length in bytes, the screen position of the character, and the
// screen width of it. The iteration stops if the callback returns true. This
// function returns true if the iteration was stopped before the last character.
func iterateStringReverse(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	type cluster struct {
		main                                       rune
		comb                                       []rune
		textPos, textWidth, screenPos, screenWidth int
	}

	// Create the grapheme clusters.
	var clusters []cluster
	iterateString(text, func(main rune, comb []rune, textPos int, textWidth int, screenPos int, screenWidth int) bool {
		clusters = append(clusters, cluster{
			main:        main,
			comb:        comb,
			textPos:     textPos,
			textWidth:   textWidth,
			screenPos:   screenPos,
			screenWidth: screenWidth,
		})
		return false
	})

	// Iterate in reverse.
	for index := len(clusters) - 1; index >= 0; index-- {
		if callback(
			clusters[index].main,
			clusters[index].comb,
			clusters[index].textPos,
			clusters[index].textWidth,
			clusters[index].screenPos,
			clusters[index].screenWidth,
		) {
			return true
		}
	}

	return false
}

// stripTags strips colour tags from the given string. (Region tags are not
// stripped.)
func stripTags(text string) string {
	stripped := colorPattern.ReplaceAllStringFunc(text, func(match string) string {
		if len(match) > 2 {
			return ""
		}
		return match
	})
	return escapePattern.ReplaceAllString(stripped, `[$1$2]`)
}
