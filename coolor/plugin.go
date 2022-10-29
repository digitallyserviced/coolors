package coolor

import (
	"container/ring"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"net/url"
	"path/filepath"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/gdamore/tcell/v2"
	"github.com/knadh/koanf"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/fsutil/finder"
	"github.com/knadh/koanf/providers/rawbytes"

	"github.com/samber/lo"
	"rogchap.com/v8go"

	"github.com/digitallyserviced/tview"

	// "github.com/digitallyserviced/coolors/coolor/events"
	. "github.com/digitallyserviced/coolors/coolor/events"
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
		zzlog.Reflect("plugin", p),
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

func eajs[R any](val R, e error) R {
	iserr := func(v interface{}) {
		if v == nil {
			return
		}
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
		err, ok := v.(error)
		if ok {
			zlog.Error(
				fmt.Sprintf("%T %v", v, err),
				zzlog.String("msg", err.Error()),
			)
			// doLog(errorx.Newf("%s %s %s %v", e.StackTrace, e.Message, e.Location, v))
			panic(err)
		}
	}
	iserr(e)
	return val
}

var pmOnLoads []PluginManagerOnLoad
func addPluginOnLoad(name string, pmol PluginManagerOnLoad){
  if pmOnLoads == nil {
    pmOnLoads = make([]PluginManagerOnLoad, 0) 
  }
  pmOnLoads = append(pmOnLoads, pmol)
}
type PluginManagerOnLoad func(pm *PluginsManager)

func (pmol *PluginManagerOnLoad) Exec(pm *PluginsManager)  {
  
}

func NewPluginsManager() *PluginsManager {
	watcher, events := plugin.Watcher()
	done := make(chan struct{})

	pm := &PluginsManager{
		watcher:   *watcher,
		fsmonitor: events,
		gv8: plugin.NewGoV8(func(gv8 *plugin.GoV8Env) error {
			gv8.VM = v8go.NewIsolate()
			gv8.Ctx = v8go.NewContext(gv8.VM)
			return nil
		}),
		bundler:       make(chan PluginEvent),
		global:       make(chan PluginEvent),
		start:         make(chan struct{}),
		done:          done,
		cancel:        make(chan struct{}),
		Plugins:       make([]*Plugin, 0),
		monitors:      make([]chan PluginEvent, 0),
		EventNotifier: NewEventNotifier("plugins_manager"),
		EventObserver: NewEventObserver("plugins_manager"),
	}

	pm.StartPluginBundler()

  for _, v := range pmOnLoads {
    v(pm)
  }

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

func (pw *PluginWorker) Work() {
	rtTimeout := time.NewTicker(workerTimeout)
	for {
		zlog.Info(
			"work loop",
			PluginLogFields.Plugin.With(zzlog.Object("plugin", pw.Plugin))...)
		select {
		case <-pw.done:
			return
    case w := <-pw.Plugin.monitor:
      pw.events <- &w
		case w := <-pw.events:
			zlog.Info(
				"worker event",
				PluginLogFields.Watcher.With(
					zzlog.Object("plugin", pw.Plugin),
					zzlog.String("event", w.eventType.String()),
				)...)
			rtTimeout.Reset(workerTimeout)
			switch {
			case w.eventType.Is(PluginCreatedKeyMap):
			case w.eventType.Is(PluginScanConfigPaths):
				w.plugin.ScanConfigPaths()
			case w.eventType.Is(PluginScanConfigResult):
			case w.eventType.Is(PluginBundled):
				err := pw.Plugin.LoadMeta()
				if err != nil {
					panic(err)
				}
			case w.eventType.Is(PluginInit), w.eventType.Is(PluginModified):
				pw.Plugin.manager.bundler <- *w
			}
		case <-rtTimeout.C:
			// pw.Plugin.ScanConfigPaths()
			zlog.Info(
				"timeout",
				PluginLogFields.Plugin.With(zzlog.Object("plugin", pw.Plugin))...)
		}
	}
}

var workerTimeout = 5000 * time.Millisecond

func (pw *PluginWorkerSandbox) DisposeRuntime() {
	if pw != nil && pw.gv8 != nil {
		if pw.gv8.Ctx != nil {
			pw.gv8.Ctx.Close()
		}
		if pw.gv8.VM != nil {
			pw.gv8.VM.Dispose()
		}
		pw.valid = false
	}
}

func (pw *PluginWorker) ResetSandbox() {
	pw.PluginWorkerSandbox.DisposeRuntime()
	pw.PluginWorkerSandbox.data = make(map[string]interface{})
}

func (pw *PluginWorker) Kill(s PluginWorkerStatus) {
	pw.once.Do(func() {
		pw.DisposeRuntime()
		pw.once = nil
	})
}

func (pw *PluginWorker) UpdateStatus(s PluginWorkerStatus) {
	go func() {
		pw.status <- s
	}()
}

func (pw *PluginWorker) LoadAndWork() {
	pw.once = &sync.Once{}
	err := pw.Plugin.LoadMeta()
	zlog.Info(
		"load and work",
		PluginLogFields.Plugin.With(zzlog.Object("plugin", pw.Plugin))...)
	if err != nil {
		panic(err)
	}
	pw.Plugin.LoadBundle(
		func(p *Plugin, gov8 *plugin.GoV8Env, exports *v8go.Value, handlers *v8go.Value) {
			pd := ErrorAssert(ModuleValueToExportsMap(gov8.Ctx, exports))
			p.PluginData = pd
			pw.PluginWorkerSandbox.gv8 = gov8
			go pw.Work()
		},
	)
}

func (p *Plugin) NewWorker(
	f func(p *Plugin),
) *PluginWorker {
	pw := &PluginWorker{
		idx:                 nextWorkerIdx(),
		Plugin:              p,
		PluginWorkerSandbox: &PluginWorkerSandbox{gv8: &plugin.GoV8Env{}},
		status:              make(chan PluginWorkerStatus, 1),
		statusHistory:       ring.New(8),
		events:              make(chan *PluginEvent, MAX_BUFFERED_EVENTS),
		done:                make(chan struct{}),
		onLoad:              func(p *Plugin, psf *PluginSchemeFile) {},
	}
	p.LoadMeta()

	p.Worker = pw

	go func() {
		p.Worker.LoadAndWork()
	}()

	return pw
}

func (p *Plugin) Bundle() (string, api.BuildResult) {
	// outfile := p.getBundledPath()
	entry := ErrorAssert(filepath.Abs(p.getPluginPath()))
	handlers := ErrorAssert(filepath.Abs(p.getHandlersPath()))
	outfile := ErrorAssert(filepath.Abs(p.getBundledPath()))
	outdir := filepath.Dir(ErrorAssert(filepath.Abs(p.getBundledPath())))
	// defer os.Remove(outfile)

	result := api.Build(api.BuildOptions{
		AbsWorkingDir: ErrorAssert(filepath.Abs(filepath.Dir(entry))),
		// NodePaths:         []string{"js/node_modules"},
		EntryPoints: []string{entry, handlers},
		Bundle:      true,
		Format:      api.FormatESModule,
		// Outfile:     outfile,
		Write:      true,
		Outdir:     outdir,
		GlobalName: "global",
		// MinifyWhitespace:  true,
		// MinifyIdentifiers: true,
		// MinifySyntax:      true,
	})
	if len(result.Errors) > 0 {
		for _, v := range result.Errors {
			zzlog.Error(
				"bundle err",
				PluginLogFields.Manager.With(zzlog.Reflect("bundleinfo", v))...)
			fmt.Printf("errs: %s %v\n", v)
		}
		// fmt.Println(result.Errors)
	}
	if len(result.Warnings) > 0 {
		for _, v := range result.Warnings {
			zzlog.Warn(
				"bundle warn",
				PluginLogFields.Manager.With(zzlog.Reflect("bundleinfo", v))...)
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

func (p *Plugin) ScanConfigPaths() []*tview.TreeNode {
	nodes := make([]*tview.TreeNode, 0)

	if len(p.PluginData.ConfigurationPaths) > 0 || true {
		expanded := lo.Map[string, string](
			p.PluginData.ConfigurationPaths,
			func(s string, i int) string {
				return fsutil.Expand(s)
			},
		)
		zlog.Info(
			"expande",
			PluginLogFields.Plugin.With(zzlog.Reflect("expanded", expanded))...)
		f := finder.NewFinder(expanded)
		f.ExcludeDotDir(false).ExcludeDotFile(false)
		f.AddFilter(finder.GlobFilterFunc(p.PluginData.Filenames, true))
		f.EachFile(func(sf *os.File) {
			zlog.Info(
				"pre parse scheme",
				PluginLogFields.Plugin.With(zzlog.Object("scheme_file", p))...)

			psf, valid := p.ParseScheme(sf)
      if !valid {
        zlog.Warn(
          "scheme not detected",
          PluginLogFields.Plugin.With(zzlog.Object("scheme_file", p))...)
        return
      }
      zlog.Debug(fmt.Sprintf("%s %s emitting result", psf.Name, psf.Author), PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
      pe := p.NewPluginEvent(psf.Name,PluginScanConfigResult,  psf)
      p.EmitEvent(pe)
		})
	}

	return nodes
}

func (p *Plugin) NewPluginEvent(name string, pet PluginEventType, refs... interface{}) *PluginEvent {
  pe := &PluginEvent{
  	eventType: pet,
  	plugin:    p,
  	name:      name,
  	refs:      refs,
  	callback: func(pe *PluginEvent, gv8 *plugin.GoV8Env, args... interface{}) error {
      return nil
  	},
  	status: PluginEventInit,
  }
  zlog.Debug(fmt.Sprintf("new event %s %s", name, pet.String()), PluginLogFields.Plugin.With(zzlog.Reflect("refs", refs))...)
  return pe
}

func (pe *PluginEvent) Status(pes PluginEventStatus) *PluginEvent {
  pe.status = pes
  return pe
}

func (p *Plugin) EmitEvent(pe *PluginEvent) *PluginEvent {
  p.manager.DispatchEvent(pe)
  return pe
}

func (pd *PluginData) GetDecoder() koanf.Parser {
	for _, v := range pd.Decoder {
		dec, ok := plugin.PluginSchemeFileDecoders[v]
		if ok {
			return dec
		}
	}
	return nil
}

func (p *Plugin) SetupKeyMaps(psf *PluginSchemeFile) {
	for _, v := range PluginKeyMapTypes {
		if p.KeyMapTypes.Is(PluginKeyMapType(v.V)) {

		}
	}
}

func (p *Plugin) HandleDetection(psf *PluginSchemeFile) bool {
	if p.DetectionType.Is(ConfigKeysDetection) || true {
		zlog.Info(
			fmt.Sprintf("%s check file", psf.file.Name()),
			PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
		if hasAll, existingKeys := p.CheckKeys(psf); !hasAll {
			zlog.Warn(
				fmt.Sprintf(
					"%s scheme file does not have wanted keys",
					psf.file.Name(),
				),
				PluginLogFields.Plugin.With(
					zzlog.Object("scheme_file", psf),
					zzlog.Reflect("existing_keys", existingKeys),
				)...)
		} else {
			zlog.Info(
				fmt.Sprintf("%s scheme file has wanted keys", psf.file.Name()),
				PluginLogFields.Plugin.With(
					zzlog.Object("scheme_file", psf),
					zzlog.Reflect("existing_keys", existingKeys),
				)...)
			psf.Plugin = p
			return true
		}
	}
	return false
}

func (p *Plugin) ParseScheme(f *os.File) (psf *PluginSchemeFile, valid bool){
	psf = p.GetSchemeFile(f)
  valid = false
	if p.HandleDetection(psf) {
    valid = true
		zlog.Info(
			"scheme file detected",
			PluginLogFields.Plugin.With(
				zzlog.Object("scheme_file", psf),
				zzlog.Reflect("configdata", psf.Data.All()),
			)...)
		p.LoadBundle(
			func(p *Plugin, gov8 *plugin.GoV8Env, exports *v8go.Value, handlers *v8go.Value) {
				defaults := eajs(ModuleExportsFromDefaults(gov8.Ctx, exports))

				if !defaults.Has("keyMapHandlers") {
					return
				}
				gov8.Ctx.Global().
					Set("configData", eajs(plugin.GoStructToV8Object(gov8.Ctx, psf.ConfigData)))

				if fn, err := handlers.Object().Get("getMetaMapped"); err == nil &&
					fn.IsFunction() {
					confData := eajs(plugin.GoStructToV8Object(gov8.Ctx, psf.ConfigData))
					mapd, er := fn.Function().Call(gov8.Ctx.Global(), confData)
					if er != nil {
						dump.P("err", er.Error())
						panic(er)
					}
					psf.Name = eajs(mapd.Object().Get("name")).DetailString()
          if psf.Name == "" || psf.Name == "undefined" {
            psf.Name =filepath.Base(psf.file.Name()) 
          }
					psf.Author = eajs(mapd.Object().Get("author")).DetailString()
					psf.OriginUrl = eajs(mapd.Object().Get("name")).DetailString()
				} else {
					panic(err)
				}
				if fn, err := handlers.Object().Get("getTagMapped"); err == nil &&
					fn.IsFunction() {
					confData := eajs(plugin.GoStructToV8Object(gov8.Ctx, psf.ConfigData))
					mapd, er := fn.Function().Call(gov8.Ctx.Global(), confData)
					if er != nil {
						dump.P("err", er.Error())
						panic(er)
					}
					tlist := GetTerminalColorsAnsiTags()
					for _, v := range tlist.items {
						k := v.GetKey()
						psf.Palette[k] = eajs(mapd.Object().Get(k)).DetailString()
					}
				} else {
					panic(err)
				}
		zlog.Info(
			"scheme file loaded",
			PluginLogFields.Plugin.With(
				zzlog.Object("scheme_file", psf),
				zzlog.Reflect("configdata", psf.Data.All()),
			)...)
			},
		)
    return
	}

	return
}

func (psf *PluginSchemeFile) GetPalette() *CoolorColorsPalette {
	cp := NewCoolorColorsPalette()
	for tagName, col := range psf.Palette {
		tag := cp.tagType.tagList.GetTagBy(tagName)
		c := cp.AddCoolorColor(NewCoolorColor(col))
		c.SetTag(tag)
	}

	return cp
}

func (cp *PluginSchemeFile) TagsKeys(random bool) CoolorPaletteTagsMeta {
	tagKeys := make(map[string]*Coolor)
	cptm := &CoolorPaletteTagsMeta{
		tagCount:     0,
		TaggedColors: tagKeys,
	}
	tlist := GetTerminalColorsAnsiTags()
	if Base16Tags.tagList == nil {
	}

	for _, v := range tlist.items {
		cptm.tagCount += 1
		// k := v.data[keyfield.name].(string)
		k := v.GetKey()
		tagKeys[k] = nil // cp.RandomColor().Coolor()
	}
	// return *cptm

	for k, v := range cp.Palette {
		cptm.TaggedColors[k] = &Coolor{
			Color: tcell.GetColor(v),
		}
		cptm.tagCount += 1
	}

	return *cptm
}

func (psf *PluginSchemeFile) Decode(
	onLoad func(p *Plugin, psf *PluginSchemeFile),
) *PluginSchemeFile {
	parser := psf.Plugin.PluginData.GetDecoder()
	// dump.P(parser, psf.file.Name())
	if parser == nil {
		zlog.Warn(
			"no decoder",
			PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
		return psf
	}
	err := psf.Data.Load(
		rawbytes.Provider(psf.buf),
		parser,
	)

	if !checkErrX(err) {
		zlog.Error(
			"err loading file ",
			PluginLogFields.Plugin.With(zzlog.Object("scheme_file", psf))...)
		return psf
	}

	psf.ConfigData = psf.Data.All()

	if onLoad != nil {
		onLoad(psf.Plugin, psf)
	}

	return psf
}
func (p *Plugin) GetSchemeFile(
	f *os.File,
	// onLoad func(p *Plugin, psf *PluginSchemeFile, f *os.File),
) *PluginSchemeFile {
	defer f.Close()

	psf := &PluginSchemeFile{
		Plugin:           p,
		PluginSchemeMeta: PluginSchemeMeta{Name: "", Author: "", OriginUrl: ""},
		file:             f,
		Data:             koanf.New("."),
		ConfigData:       map[string]interface{}{},
		Palette:          make(map[string]string),
		buf:              ErrorAssert(ioutil.ReadAll(f)),
	}

	psf.Decode(func(p *Plugin, psf *PluginSchemeFile) {
		zlog.Info(
			"decoding",
			PluginLogFields.Plugin.With(
				zzlog.Object("scheme_file", psf),
				zzlog.Reflect("configdata", psf.Data.All()),
			)...)
		// dump.P(psf.ConfigData)
	})
	psf.file.Close()
	psf.buf = []byte{}

	return psf
}

func (p *Plugin) CheckKeys(
	psf *PluginSchemeFile,
) (hasAll bool, existsKeys map[string]bool) {
	hasAll = false
	existsKeys = make(map[string]bool)
	zlog.Info(
		"checking config keys",
		PluginLogFields.Plugin.With(zzlog.Object("scheme_file", p.PluginData))...)

	for _, k := range p.PluginData.ConfigKeys {
		hasAll = psf.Data.Exists(k)
		existsKeys[k] = psf.Data.Exists(k)
	}
	zlog.Info(
		"finished checking keys",
		PluginLogFields.Plugin.With(
			zzlog.Object("scheme_file", p.PluginData),
			zzlog.Reflect("existing_keys", existsKeys),
		)...)
	return
}

func (pm *Plugin) GetPluginRuntime(
	i func(gv8 *plugin.GoV8Env) error,
) *plugin.GoV8Env {
	gov8 := pm.manager.GetRuntime(func(gv8 *plugin.GoV8Env) error {
		if i != nil {
			i(gv8)
		}
		return nil
	})
	return gov8
}

func (p *Plugin) LoadBundle(
	callbacks ...func(p *Plugin, gov8 *plugin.GoV8Env, exports *v8go.Value, handlers *v8go.Value),
) error {
	gov8 := p.manager.GetRuntime(func(gv8 *plugin.GoV8Env) error {
		InjectBaseObjects(gv8.Ctx, gv8.VM)
		return nil
	})

	defer gov8.VM.Dispose()
	defer gov8.Ctx.Close()

	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				fmt.Println(e)
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
	handlerSrc := ErrorAssert(ioutil.ReadFile(p.getBundledHandlersPath()))
	val := eajs(
		gov8.Ctx.RunModule(string(modSrc), p.getPluginPath()),
	)
	hdlrs := eajs(
		gov8.Ctx.RunModule(string(handlerSrc), p.getMetaPath()),
	)

	if len(callbacks) > 0 {
		for _, cb := range callbacks {
			cb(p, gov8, val, hdlrs)
		}
	}

	return nil
}

func (p *Plugin) InjectHandlerData(
	gov8 *plugin.GoV8Env,
	psf *PluginSchemeFile,
) {

}

func (p *Plugin) LoadMeta() error {
	path := p.Path
	metaSrc, err := ioutil.ReadFile(filepath.Join(path, "meta.js"))
	if !checkErrX(err, path) {
		return err
	}

	gov8 := p.manager.GetRuntime(func(gv8 *plugin.GoV8Env) error {
		InjectBaseObjects(gv8.Ctx, gv8.VM)
		return nil
	})

	defer gov8.VM.Dispose()
	defer gov8.Ctx.Close()

	val := eajs(gov8.Ctx.RunModule(string(metaSrc), "meta.js"))

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
	pluginType := eajs(metaObj.Get("pluginType")).BigInt().Uint64()
	detectionType := eajs(metaObj.Get("detection")).BigInt().Uint64()
	eventType := eajs(metaObj.Get("events")).BigInt().Uint64()
	p.Name = name
	p.PluginType = PluginType(pluginType)
	p.DetectionType = PluginDetectionType(detectionType)
	p.EventTypes = PluginEventType(eventType)
	url := eajs(url.Parse(originUrl))
	p.OriginUrl = *url

	return nil
}

func (pm *PluginsManager) GetRuntime(
	i func(gv8 *plugin.GoV8Env) error,
) *plugin.GoV8Env {
	gov8 := plugin.NewGoV8(func(gv8 *plugin.GoV8Env) error {
		gv8.VM = v8go.NewIsolate()
		gv8.Ctx = v8go.NewContext(gv8.VM)
		if i != nil {
			i(gv8)
		}
		return nil
	})
	return gov8
}

func (pm *PluginsManager) StartPluginBundler() error {
	go func() {
		for v := range pm.bundler {
			if v.eventType&PluginModified != 0 || true {
				_, result := v.plugin.Bundle()
				_ = result
				p, m := pm.getPluginByPath(v.plugin.Path)
				util.LossySend(
					make(<-chan struct{}),
					m,
					*NewPluginEvent(PluginBundled, "BOO", p),
					time.Second*1,
				)
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
			return names
		}),
	)
}

func (pm *PluginsManager) GetTreeEntries() func(f *tview.TreeNode) []*tview.TreeNode {
	return func(f *tview.TreeNode) []*tview.TreeNode {
		nodes := make([]*tview.TreeNode, 0)
		for _, v := range PluginManager.Plugins {
			vn := tree.NewPluginNode(v.Name, "", v)
      vn.Node.SetChildren(v.ScanConfigPaths())
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

func (pm *PluginsManager) DispatchEvent(
	pe *PluginEvent,
) error {
  if pe.plugin == nil {
    for i, p := range pm.Plugins {
      zlog.Debug(
        "dispatch event",
        zzlog.Reflect("plugin", p.String()),
        zzlog.Reflect("event", pe),
      )
      pm.monitors[i] <- *NewPluginEvent(pe.eventType, pe.name, p)
    }
  } else {
      // zlog.Debug(
      //   "dispatch global event",
      //   zzlog.Reflect("event", pe),
      // )
    pm.global <- *pe
      // zlog.Debug(
      //   "post dispatch global event",
      //   zzlog.Reflect("event", pe),
      // )
  }

	return nil
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

// func (p *Plugin) StartMonitor() error {
// 	go func() {
// 		tick := time.NewTicker(time.Millisecond * 100)
// 		defer tick.Stop()
// 		for {
// 			select {
// 			case e := <-p.monitor:
// 				zlog.Info(
// 					"recv pe monitor event",
// 					PluginLogFields.Watcher.With(
// 						zzlog.String("plugin", p.Name),
// 						zzlog.String("path", p.Path),
// 						zzlog.String("eventType", e.eventType.String()),
// 						zzlog.String("dir", "before"),
// 					)...)
// 				switch {
// 				case e.eventType&PluginBundled != 0:
// 					p.LoadBundle(
// 						func(p *Plugin, gov8 *plugin.GoV8Env, exports *v8go.Value) {
// 							pd := ErrorAssert(ModuleValueToExportsMap(gov8.Ctx, exports))
// 							p.PluginData = pd
// 							zlog.Info(
// 								"load bundle val",
// 								PluginLogFields.Plugin.With(
// 									zzlog.String("plugin", p.String()),
// 									zzlog.Reflect("dump", pd),
// 								)...)
// 						},
// 					)
// 				case e.eventType == PluginWindBlows:
// 					fallthrough
// 				case e.eventType&PluginModified != 0:
// 					fallthrough
// 				case e.eventType&PluginInit != 0:
// 					// zlog.Info("Monitor", PluginLogFields.Watcher.With(zzlog.String("eventType", e.eventType.String()), zzlog.String("dir", "before"))...)
// 					// fmt.Printf("monitor_plugin_init: %s %s %+v\n", p.Name, p.Path, e)
//
// 				}
// 			case <-tick.C:
// 			}
// 		}
// 	}()
// 	return nil
// }

func (pm *PluginsManager) StartPluginMonitor() error {
	pm.Scan("js/plugins")
	debouncedChan := debounce(
		50*time.Millisecond,
		200*time.Millisecond,
		pm.fsmonitor,
	)

	go func() {
		for {
			select {
      case i := <-pm.global:
      // zlog.Debug("received global pe", PluginLogFields.Manager.With(
      //   zzlog.Reflect("manager", pm),
      //   zzlog.Reflect("event", i),
      //   )...)
      // events.Global.Notify(*pm.NewObservableEvent(PluginEvents, i.name, &i, pm))
        pm.Notify(*pm.NewObservableEvent(PluginEvents, i.name, &i, pm))

      // zlog.Debug("sent observable notify", PluginLogFields.Manager.With(
      //   zzlog.Reflect("manager", pm),
      //   zzlog.Reflect("event", i),
      //   )...)
			case i := <-debouncedChan:
				event := i.(fsnotify.Event)
				p, m := pm.getPluginByPath(filepath.Dir(event.Name))
				if p == nil {
					pm.Scan("js/plugins")
				} else {
					e := NewPluginEvent(PluginModified, "fsmods", p)
					m <- *e
      zlog.Debug("sent plugin event", PluginLogFields.Manager.With(
        zzlog.Reflect("event", e),
        )...)
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

	_, _ = p.Bundle()

	// dump.P(result)
	pm.monitors = append(pm.monitors, monitor)

	p.LoadMeta()
	pm.watcher.Add(ppath)
	p.NewWorker(func(p *Plugin) {

	})
	// p.StartMonitor()
  // zlog.Info(fmt.Sprintf("plugin %s", p.Name), zzlog.Object("plugin", p))

	pm.Plugins = append(pm.Plugins, p)

	p, m := pm.getPluginByPath(ppath)
	m <- *NewPluginEvent(PluginModified, "modded", p)
	// m <- *NewPluginEvent(PluginModified, "modded", nil)
  pm.DispatchEvent(NewPluginEvent(PluginInit, p.Name, p))

	return nil
}

func (pm *PluginsManager) Init() error {
	//  return nil
	return nil
}

// func Handle
