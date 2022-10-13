package coolor

import (
	"fmt"
	"io/ioutil"
	"os"

	// "log"
	"net/url"
	"path/filepath"
	"time"

	// "path/filepath"
	// "strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf"
	// "github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"

	// "github.com/gookit/goutil/dump"
	// "github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/fsutil/finder"
	"github.com/samber/lo"
	"rogchap.com/v8go"

	// "github.com/digitallyserviced/coolors/coolor/zzlog"
	"github.com/digitallyserviced/tview"

	"github.com/digitallyserviced/coolors/coolor/plugin"
	"github.com/digitallyserviced/coolors/coolor/util"
	"github.com/digitallyserviced/coolors/coolor/zzlog"
	"github.com/digitallyserviced/coolors/tree"
)

var PluginManager *PluginsManager

func NewPluginEvent(ev PluginEventType, name string, p *Plugin) *PluginEvent {
	pe := &PluginEvent{
		eventType: ev,
		plugin:    p,
		name:      name,
	}
	zlog.Debug(
		"new_pe",
		zzlog.Reflect("plugin", p.String()),
		zzlog.Reflect("event", ev.String()),
	)
	// fmt.Printf("new_pe: %+v\n", pe)
	return pe
}

// Name description
func InitPlugins() error {
	PluginManager = NewPluginsManager()
	go PluginManager.StartPluginBundler()
	go PluginManager.StartPluginMonitor()

	return nil
}

func eajs[R any](v R, e error) R {
	iserr := func(v interface{}) {
		// if v == nil {
		// 	return
		// }
		e, ok := v.(*v8go.JSError)
		if ok {
			zlog.Error(
				fmt.Sprintf("%T %v", v, e),
				zzlog.String("msg", e.Message),
				zzlog.String("loc", e.Location),
			)
			// doLog(errorx.Newf("%s %s %s %v", e.StackTrace, e.Message, e.Location, v))
			panic(e)
		}
	}
	iserr(e)
	return v
}

func NewPluginsManager() *PluginsManager {
	watcher, events := plugin.Watcher()
	done := make(chan struct{})

	pm := &PluginsManager{
		watcher:   *watcher,
		fsmonitor: events,
		monitors:  make([]chan PluginEvent, 0),
		bundler:   make(chan PluginEvent),
		start:     make(chan struct{}),
		done:      done,
		cancel:    make(chan struct{}),
		Plugins:   make([]*Plugin, 0),
		gv8: plugin.NewGoV8(func(gv8 *plugin.GoV8Env) error {
			gv8.VM = v8go.NewIsolate()
			gv8.Ctx = v8go.NewContext(gv8.VM)
			return nil
		}),
	}

	pm.StartPluginBundler()
	return pm
}

type LogFields []zzlog.Field

// struct{
//   Plugin LogFields
//   Manager LogFields
//   Watcher LogFields
// } =
var PluginLogFields = struct {
	Plugin  LogFields
	Manager LogFields
	Watcher LogFields
}{
	Plugin:  LogFields{zzlog.String("from", "plugins")},
	Manager: LogFields{zzlog.String("from", "pluginmanager")},
	Watcher: LogFields{zzlog.String("from", "fsnotify")},
}

func (lf LogFields) With(args ...zzlog.Field) LogFields {
	return append(lf, args...)
}

func (p *Plugin) String() string {
	str := fmt.Sprintf(
		"%s (%s)\nType[%s]\nDetection[%s]\nEvents[%s]",
		p.Name,
		p.Path,
		p.PluginType.String(),
		p.DetectionType.String(),
		p.EventTypes.String(),
	)
	return str
}

func (p *Plugin) HandleEvent(e PluginEvent) error {
	// p.manager.gv8.Ctx.RunScript(source string, origin string)
	return nil
}

