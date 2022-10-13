package coolor

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/goutil/dump"
	"github.com/knadh/koanf"
	"go.uber.org/zap/zapcore"
	"rogchap.com/v8go"

	"github.com/digitallyserviced/coolors/coolor/plugin"
	"github.com/digitallyserviced/coolors/coolor/util"
)

type PluginSchemeMeta struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	OriginUrl string `json:"origin_url"`
}

type PluginSchemeFile struct {
	Plugin *Plugin
	PluginSchemeMeta
	file *os.File
	Data *koanf.Koanf
  ConfigData map[string]interface{}
}

type PluginFileHandler struct {
	Import *v8go.Function
	Export *v8go.Function
	TagMap *v8go.Object
	Misc   *v8go.Object
}

type Plugin struct {
	Name, Path    string
	OriginUrl     url.URL
	PluginType    PluginType
	manager       *PluginsManager
	monitor       <-chan PluginEvent
	DetectionType PluginDetectionType
	EventTypes    PluginEventType
	*PluginData
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (psm *PluginSchemeMeta) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString("scheme.name", psm.Name)
	oe.AddString("scheme.author", psm.Author)
	return nil
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (psf *PluginSchemeFile) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddObject("plugin", psf.Plugin)
	if psf.Data != nil {
		oe.AddReflected("plugin.data", psf.Data.All())
	}
	return nil
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (pe *PluginEvent) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddObject("plugin", pe.plugin)
	oe.AddString("plugin.event", pe.eventType.String())
	return nil
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (p *Plugin) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddString("plugin.name", p.Name)
	oe.AddString("plugin.path", p.Path)
	return nil
}

type PluginEvent struct {
	eventType PluginEventType
	plugin    *Plugin
	name      string
}

type PluginData struct {
	ConfigKeys         []string `json:"configKeys"`
	Handlers           []string `json:"handlers"`
	Filenames          []string `json:"filenames"`
	ConfigurationPaths []string `json:"configurationPaths"`
	Decoder            []string `json:"decoder"`
}

// MarshalLogObject implements zapcore.ObjectMarshaler
func (pd *PluginData) MarshalLogObject(oe zapcore.ObjectEncoder) error {
  oe.AddReflected("configKeys", pd.ConfigKeys)
  return nil
}

type PluginsManager struct {
	watcher   fsnotify.Watcher
	fsmonitor <-chan interface{}
	gv8       *plugin.GoV8Env
	bundler   chan PluginEvent
	start     chan struct{}
	done      chan struct{}
	cancel    chan struct{}
	Plugins   []*Plugin
	monitors  []chan PluginEvent
}

const (
	pluginsPath        = "js/plugins"
	bundledPluginsPath = "js/plugins/.bundled"
)

type (
	PluginType          uint
	PluginEventType     uint
	PluginDetectionType uint
)

const (
	ColorSchemeImportPlugin PluginType = 1 << iota
	ColorSchemeExportPlugin
	ColorSchemeDetectPlugin
	ColorSchemeTaggedPlugin

	ColorModPlugin
	ColorPalettePlugin
	ColorPaletteGeneratorPlugin

	PreviewPlugin
	FilePlugin

	CommandPlugin
	ActionPlugin

	GlobalPlugin
	HookPlugin

	PlaygroundPlugin

	FullColorSchemePlugin = ColorSchemeImportPlugin | ColorSchemeExportPlugin | ColorSchemeDetectPlugin | ColorSchemeTaggedPlugin
)

const (
	FilenameDetection PluginDetectionType = 1 << iota
	ConfigKeysDetection
	RegexDetection
	FunctionDetection
)

const (
	PluginInit PluginEventType = 1 << iota
	PluginModified
	PluginBundled
	PluginWindBlows = PluginInit | PluginModified | PluginBundled
)

var (
// PluginLogMarshaler zzlog.ObjectMarshalerFunc = func(oe zapcore.ObjectEncoder) error {
//
// }
)

