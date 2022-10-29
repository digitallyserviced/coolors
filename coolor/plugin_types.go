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
	{uint64(ColorSchemeImportPlugin), "ColorSchemeImportPlugin"},
	{uint64(ColorSchemeExportPlugin), "ColorSchemeExportPlugin"},
	{uint64(ColorSchemeDetectPlugin), "ColorSchemeDetectPlugin"},
	{uint64(ColorSchemeTaggedPlugin), "ColorSchemeTaggedPlugin"},

	{uint64(ColorModPlugin), "ColorModPlugin"},
	{uint64(ColorPalettePlugin), "ColorPalettePlugin"},
	{uint64(ColorPaletteGeneratorPlugin), "ColorPaletteGeneratorPlugin"},

	{uint64(PreviewPlugin), "PreviewPlugin"},
	{uint64(FilePlugin), "FilePlugin"},

	{uint64(CommandPlugin), "CommandPlugin"},
	{uint64(ActionPlugin), "ActionPlugin"},

	{uint64(GlobalPlugin), "GlobalPlugin"},
	{uint64(HookPlugin), "HookPlugin"},

	{uint64(PlaygroundPlugin), "PlaygroundPlugin"},

	{uint64(FullColorSchemePlugin), "FullColorSchemePlugin"},
}

var PluginEventTypes = []EnumName{
	{uint64(PluginInit), "PluginInit"},
	{uint64(PluginModified), "PluginModified"},
	{uint64(PluginReloaded), "PluginReloaded"},
	{uint64(PluginBundled), "PluginBundled"},
	{uint64(PluginScanConfigPaths), "PluginScanConfigPaths"},
	{uint64(PluginScanConfigResult), "PluginScanConfigResult"},
	{uint64(PluginCreatedKeyMap), "PluginCreatedKeyMap"},
	{uint64(PluginWindBlows), "PluginWindBlows"},
}

var PluginKeyMapTypes = []EnumName{
	{uint64(PluginTagsKeyMap), "PluginTagsKeyMap"},
	{uint64(PluginMetaKeyMap), "PluginMetaKeyMap"},
	{uint64(PluginExtrasKeyMap), "PluginExtrasKeyMap"},
	{uint64(PluginColorsKeyMap), "PluginColorsKeyMap"},
	{uint64(PluginFeaturesKeyMap), "PluginFeaturesKeyMap"},
}

var PluginDetectionTypes = []EnumName{
	{uint64(FilenameDetection), "FilenameDetection"},
	{uint64(ConfigKeysDetection), "ConfigKeysDetection"},
	{uint64(RegexDetection), "RegexDetection"},
	{uint64(FunctionDetection), "FunctionDetection"},
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