func (p *Plugin) StartMonitor() error {
	go func() {
		tick := time.NewTicker(time.Millisecond * 100)
		defer tick.Stop()
		for {
			select {
			case e := <-p.monitor:
				zlog.Info(
					"recv pe monitor event",
					PluginLogFields.Watcher.With(
						zzlog.String("plugin", p.Name),
						zzlog.String("path", p.Path),
						zzlog.String("eventType", e.eventType.String()),
						zzlog.String("dir", "before"),
					)...)
				switch {
				case e.eventType&PluginBundled != 0:
					p.LoadBundle()
				case e.eventType == PluginWindBlows:
					fallthrough
				case e.eventType&PluginModified != 0:
					fallthrough
				case e.eventType&PluginInit != 0:
					// zlog.Info("Monitor", PluginLogFields.Watcher.With(zzlog.String("eventType", e.eventType.String()), zzlog.String("dir", "before"))...)
					// fmt.Printf("monitor_plugin_init: %s %s %+v\n", p.Name, p.Path, e)

					p.manager.bundler <- e
				}
			case <-tick.C:
			}
		}
	}()
	return nil
}
func (p *Plugin) Bundle() (string, api.BuildResult) {
	// outfile := p.getBundledPath()
	entry := ErrorAssert(filepath.Abs(p.getPluginPath()))
	outfile := ErrorAssert(filepath.Abs(p.getBundledPath()))
	// defer os.Remove(outfile)

	result := api.Build(api.BuildOptions{
		AbsWorkingDir: ErrorAssert(filepath.Abs(filepath.Dir(entry))),
		// NodePaths:         []string{"js/node_modules"},
		EntryPoints: []string{entry},
		Bundle:      true,
		Format:      api.FormatESModule,
		Outfile:     outfile,
		Write:       true,
		GlobalName:  "global",
		// MinifyWhitespace:  true,
		// MinifyIdentifiers: true,
		// MinifySyntax:      true,
	})
	if len(result.Errors) > 0 {
		for _, v := range result.Errors {
			fmt.Printf("errs: %s %v\n", v)
		}
		// fmt.Println(result.Errors)
	}
	if len(result.Warnings) > 0 {
		for _, v := range result.Warnings {
			fmt.Printf("warns: %s %v\n", v)
		}
		// fmt.Println(result.Errors)
	}

	if len(result.Errors) != 0 {
		return "", result
	}

	if bytes, err := ioutil.ReadFile(outfile); err != nil {
		return "", result
	} else {
		return string(bytes), result
	}
}

func (p *Plugin) getFilenameFilters() []fsutil.FilterFunc {
	filters := make([]fsutil.FilterFunc, 0)
	for _, v := range p.Filenames {
		filters = append(filters, func(fPath string, fi os.FileInfo) bool {
			return ErrorAssert(filepath.Match(v, fPath))
		})
	}
	return filters
}

func (p *Plugin) FindColorSchemes(path string) []*tview.TreeNode {
	schemes := make([]*tview.TreeNode, 0)
	filters := p.getFilenameFilters()
	fsutil.FindInDir(
		fsutil.Expand(path),
		func(fPath string, fi os.FileInfo) error {
			if !fi.IsDir() {
				node := tree.NewNode(fPath, fi)
				schemes = append(schemes, node)
			}
			return nil
		},
		filters...)
	return schemes
}

func (p *Plugin) findPluginInterests() []*tview.TreeNode {
	nodes := make([]*tview.TreeNode, 0)

	if len(p.PluginData.ConfigurationPaths) > 0 {
		expanded := lo.Map[string, string](
			p.PluginData.ConfigurationPaths,
			func(s string, i int) string {
				return fsutil.Expand(s)
			},
		)
		f := finder.NewFinder(expanded)
		f.ExcludeDotDir(false).ExcludeDotFile(false)
		f.AddFilter(finder.GlobFilterFunc(p.PluginData.Filenames, true))
		// f.Each(func(fi os.FileInfo, filePath string) {
		//
		// })
		f.EachFile(func(sf *os.File) {
      psf := p.ParseScheme(sf)
      if psf.ConfigData != nil {

      }
		})

		// for _, v := range p.PluginData.ConfigurationPaths {
		//   // finder.NewFinder(dirPaths []string, filePaths ...string)
		//   nodes = append(nodes, p.FindColorSchemes(v)...)
		// }
	}

	return nodes
}

// func (p *Plugin) GetDecoder(psf *PluginSchemeFile) {
//
// }

func (pd *PluginData) GetDecoder() koanf.Parser {
  for _, v := range pd.Decoder {
    dec, ok := plugin.PluginSchemeFileDecoders[v]
    if ok {
      return dec
    }
  }
  return nil
}

