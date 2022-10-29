package coolor

import (
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf"
	"go.uber.org/zap/zapcore"

	"github.com/digitallyserviced/coolors/coolor/plugin"
	"github.com/digitallyserviced/coolors/coolor/util"
	. "github.com/digitallyserviced/coolors/coolor/events"
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
	Name, Path    string
	OriginUrl     url.URL
	PluginType    PluginType
	manager       *PluginsManager
	monitor       <-chan PluginEvent
	DetectionType PluginDetectionType
	EventTypes    PluginEventType
	KeyMapTypes   PluginKeyMapType
	Worker        *PluginWorker
	*PluginData
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
	Plugins   []*Plugin
	monitors  []chan PluginEvent
	*EventNotifier
	*EventObserver
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

var PluginEventTypes = []EnumName{
	{uint32(PluginInit), "PluginInit"},
	{uint32(PluginModified), "PluginModified"},
	{uint32(PluginReloaded), "PluginReloaded"},
	{uint32(PluginBundled), "PluginBundled"},
	{uint32(PluginScanConfigPaths), "PluginScanConfigPaths"},
	{uint32(PluginScanConfigResult), "PluginScanConfigResult"},
	{uint32(PluginCreatedKeyMap), "PluginCreatedKeyMap"},
	{uint32(PluginWindBlows), "PluginWindBlows"},
}

var PluginKeyMapTypes = []EnumName{
	{uint32(PluginTagsKeyMap), "PluginTagsKeyMap"},
	{uint32(PluginMetaKeyMap), "PluginMetaKeyMap"},
	{uint32(PluginExtrasKeyMap), "PluginExtrasKeyMap"},
	{uint32(PluginColorsKeyMap), "PluginColorsKeyMap"},
	{uint32(PluginFeaturesKeyMap), "PluginFeaturesKeyMap"},
}

var PluginDetectionTypes = []EnumName{
	{uint32(FilenameDetection), "FilenameDetection"},
	{uint32(ConfigKeysDetection), "ConfigKeysDetection"},
	{uint32(RegexDetection), "RegexDetection"},
	{uint32(FunctionDetection), "FunctionDetection"},
}

func (a PluginDetectionType) Is(b PluginDetectionType) bool {
	return util.BitAnd(a, b)
}
func (v PluginDetectionType) String() string {
	return EnumString(uint32(v), PluginDetectionTypes, false)
}
func (v PluginDetectionType) GoString() string {
	return EnumString(uint32(v), PluginDetectionTypes, true)
}

func (a PluginType) Is(b PluginType) bool {
	return util.BitAnd(a, b)
}
func (v PluginType) String() string {
	return EnumString(uint32(v), PluginTypes, false)
}
func (v PluginType) GoString() string {
	return EnumString(uint32(v), PluginTypes, true)
}

func (a PluginKeyMapType) Is(b PluginKeyMapType) bool {
	return util.BitAnd(a, b)
}
func (v PluginKeyMapType) String() string {
	return EnumString(uint32(v), PluginKeyMapTypes, false)
}
func (v PluginKeyMapType) GoString() string {
	return EnumString(uint32(v), PluginKeyMapTypes, true)
}

func (a PluginEventType) Is(b PluginEventType) bool {
	return util.BitAnd(a, b)
}
func (v PluginEventType) String() string {
	return EnumString(uint32(v), PluginEventTypes, false)
}
func (v PluginEventType) GoString() string {
	return EnumString(uint32(v), PluginEventTypes, true)
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
