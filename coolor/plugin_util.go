package coolor

import (
	"encoding/json"
	"path/filepath"

	"github.com/gookit/goutil/dump"
	"rogchap.com/v8go"

	"github.com/digitallyserviced/coolors/coolor/plugin"
)

func GetKeyMap(
	mapType PluginKeyMapType,
	ctx *v8go.Context,
	obj *v8go.Value,
) (pd *PluginData, err error) {

	return nil, nil
}

func ModuleExportsFromDefaults(
	ctx *v8go.Context,
	obj *v8go.Value,
) (exports *v8go.Object, err error) {
	if !obj.IsObject() {
		return
	}
	exports = &v8go.Object{}

	if obj.IsObject() && eajs(obj.AsObject()).Has("default") {
		exports = eajs(eajs(eajs(obj.AsObject()).Get("default")).AsObject())
	} else {
		exports = eajs(obj.AsObject())
	}
	return
}
var workerIdx = 0

func nextWorkerIdx() uint {
	workerIdx += 1
	return uint(workerIdx)
}


func InjectBaseObjects(ctx *v8go.Context, iso *v8go.Isolate) {
	obTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginTypes {
		obTpl.Set(v.S, v.V, v8go.None)
	}
	peTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginEventTypes {
		peTpl.Set(v.S, v.V, v8go.None)
	}
	pkmTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginKeyMapTypes {
		pkmTpl.Set(v.S, v.V, v8go.None)
	}
	pdTpl := v8go.NewObjectTemplate(iso)
	for _, v := range PluginDetectionTypes {
		pdTpl.Set(v.S, v.V, v8go.None)
	}
	if err := ctx.Global().Set("pluginEventTypes", eajs(peTpl.NewInstance(ctx))); err != nil {
		panic(err)
	}
	if err := ctx.Global().Set("pluginDetectionTypes", eajs(pdTpl.NewInstance(ctx))); err != nil {
		panic(err)
	}
	if err := ctx.Global().Set("pluginKeyMapTypes", eajs(obTpl.NewInstance(ctx))); err != nil {
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

func ModuleValueToExportsMap(
	ctx *v8go.Context,
	obj *v8go.Value,
) (pd *PluginData, err error) {
	pd = &PluginData{}

	exports := eajs(ModuleExportsFromDefaults(ctx, obj))
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
}
func (p *Plugin) getPluginPath() string {
	return filepath.Join(p.Path, "index.js")
}
func (p *Plugin) getHandlersPath() string {
	return filepath.Join(p.Path, "handlers.js")
}
func (p *Plugin) getBundledHandlersPath() string {
	return filepath.Join(p.getBundleDir(), "handlers.js")
}
func (p *Plugin) getMetaPath() string {
	return filepath.Join(p.Path, "meta.js")
}
func (p *Plugin) getBundleDir() string {
	return filepath.Join(p.Path, ".bundled")
}
func (p *Plugin) getBundledPath() string {
	return filepath.Join(p.getBundleDir(), "index.js")
}