func (p *Plugin) ParseScheme(f *os.File) *PluginSchemeFile {
  psf := p.LoadSchemeFile(f)
  return psf
}

func (p *Plugin) LoadSchemeFile(f *os.File) *PluginSchemeFile {
	defer f.Close()
	psf := &PluginSchemeFile{
		Plugin: p,
		PluginSchemeMeta: PluginSchemeMeta{
			Name:      "",
			Author:    "",
			OriginUrl: "",
		},
		file: f,
		Data: koanf.New("."),

	}
	parser := p.PluginData.GetDecoder()
	if parser == nil {
  zlog.Info("no decoder", PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
		return psf
	}
	err := psf.Data.Load(
		rawbytes.Provider(ErrorAssert(ioutil.ReadAll(f))),
		parser,
	)
	if !checkErrX(err) {
		return nil
	}

  if p.HasMatchingKeys(psf) {
    zlog.Info("scheme file detected", PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
    psf.ConfigData = psf.Data.All()
  }
	return psf
}

func (p *Plugin) HasMatchingKeys(psf *PluginSchemeFile) bool {
	hasKeys := false
	if p.DetectionType.Is(ConfigKeysDetection) {
    zlog.Info("check config keys", PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf), zzlog.Object("scheme_file", p.PluginData))...)
		for _, k := range p.PluginData.ConfigKeys {
			hasKeys = psf.Data.Exists(k)
		}
	} else {
		hasKeys = true
	}
	return hasKeys
}
func (p *Plugin) getPluginPath() string {
	return filepath.Join(p.Path, "index.js")
}
func (p *Plugin) getMetaPath() string {
	return filepath.Join(p.Path, "meta.js")
}
func (p *Plugin) getBundledPath() string {
	return filepath.Join(p.Path, ".bundled", "index.js")
}

func (p *Plugin) LoadBundle() error {
	gov8 := plugin.NewGoV8(func(gv8 *plugin.GoV8Env) error {
		gv8.VM = v8go.NewIsolate()
		gv8.Ctx = v8go.NewContext(gv8.VM)
		SetupContexts(gv8.Ctx, gv8.VM)
		return nil
	})
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				eajs([]string{}, e)
				zlog.Error(
					"js script recovered",
					PluginLogFields.Plugin.With(
						zzlog.String("plugin", p.String()),
					)...)

			}
		}
	}()
	// SetupContexts(gov8.Ctx, gov8.VM)
	gov8.DoBindings()
	zlog.Info(
		"load bundle",
		PluginLogFields.Plugin.With(
			zzlog.String("plugin", p.Name),
			zzlog.String("path", p.Path),
			zzlog.String("bundle_path", p.getBundledPath()),
			zzlog.String("plugin_path", p.getPluginPath()),
			zzlog.String(
				"abs_bplugin_path",
				ErrorAssert(filepath.Abs(p.getBundledPath())),
			),
		)...)
	modSrc := ErrorAssert(ioutil.ReadFile(p.getBundledPath()))
	val := eajs(
		gov8.Ctx.RunModule(string(modSrc), p.getPluginPath()),
	)
	pd := ErrorAssert(ModuleValueToExportsMap(gov8.Ctx, val))
	p.PluginData = pd
	zlog.Info(
		"load bundle val",
		PluginLogFields.Plugin.With(
			zzlog.String("plugin", p.String()),
			zzlog.Reflect("dump", pd),
		)...)

	return nil
}

