package coolor

import (
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf"
	"go.uber.org/zap/zapcore"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/plugin"
	"github.com/digitallyserviced/coolors/coolor/util"
)

type (
	PluginType          uint
	PluginEventType     uint
	PluginEventStatus   uint
	PluginWorkerStatus  uint
	PluginDetectionType uint
	PluginKeyMapType    uint
)

type PluginSchemeMeta struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	OriginUrl string `json:"origin_url"`
}

type PluginSchemeFile struct {
	Plugin *Plugin
	PluginSchemeMeta
	file       *os.File
	Data       *koanf.Koanf
	ConfigData map[string]interface{}
	Palette    map[string]string
	buf        []byte
}

type Plugin struct {
	manager *PluginsManager
	monitor <-chan PluginEvent
	Worker  *PluginWorker
	*PluginData
	OriginUrl     url.URL
	Name          string
	Path          string
	PluginType    PluginType
	DetectionType PluginDetectionType
	EventTypes    PluginEventType
	KeyMapTypes   PluginKeyMapType
}
type PluginData struct {
	ConfigKeys         []string `json:"configKeys"`
	Handlers           []string `json:"handlers"`
	KeyMapTypes        []string `json:"keyMapHandlers"`
	Filenames          []string `json:"filenames"`
	ConfigurationPaths []string `json:"configurationPaths"`
	Decoder            []string `json:"decoder"`
}

type PluginsManager struct {
	watcher   fsnotify.Watcher
	fsmonitor <-chan interface{}
	gv8       *plugin.GoV8Env
	bundler   chan PluginEvent
	global    chan PluginEvent
	start     chan struct{}
	done      chan struct{}
	cancel    chan struct{}
	*EventNotifier
	*EventObserver
	Plugins  []*Plugin
	monitors []chan PluginEvent
}

// GetRef implements Referenced
func (pm *PluginsManager) GetRef() interface{} {
  return pm
}

const (
	pluginsPath        = "js/plugins"
	bundledPluginsPath = "js/plugins/.bundled"
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

  ColorSchemeFilesPlugin = ColorSchemeImportPlugin | ColorSchemeDetectPlugin
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
	PluginReloaded
	PluginBundled
	PluginScanConfigPaths
	PluginScanConfigResult
	PluginCreatedKeyMap
	PluginWindBlows = PluginInit | PluginModified | PluginReloaded | PluginBundled | PluginScanConfigPaths | PluginScanConfigResult | PluginCreatedKeyMap
)

const (
	PluginTagsKeyMap PluginKeyMapType = 1 << iota
	PluginMetaKeyMap
	PluginExtrasKeyMap
	PluginColorsKeyMap
	PluginFeaturesKeyMap
)

const (
	PluginWorkerNop PluginWorkerStatus = 1 << iota
	PluginWorkerIdle
	PluginWorkerDone
	PluginWorkerReset
	PluginWorkerStartup
	PluginWorkerBusy
)
const (
	PluginEventNop PluginEventStatus = 1 << iota
	PluginEventInit
	PluginEventEmitted
	PluginEventReceived
	PluginEventHandled
	PluginEventFailure
)

var (
// PluginLogMarshaler zzlog.ObjectMarshalerFunc = func(oe zapcore.ObjectEncoder) error {
//
// }
)

