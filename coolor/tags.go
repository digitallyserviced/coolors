package coolor

import (
	"fmt"
	"reflect"

	// "strings"
	"github.com/digitallyserviced/coolors/coolor/lister"
	"github.com/digitallyserviced/coolors/coolor/shortcuts"

	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
)

type (
  TagTypeFieldType uint32
)

const (
	FieldKey TagTypeFieldType = 1 << iota
	FieldYesNo
	FieldString
	FieldOptions
	FieldRequired
	FieldDynamic
	FieldMultiple
	FieldMutuallyExclusive

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
TagTypes map[string]TagType = make(map[string]TagType)
	TagKey         TagTypeField = NewTagTypeField("key", "tag key", reflect.TypeOf(""), FieldString|FieldKey)
	TagName        TagTypeField = NewTagTypeField("name", "tag name", reflect.TypeOf(""), ListMainTextField)
	TagDescription TagTypeField = NewTagTypeField("description", "tag description", reflect.TypeOf(""), ListSecondaryTextField)
	TagRequired    TagTypeField = NewTagTypeField("required", "field required", reflect.TypeOf(false), FieldRequired)
	TagDynamic    TagTypeField = NewTagTypeField("dynamic", "field dynamic", reflect.TypeOf(true), FieldDynamic|FieldRequired)
	options        []string     = []string{"one", "two"}
	TagOptions     TagTypeField = NewTagTypeField("options", options, reflect.TypeOf(options), FieldOptions)
	Base16Tags     TagType
	StatusTags     TagType
)

func init() {
	TagTypes = make(map[string]TagType)
	Base16Tags = NewTagTypeInfo("base16", "base16 color scheme tags", true, []TagTypeField{
		TagKey, TagName, TagDescription, TagRequired, TagDynamic,
	})
}

func (f *TagTypeField) SetOptions(data interface{}) {
	if f.typed == reflect.TypeOf(data) {
		f.data = data
	}
}

// type TaggableItems Tagger
// type Taggables struct {
//   *TaggableItems
// }

type Tagging interface {
	// AddTag(*TagItem)
  GetItem()
	SetTag(*TagItem)
	ClearTags()
}

type Tagger interface {
	Tagging
	SetTagType(*TagType)
	GetTagged(tag *TagItem)
	GetTags() []*TagItem
	SetTags([]*TagItem)
	GetTagsType() *TagType
}

type Tagged struct {
  TagType *TagType `msgpack:"-" clover:"-"`
	Item    interface{} `msgpack:"-" clover:"-"`
	Tags    []*TagItem `msgpack:"-" clover:"-"`
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
	*TagType
	*shortcuts.ScriptShortcuts
	name string
}

type TagTypeFieldData struct {
	data interface{}
}
type TagItemData struct {
	data map[string]interface{}
}

type TagType struct {
	tagList *TagList
	*TagTypeCallbacks
	name      string
	desc      string
	fields    []TagTypeField
	exclusive bool
}

type TagTypeCallbackFunc func(tti *TagType, ti *TagItem, tgd *Tagged)

type TagTypeCallbacks struct {
	callbacks map[string]TagTypeCallbackFunc
}

type TagItem struct {
	*TagItemData
	*TagType
	*shortcuts.ScriptShortcut
	idx int
}
type TagTypeFieldValue struct {
	*TagTypeField
	*TagTypeFieldData
}

type TagListItem struct {
	main, secondary string
	shortcut        rune
}

type TagsData struct {
	items []*TagItemData
}

