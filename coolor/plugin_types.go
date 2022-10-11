package coolor

import (
	"encoding/json"
	"net/url"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/goutil/dump"
	"rogchap.com/v8go"

	"github.com/digitallyserviced/coolors/coolor/plugin"
)

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
type PluginEvent struct {
	eventType PluginEventType
	plugin    *Plugin
	name      string
}

type PluginData struct {
  ConfigKeys []string `json:"configKeys"`
  Handlers []string `json:"handlers"`
  Filenames []string `json:"filenames"`
  ConfigurationPaths []string `json:"configurationPaths"`
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
	PluginType                 uint
	PluginEventType            uint
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
  PluginWindBlows  = PluginInit|PluginModified|PluginBundled
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

func (v PluginDetectionType) String() string {
	return enumString(uint32(v), PluginDetectionTypes, false)
}

func (v PluginDetectionType) GoString() string {
	return enumString(uint32(v), PluginDetectionTypes, true)
}

func (v PluginType) String() string {
	return enumString(uint32(v), PluginTypes, false)
}

func (v PluginType) GoString() string {
	return enumString(uint32(v), PluginTypes, true)
}

func (v PluginEventType) String() string {
	return enumString(uint32(v), PluginEventTypes, false)
}

func (v PluginEventType) GoString() string {
	return enumString(uint32(v), PluginEventTypes, true)
}

func ModuleValueToExportsMap(ctx *v8go.Context, obj *v8go.Value) (pd *PluginData, err error){
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

  dump.P(pd)

	return
  // strs, err := V8StringArrayToGoStringArray(arr)
  // if checkErrX(err){
  //   for k, v := range a["default"] {
  //     
  //   }
  // }
}