var PluginTypes = []EnumName{
	{V: uint64(ColorSchemeImportPlugin), S: "ColorSchemeImportPlugin"},
	{V: uint64(ColorSchemeExportPlugin), S: "ColorSchemeExportPlugin"},
	{V: uint64(ColorSchemeDetectPlugin), S: "ColorSchemeDetectPlugin"},
	{V: uint64(ColorSchemeTaggedPlugin), S: "ColorSchemeTaggedPlugin"},

	{V: uint64(ColorModPlugin), S: "ColorModPlugin"},
	{V: uint64(ColorPalettePlugin), S: "ColorPalettePlugin"},
	{V: uint64(ColorPaletteGeneratorPlugin), S: "ColorPaletteGeneratorPlugin"},

	{V: uint64(PreviewPlugin), S: "PreviewPlugin"},
	{V: uint64(FilePlugin), S: "FilePlugin"},

	{V: uint64(CommandPlugin), S: "CommandPlugin"},
	{V: uint64(ActionPlugin), S: "ActionPlugin"},

	{V: uint64(GlobalPlugin), S: "GlobalPlugin"},
	{V: uint64(HookPlugin), S: "HookPlugin"},

	{V: uint64(PlaygroundPlugin), S: "PlaygroundPlugin"},

	{V: uint64(ColorSchemeFilesPlugin), S: "ColorSchemeFilesPlugin"},
	{V: uint64(FullColorSchemePlugin), S: "FullColorSchemePlugin"},
}

var PluginEventTypes = []EnumName{
	{V: uint64(PluginInit), S: "PluginInit"},
	{V: uint64(PluginModified), S: "PluginModified"},
	{V: uint64(PluginReloaded), S: "PluginReloaded"},
	{V: uint64(PluginBundled), S: "PluginBundled"},
	{V: uint64(PluginScanConfigPaths), S: "PluginScanConfigPaths"},
	{V: uint64(PluginScanConfigResult), S: "PluginScanConfigResult"},
	{V: uint64(PluginCreatedKeyMap), S: "PluginCreatedKeyMap"},
	{V: uint64(PluginWindBlows), S: "PluginWindBlows"},
}

var PluginKeyMapTypes = []EnumName{
	{V: uint64(PluginTagsKeyMap), S: "PluginTagsKeyMap"},
	{V: uint64(PluginMetaKeyMap), S: "PluginMetaKeyMap"},
	{V: uint64(PluginExtrasKeyMap), S: "PluginExtrasKeyMap"},
	{V: uint64(PluginColorsKeyMap), S: "PluginColorsKeyMap"},
	{V: uint64(PluginFeaturesKeyMap), S: "PluginFeaturesKeyMap"},
}

var PluginDetectionTypes = []EnumName{
	{V: uint64(FilenameDetection), S: "FilenameDetection"},
	{V: uint64(ConfigKeysDetection), S: "ConfigKeysDetection"},
	{V: uint64(RegexDetection), S: "RegexDetection"},
	{V: uint64(FunctionDetection), S: "FunctionDetection"},
}

func (a PluginDetectionType) Is(b PluginDetectionType) bool {
	return util.BitAnd(a, b)
}
func (v PluginDetectionType) String() string {
	return EnumString(uint64(v), PluginDetectionTypes, false)
}
func (v PluginDetectionType) GoString() string {
	return EnumString(uint64(v), PluginDetectionTypes, true)
}

func (a PluginType) Is(b PluginType) bool {
	return util.BitAnd(a, b)
}
func (v PluginType) String() string {
	return EnumString(uint64(v), PluginTypes, false)
}
func (v PluginType) GoString() string {
	return EnumString(uint64(v), PluginTypes, true)
}

func (a PluginKeyMapType) Is(b PluginKeyMapType) bool {
	return util.BitAnd(a, b)
}
func (v PluginKeyMapType) String() string {
	return EnumString(uint64(v), PluginKeyMapTypes, false)
}
func (v PluginKeyMapType) GoString() string {
	return EnumString(uint64(v), PluginKeyMapTypes, true)
}

func (a PluginEventType) Is(b PluginEventType) bool {
	return util.BitAnd(a, b)
}
func (v PluginEventType) String() string {
	return EnumString(uint64(v), PluginEventTypes, false)
}
func (v PluginEventType) GoString() string {
	return EnumString(uint64(v), PluginEventTypes, true)
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

// MarshalLogObject implements zapcore.ObjectMarshaler
func (pd *PluginData) MarshalLogObject(oe zapcore.ObjectEncoder) error {
	oe.AddReflected("configKeys", pd.ConfigKeys)
	return nil
}