func (p *Plugin) LoadMeta() error {
	path := p.Path
	metaSrc, err := ioutil.ReadFile(filepath.Join(path, "meta.js"))
	if !checkErrX(err, path) {
		return err
	}

	SetupContexts(p.manager.gv8.Ctx, p.manager.gv8.VM)

	val := eajs(p.manager.gv8.Ctx.RunModule(string(metaSrc), "meta.js"))

	if !val.IsModuleNamespaceObject() {
		return errorx.Newf(
			"Plugin %s returned an unexpected value %+v",
			val.DetailString(),
		)
	}

	exp, _ := val.AsObject()
	meta, _ := exp.Get("meta")
	if !meta.IsObject() {
		return errorx.Newf(
			"Plugin meta module did not return expected values got: %+v %+v",
			exp,
			meta,
		)
	}

	metaObj := eajs(meta.AsObject())
	// dump.P(JSErrorAssert(metaObj.Get("shit")).DetailString())
	name := eajs(metaObj.Get("name")).DetailString()
	originUrl := eajs(metaObj.Get("originUrl")).DetailString()
	pluginType := eajs(metaObj.Get("pluginType")).Uint32()
	detectionType := eajs(metaObj.Get("detection")).Uint32()
	eventType := eajs(metaObj.Get("events")).Uint32()
	p.Name = name
	p.PluginType = PluginType(pluginType)
	p.DetectionType = PluginDetectionType(detectionType)
	p.EventTypes = PluginEventType(eventType)
	url := eajs(url.Parse(originUrl))
	p.OriginUrl = *url
	zlog.Info(
		"load meta val",
		PluginLogFields.Plugin.With(
			zzlog.String("plugin", p.String()),
			zzlog.Reflect("dump", p),
		)...)

	p.manager.gv8.VM.Dispose()

	return nil
}

// func (pm *PluginsManager) DetectColorSchemeType() error {
//
// }

func (pm *PluginsManager) StartPluginBundler() error {
	go func() {
		for v := range pm.bundler {
			if v.eventType&PluginModified != 0 {
				_, result := v.plugin.Bundle()
				_ = result
				p, m := pm.getPluginByPath(v.plugin.Path)
				util.LossySend(
					make(<-chan struct{}),
					m,
					*NewPluginEvent(PluginBundled, "BOO", p),
					time.Second*1,
				)
				// m <-
				// _ = p
				// _ = src
				// plugin.Bundle(filepath.Join(v.plugin.Path, "index.js"))
			}
		}
	}()
	return nil
}

func (pm *PluginsManager) SupportedFilenames() []string {
	if pm == nil {
		return nil
	}
	return lo.Flatten[string](
		lo.Map[*Plugin, []string](pm.Plugins, func(p *Plugin, i int) []string {
			if p.PluginData == nil {
				return []string{}
			}
      names := p.PluginData.Filenames
			return 	names	
    }),
	)
}

func (pm *PluginsManager) GetTreeEntries() func(f *tview.TreeNode) []*tview.TreeNode {
	return func(f *tview.TreeNode) []*tview.TreeNode {
		nodes := make([]*tview.TreeNode, 0)
		for _, v := range PluginManager.Plugins {
			vn := tree.NewVirtualNode(v.Name, "", "")
			vn.Node.SetChildren(v.findPluginInterests())
			nodes = append(nodes, vn.Node)
		}
		return nodes
	}
}

func (pm *PluginsManager) Each(
	f func(p *Plugin, pm *PluginsManager, idx int) error,
) {
	lo.ForEach[*Plugin](pm.Plugins, func(p *Plugin, i int) {
		fmt.Printf("each: %v %v\n", p, i)
		err := f(p, pm, i)
		if err != nil {
			fmt.Printf("each: %v %v\n", p, err)
		}
	})
}

func (pm *PluginsManager) getPluginByPath(
	path string,
) (plug *Plugin, monitor chan<- PluginEvent) {
	pm.Each(func(p *Plugin, pm *PluginsManager, idx int) error {
		if ErrorAssert(filepath.Match(fmt.Sprintf("*%s*", path), p.Path)) {
			plug = p
			monitor = pm.monitors[idx]
		}
		return nil
	})
	return
}

func (pm *PluginsManager) HandleWatcherEvent(e fsnotify.Event) {
	p, mon := pm.getPluginByPath(e.Name)
	zlog.Debug(
		"watcherevent",
		PluginLogFields.Watcher.With(
			zzlog.String("watcherpath", e.Name),
			zzlog.String("fsevent", e.Op.String()),
			zzlog.String("plugin", p.String()),
		)...)
	ev := NewPluginEvent(
		PluginModified,
		fmt.Sprintf("FSNotify %s ", e.Op.String()),
		p,
	)
	mon <- *ev
	// util.LossySend[T any](done <-chan struct{}, valueChan chan<- util.T, value T, t time.Duration)
}