func NewTaggable(tti *TagType) *Tagged {
	tgbl := &Tagged{
		TagType: tti,
		Tags:     make([]*TagItem, 0),
	}
	return tgbl
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

func NewTagTypeInfo(name, desc string, exclusive bool, fields []TagTypeField) TagType {
	tti := &TagType{
		TagTypeCallbacks: &TagTypeCallbacks{
			callbacks: make(map[string]TagTypeCallbackFunc),
		},
		name:             name,
		desc:             desc,
		fields:           fields,
		exclusive:        exclusive,
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
func (tti *TagType) NewTagList(name string) *TagList {
	tl := &TagList{
		TagListItems:    &TagListItems{items: make([]*TagItem, 0)},
		TagType:     tti,
		ScriptShortcuts: shortcuts.NewSubScriptShortcuts(),
		name:            name,
	}
	tti.tagList = tl
	return tl
}

func (tti *TagTypeCallbacks) SetCallback(name string, f TagTypeCallbackFunc) {
  tti.callbacks[name] = f
}

func (tti *TagTypeCallbacks) Callback(name string, tgd *Tagged, ti *TagItem) {
  cb, ok := tti.callbacks[name]
  if ok && cb != nil {
    cb(ti.TagType, ti, tgd)
  }
}

// Changed implements ListItem
func (f TagItem) SecondaryText() string {
	ttf := f.GetFlag(SecondaryTextField)
	// dump.P(ttf, ttf.name, f.data)
	return f.data[ttf.name].(string)
}

func (f TagItem) GetKey() string {
	ttf := f.GetFlag(FieldKey)
	// dump.P(ttf, MainTextField)
	return f.data[ttf.name].(string)
}

func (f TagItem) MainText() string {
	ttf := f.GetFlag(MainTextField)
	// dump.P(ttf, MainTextField)
	return f.data[ttf.name].(string)
}

func (ti *TagItem) Shortcut() shortcuts.ScriptShortcut {
	return *ti.ScriptShortcut
}

func (*TagItem) Cancelled(idx int, i interface{}, lis []*lister.ListItem) {
}

func (*TagItem) Changed(idx int, selected bool, i interface{}, lis []*lister.ListItem) {
	// dump.P(idx, selected, i, lis)
}

// Selected implements ListItem
func (ti *TagItem) Selected(idx int, i interface{}, lis []*lister.ListItem) {
}

func (*TagItem) Visibility() lister.ListItemsVisibility {
	return lister.ListItemVisible
}

func (tti *TagList) AddItem(ti *TagItem) {
	tti.items = append(tti.items, ti)
}

func (tti *TagList) AddTagItemWithData(args ...interface{}) *TagItem {
	ss := tti.TakeNext()
	ti := &TagItem{
		TagItemData:    tti.NewTagItemData(args...),
		TagType:    tti.TagType,
		ScriptShortcut: &ss,
		idx:            len(tti.items),
	}
	tti.items = append(tti.items, ti)
	return ti
}

func (lis *TagList) GetListItems() []*lister.ListItem {
	lits := make([]*lister.ListItem, 0)
	for _, v := range lis.items {
		litem := lister.ListItem(v)
		li := &litem
		lits = append(lits, li)
	}
	return lits
}

func (tti *TagList) NewTagItem(tid *TagItemData) *TagItem {
	ss := tti.TakeNext()
	ti := &TagItem{
		TagItemData:    tid,
		TagType:    tti.TagType,
		ScriptShortcut: &ss,
		idx:            len(tti.items),
	}
	return ti
}

func (tti *TagType) NewTagItemData(args ...interface{}) *TagItemData {
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

func (ttf TagType) GetFlag(flag TagTypeFieldType) *TagTypeField {
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


// func (tl *TagList) AddTag(ti TagListItem) {
// 	tl.ListItems(ti)
// }
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

func (t *Tagged) init() (notNil bool) {
	if t == nil {
		return false
	}
	notNil = true
	if t.Tags == nil {
		t.Tags = make([]*TagItem, 0)
		notNil = false
	}
	return notNil
}

func (t *Tagged) AddTag(ti *TagItem) {
	t.init()
	t.Tags = append(t.Tags, ti)
}

func (t *Tagged) SetTagType(tt *TagType) {
	t.TagType = tt
	t.init()
	t.ClearTags()
}

func (t *Tagged) SetTag(ti *TagItem) {
	t.init()
	t.ClearTags()
	t.Tags = append(t.Tags, ti)
  t.TagType.Callback("set", t, ti)
}

func (t *Tagged) SetTags(tis []*TagItem, appendItems bool) {
	if t.init() && !appendItems {
		t.Tags = make([]*TagItem, 0)
	}
	t.Tags = append(t.Tags, tis...)
}

func (t *Tagged) ClearTags() {
	t.init()
	t.Tags = make([]*TagItem, 0)
}

func (t *Tagged) GetTag(idx int) *TagItem {
	t.init()
  if t == nil {
    return nil
  }
  if len(t.Tags) == 0 {
    return nil
  }
	if len(t.Tags)-1 >= idx {
		return t.Tags[idx]
	}
	return nil
}

func (t *Tagged) GetTagsType() *TagType {
	if t.TagType != nil {
		return t.TagType
	}
	return nil
}

func (t *Tagged) GetTags() []*TagItem {
	if t == nil {
		return nil
	}
	t.init()
	if len(t.Tags) > 0 {
		return t.Tags
	}
	return nil
}

func (tl *TagList) GetTagBy(tagField string) *TagItem {
  for _, v := range tl.items {
    if v.GetKey() == tagField {
      return v
    }
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

func GetTerminalColorsAnsiTags() *TagList {
	items := Base16Tags.NewTagList("base16 color scheme")
  // Base16Tags.SetCallback("set", func(tti *TagType, ti *TagItem, tgd *Tagged) {
  //
  // })
	items.AddTagItemWithData("fg", "foreground", "default foreground", true, true)
	items.AddTagItemWithData("bg", "background", "default background", true, true)
	items.AddTagItemWithData("cursor", "cursor", "cursor color", true, true)
	for i, name := range baseXtermAnsiColorNames {
		desc := fmt.Sprintf("4-bit color (%d) [%s]", i, name)
		items.AddTagItemWithData(name, name, desc, true, false)
	}
	return items
}
