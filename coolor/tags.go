package coolor

import (
	"fmt"
	"reflect"

	// "strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

var TagTypes map[string]TagTypeInfo

type TagTypeFieldType uint32

const (
	FieldKey TagTypeFieldType = 1 << iota
	FieldYesNo
	FieldString
	FieldOptions
	FieldMultiple
	FieldRequired
	FieldDynamic

	MainTextField
	SecondaryTextField

	ListMainTextField      = FieldString | MainTextField
	ListSecondaryTextField = FieldString | SecondaryTextField

	RequiredFieldYesNo   = FieldYesNo | FieldRequired
	RequiredFieldString  = FieldString | FieldRequired
	RequiredFieldOptions = FieldOptions | FieldRequired
	shortcutChars        = "0123456789!@#$%^&*()"
)

var (
	TagKey         TagTypeField = NewTagTypeField("key", "tag key", reflect.TypeOf(""), FieldString)
	TagName        TagTypeField = NewTagTypeField("name", "tag name", reflect.TypeOf(""), ListMainTextField)
	TagDescription TagTypeField = NewTagTypeField("description", "tag description", reflect.TypeOf(""), ListSecondaryTextField)
	TagRequired    TagTypeField = NewTagTypeField("required", "field required", reflect.TypeOf(false), FieldRequired)
	options        []string     = []string{"one", "two"}
	TagOptions     TagTypeField = NewTagTypeField("options", options, reflect.TypeOf(options), FieldOptions)
	Base16Tags     TagTypeInfo
	StatusTags     TagTypeInfo
)

func init() {
	TagTypes = make(map[string]TagTypeInfo)
	Base16Tags = NewTagType("base16", "console base16 ansi colors", []TagTypeField{
		TagKey, TagName, TagDescription, TagRequired,
	})
}

func (f *TagTypeField) SetOptions(data interface{}) {
	if f.typed == reflect.TypeOf(data) {
		f.data = data
	}
}

type Tagger interface {
	AddTag(*TagItem)
	SetTag(*TagItem)
	SetTags([]*TagItem)
	ClearTags()
}

type Tagged interface {
	Tagger
	SetTagType(*TagTypeInfo)
	GetTag(idx int) *TagItem
	GetTags() []*TagItem
	GetTagsType() *TagTypeInfo
}

type Taggable struct {
	TagsType *TagTypeInfo
	Tags     []*TagItem
}

type TagTypeField struct {
	data     interface{}
	typed    reflect.Type
	name     string
	typeFlag TagTypeFieldType
}

type TagListItems struct {
	items []*TagItem
}

type TagList struct {
	*TagListItems
	*TagTypeInfo
  *ScriptShortcuts
	name string
}

type TagTypeFieldData struct {
	data interface{}
}
type TagItemData struct {
	data map[string]interface{}
}

type TagTypeInfo struct {
	name   string
	desc   string
	fields []TagTypeField
  tagList *TagList
}

type TagTypeCallbackFunc func(tti *TagTypeInfo, ti *TagItem)

type TagTypeCallback struct {
	callbacks map[string]*TagTypeCallbackFunc
}

type TagItem struct {
	*TagItemData
	*TagTypeInfo
  *ScriptShortcut
	idx int
}

// Changed implements ListItem
func (f TagItem) SecondaryText() string {
	ttf := f.GetFlag(SecondaryTextField)
  dump.P(ttf, ttf.name, f.data)
	return f.data[ttf.name].(string)
}
func (f TagItem) MainText() string {
	ttf := f.GetFlag(MainTextField)
	dump.P(ttf, MainTextField)
	return f.data[ttf.name].(string)
}
func (ti *TagItem) Shortcut() ScriptShortcut {
  return *ti.ScriptShortcut
}

func (*TagItem) Cancelled(idx int, i interface{}, lis []*ListItem) {
}

func (*TagItem) Changed(idx int, selected bool, i interface{}, lis []*ListItem) {
    dump.P(idx, selected, i, lis)
}

// Selected implements ListItem
func (ti *TagItem) Selected(idx int, i interface{}, lis []*ListItem) {

}

func (*TagItem) Visibility()ListItemsVisibility {
  return ListItemVisible
}


type TagTypeFieldValue struct {
	*TagTypeField
	*TagTypeFieldData
}

type TagListItem struct {
	main, secondary string
	shortcut        rune
}

func (tti *TagList) AddItem(ti *TagItem) {
	tti.items = append(tti.items, ti)
}

func (tti *TagList) AddTagItemWithData(args ...interface{}) *TagItem {
  ss := tti.ScriptShortcuts.TakeNext()
	ti := &TagItem{
		TagItemData:    tti.NewTagItemData(args...),
		TagTypeInfo:    tti.TagTypeInfo,
		ScriptShortcut: &ss,
		idx:            len(tti.items),
	}
	tti.items = append(tti.items, ti)
	return ti
}

func (lis *TagList) GetListItems() []*ListItem {
	lits := make([]*ListItem, 0)
	for _, v := range lis.items {
	litem := ListItem(v)
	var li *ListItem = &litem
		lits = append(lits, li)
	}
	return lits
}

func (tti *TagList) NewTagItem(tid *TagItemData) *TagItem {
  ss := tti.ScriptShortcuts.TakeNext()
	ti := &TagItem{
		TagItemData:    tid,
		TagTypeInfo:    tti.TagTypeInfo,
		ScriptShortcut: &ss,
		idx:            len(tti.items),
	}
	return ti
}

func (tti *TagTypeInfo) NewTagItemData(args ...interface{}) *TagItemData {
	// if len(tti.fields) != len(args) {
	//
	// }
	tid := &TagItemData{
		data: make(map[string]interface{}),
	}
	for i, v := range tti.fields {
		if i > len(args)-1 {
			break
		}
		tid.data[v.name] = args[i]
	}
	return tid
}

func (ttf TagTypeInfo) GetFlag(flag TagTypeFieldType) *TagTypeField {
	for _, v := range ttf.fields {
		if v.HasFlag(flag) {
			return &v
		}
	}
	return nil
}

func (ttf TagTypeField) HasFlag(flag TagTypeFieldType) bool {
	return flag&ttf.typeFlag != 0
}

func NewTagTypeField(name string, data interface{}, typed reflect.Type, flag TagTypeFieldType) TagTypeField {
	ttf := &TagTypeField{
		name:     name,
		typeFlag: flag,
		data:     data,
		typed:    typed,
	}
	return *ttf
}

func NewTagType(name, desc string, fields []TagTypeField) TagTypeInfo {
	tti := &TagTypeInfo{
		name:   name,
		desc:   desc,
		fields: fields,
	}
	TagTypes[name] = *tti
	return TagTypes[name]
}

func NewTagListItem(main, sec string, s rune) *TagListItem {
	ti := &TagListItem{
		main:      main,
		secondary: sec,
		shortcut:  s,
	}

	return ti
}

func (tti *TagTypeInfo) NewTagList(name string) *TagList {
	tl := &TagList{
		TagListItems:    &TagListItems{items: make([]*TagItem, 0)},
		TagTypeInfo:     tti,
		ScriptShortcuts: NewSubScriptShortcuts(),
		name:            name,
	}
  tti.tagList = tl
	return tl
}

// func (tl *TagList) AddTag(ti TagListItem) {
// 	tl.ListItems(ti)
// }

type TagsData struct {
	items []*TagItemData
}

func (lis *TagsData) GetItemCount() int {
	return len(lis.items)
}

func (lis *TagsData) GetItem(idx int) *TagItemData {
	if len(lis.items) > idx && idx > 0 {
		return lis.items[idx]
	}
	return nil
}

func (lis *TagsData) AddItem(li TagItemData) {
	lis.items = append(lis.items, &li)
}

func (lis *TagsData) GetTagsData() []*TagItemData {
	return lis.items
}

// func (f TagItem) Selected() {
// }
//
func (f TagItem) GetShortcut() rune {
	return rune(shortcutChars[f.idx])
}

func (f *TagListItem) GetMainTextStyle() tcell.Style {
	return tcell.Style{}
}

func (t *Taggable) init() (notNil bool) {
  if t == nil {
    return false
  }
	notNil = true
	if  t.Tags == nil {
		t.Tags = make([]*TagItem, 0)
		notNil = false
	}
	return notNil
}

func (t *Taggable) AddTag(ti *TagItem) {
	t.init()
	t.Tags = append(t.Tags, ti)
}

func (t *Taggable) SetTagType(tt *TagTypeInfo) {
	t.TagsType = tt
	t.init()
	t.ClearTags()
}

func (t *Taggable) SetTag(ti *TagItem) {
	t.init()
	// t.Tags = make([]*TagItem, 0)
	t.ClearTags()
	t.Tags = append(t.Tags, ti)
}

func (t *Taggable) SetTags(tis []*TagItem, appendItems bool) {
	if t.init() && !appendItems {
		t.Tags = make([]*TagItem, 0)
	}
	t.Tags = append(t.Tags, tis...)
}

func (t *Taggable) ClearTags() {
	t.init()
	t.Tags = make([]*TagItem, 0)
}

func (t *Taggable) GetTag(idx int) *TagItem {
	t.init()
	if len(t.Tags)-1 > idx {
		return t.Tags[idx]
	}
	return nil
}

func (t *Taggable) GetTagsType() *TagTypeInfo {
	if t.TagsType != nil {
		return t.TagsType
	}
	return nil
}

func (t *Taggable) GetTags() []*TagItem {
  if t == nil {
    return nil
  }
	t.init()
	if len(t.Tags) > 0 {
		return t.Tags
	}
	return nil
}

func (tl *TagList) GetTag(idx int) *TagItem {
	if idx < len(tl.items) {
		return tl.items[idx]
	}
	return nil
	// return *tl.GetItem(idx).(Tagged).GetTag(idx)
}

func NewTaggable(tti *TagTypeInfo) *Taggable {
  tgbl := &Taggable{
  	TagsType: tti,
  	Tags:     make([]*TagItem, 0),
  }
  return tgbl
}

func GetTerminalColorsAnsiTags() *TagList {
	items := Base16Tags.NewTagList("color tags")
	// items.ListItems =
	// it := Base16Tags.NewTagItemData("foreground", "fg", "default foreground", true)
	// Base16Tags.NewTagList
	// items.TagListItems
	// items.TagListItems.items = append(items.TagListItems.items, Base16Tags.TagListItems.NewTagItemWithData("foreground", "fg", "default foreground", true))
	items.AddTagItemWithData("fg", "foreground", "default foreground", true)
	items.AddTagItemWithData("bg", "background", "default background", true)
	items.AddTagItemWithData("cursor", "cursor", "cursor color", true)
	// items.AddItem(NewTagListItem("foreground", "default foreground", 'f'))
	// items.AddItem(NewTagListItem("background", "default background", 'b'))
	// items.AddItem(NewTagListItem("cursor", "cursor color", 'c'))
	ansiNames := make([]string, 0)
	for _, v := range baseAnsiNames {
		name := v
		ansiNames = append(ansiNames, name)
	}
	for _, v := range baseAnsiNames {
		name := fmt.Sprintf("%s %s", brightAnsiPrefix, v)
		ansiNames = append(ansiNames, name)
	}
	for i, name := range ansiNames {
		// s := rune(shortcutChars[i])
		// items.AddItem(NewTagListItem(name, fmt.Sprintf("%d/16 ansi color %s", i, name), s))
		// name := fmt.Sprintf("", i, name)
		// key := strings.Replace(name, " ", "_", -1)
		desc := fmt.Sprintf("4-bit color (%d) [%s]", i, name)
		items.AddTagItemWithData(name, name, desc, false)
	}
	return items
}