var PluginTypes = []enumName{
	{uint32(ColorSchemeImportPlugin), "ColorSchemeImportPlugin"},
	{uint32(ColorSchemeExportPlugin), "ColorSchemeExportPlugin"},
	{uint32(ColorSchemeDetectPlugin), "ColorSchemeDetectPlugin"},
	{uint32(ColorSchemeTaggedPlugin), "ColorSchemeTaggedPlugin"},

	{uint32(ColorModPlugin), "ColorModPlugin"},
	{uint32(ColorPalettePlugin), "ColorPalettePlugin"},
	{uint32(ColorPaletteGeneratorPlugin), "ColorPaletteGeneratorPlugin"},

	{uint32(PreviewPlugin), "PreviewPlugin"},
	{uint32(FilePlugin), "FilePlugin"},

	{uint32(CommandPlugin), "CommandPlugin"},
	{uint32(ActionPlugin), "ActionPlugin"},

	{uint32(GlobalPlugin), "GlobalPlugin"},
	{uint32(HookPlugin), "HookPlugin"},

	{uint32(PlaygroundPlugin), "PlaygroundPlugin"},

	{uint32(FullColorSchemePlugin), "FullColorSchemePlugin"},
}

var PluginEventTypes = []enumName{
	{uint32(PluginInit), "PluginInit"},
	{uint32(PluginBundled), "PluginBundled"},
	{uint32(PluginModified), "PluginModified"},
	{uint32(PluginWindBlows), "PluginWindBlows"},
}

var PluginDetectionTypes = []enumName{
	{uint32(FilenameDetection), "FilenameDetection"},
	{uint32(ConfigKeysDetection), "ConfigKeysDetection"},
	{uint32(RegexDetection), "RegexDetection"},
	{uint32(FunctionDetection), "FunctionDetection"},
}

func (a PluginDetectionType) Is(b PluginDetectionType) bool {
	return util.BitAnd(a, b)
}
func (v PluginDetectionType) String() string {
	return enumString(uint32(v), PluginDetectionTypes, false)
}
func (v PluginDetectionType) GoString() string {
	return enumString(uint32(v), PluginDetectionTypes, true)
}

func (a PluginType) Is(b PluginType) bool {
	return util.BitAnd(a, b)
}
func (v PluginType) String() string {
	return enumString(uint32(v), PluginTypes, false)
}
func (v PluginType) GoString() string {
	return enumString(uint32(v), PluginTypes, true)
}

func (a PluginEventType) Is(b PluginEventType) bool {
	return util.BitAnd(a, b)
}
func (v PluginEventType) String() string {
	return enumString(uint32(v), PluginEventTypes, false)
}
func (v PluginEventType) GoString() string {
	return enumString(uint32(v), PluginEventTypes, true)
}

func SetupContexts(ctx *v8go.Context, iso *v8go.Isolate) {
	obTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginTypes {
		obTpl.Set(v.s, v.v, v8go.None)
	}
	peTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginEventTypes {
		peTpl.Set(v.s, v.v, v8go.None)
	}
	pdTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginDetectionTypes {
		pdTpl.Set(v.s, v.v, v8go.None)
	}
	if err := ctx.Global().Set("pluginEventTypes", eajs(peTpl.NewInstance(ctx))); err != nil {
		panic(err)
	}
	if err := ctx.Global().Set("pluginDetectionTypes", eajs(pdTpl.NewInstance(ctx))); err != nil {
		panic(err)
	}
	if err := ctx.Global().Set("pluginTypes", eajs(obTpl.NewInstance(ctx))); err != nil {
		panic(err)
	}
	if err := ctx.Global().Set("xtermAnsiNames", eajs(plugin.GoStructToV8Object(ctx, baseXtermAnsiColorNames))); err != nil {
		panic(err)
	}

	// ansiNames, _ := GoStructToV8Object(gv8.Ctx, baseXtermAnsiColorNames)
	// gv8.Ctx.Global().Set("xtermAnsiNames", ansiNames)
}

func ModuleValueToExportsMap(ctx *v8go.Context, obj *v8go.Value) (pd *PluginData, err error) {
	pd = &PluginData{}
	var exports *v8go.Object
	if !obj.IsObject() {
		return
	}
	if obj.IsObject() && eajs(obj.AsObject()).Has("default") {
		exports = eajs(eajs(eajs(obj.AsObject()).Get("default")).AsObject())
	} else {
		exports = eajs(obj.AsObject())
	}
	expJson, err := v8go.JSONStringify(ctx, exports)
	dump.P(expJson)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(expJson), pd)
	if err != nil {
		return
	}

	return
	// strs, err := V8StringArrayToGoStringArray(arr)
	// if checkErrX(err){
	//   for k, v := range a["default"] {
	//
	//   }
	// }
}