func (pm *PluginsManager) Loaded(path string) bool {
	if lo.ContainsBy[*Plugin](pm.Plugins, func(p *Plugin) bool {
		return p.Path == path
	}) {
		return true
	}

	return false
}

func (pm *PluginsManager) Watching() finder.FilterFunc {
	watchingPaths := pm.watcher.WatchList()
	watchingPaths = append(
		watchingPaths,
		lo.Map[*Plugin, string](pm.Plugins, func(p *Plugin, i int) string {
			return fmt.Sprintf("*%s*", p.Path)
		})...)
	return finder.FilterFunc(finder.GlobFilterFunc(watchingPaths, false))
}

func (pm *PluginsManager) Scan(path string) error {
	jsff := finder.EmptyFinder()
	jsff.AddDir(pluginsPath).
		AddDirFilter(finder.DirNameFilterFunc([]string{".bundle"}, false)).
		AddFileFilter(finder.SuffixFilterFunc([]string{"meta.js", "index.js"}, true))
	jsff.Each(func(filePath string) {
		zlog.Info("scanner found", zzlog.String("found", filePath))
		if ErrorAssert(filepath.Match("js/plugins/*/meta.js", filePath)) {
			pm.InitPlugin(filepath.Dir(filePath))
		}
	})
	return nil
}

func (pm *PluginsManager) StartPluginMonitor() error {
	pm.Scan("js/plugins")

	// pluginEvents := util.TakeFn(pm.done, pm.fsmonitor, func(i interface{}) bool {
	//    return true
	// })
	//
	//
	debouncedChan := debounce(
		50*time.Millisecond,
		200*time.Millisecond,
		pm.fsmonitor,
		// util.TakeFn(pm.done, pm.fsmonitor, func(i interface{}) bool {
		//     event := i.(fsnotify.Event)
		//     // fmt.Println(filepath.Dir(event.Name))
		//     p, m := pm.getPluginByPath(filepath.Dir(event.Name))
		//     // fmt.Println(p, m)
		//     if p == nil {
		//       pm.Scan("js/plugins")
		//     } else {
		//       e := NewPluginEvent(PluginModified, "fsmods", p)
		//       // fmt.Printf("fsmod: %s %v\n", p.Path, e)
		//       m <- *e
		//     }
		//     return true
		// }),
	)

	go func() {
		for {
			// fmt.Printf("deb: \n")
			select {
			case i := <-debouncedChan:
				event := i.(fsnotify.Event)
				// fmt.Println(filepath.Dir(event.Name))
				p, m := pm.getPluginByPath(filepath.Dir(event.Name))
				// fmt.Println(p, m)
				if p == nil {
					pm.Scan("js/plugins")
				} else {
					e := NewPluginEvent(PluginModified, "fsmods", p)
					m <- *e
				}
			}
		}

	}()
	// <-pluginEvents
	//
	// watcher.Add("js/lib")
	// watcher.Add("js/lib/schemes")
	return nil
}

func (pm *PluginsManager) InitPlugin(ppath string) error {
	fmt.Println(ppath)
	if pm.Loaded(ppath) {
		p, m := pm.getPluginByPath(ppath)
		m <- *NewPluginEvent(PluginModified, "modded", p)
		return nil
	}

	monitor := make(chan PluginEvent)

	p := &Plugin{
		Name:       "",
		Path:       ppath,
		OriginUrl:  url.URL{},
		PluginType: 0,
		PluginData: &PluginData{
			ConfigKeys:         make([]string, 0),
			Handlers:           make([]string, 0),
			Filenames:          make([]string, 0),
			ConfigurationPaths: make([]string, 0),
			Decoder:            make([]string, 0),
		},
		manager: pm,
		monitor: monitor,
	}

	pm.monitors = append(pm.monitors, monitor)

	p.LoadMeta()
	pm.watcher.Add(ppath)
	p.StartMonitor()

	pm.Plugins = append(pm.Plugins, p)

	p, m := pm.getPluginByPath(ppath)
	m <- *NewPluginEvent(PluginModified, "modded", p)

	return nil
}

func (pm *PluginsManager) Init() error {
	//  return nil
	return nil
}

// func Handle
