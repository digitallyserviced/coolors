package coolor_test

import (
	"testing"

	"github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/coolor"
)

func TestLoadMeta(t *testing.T) {
  pm := coolor.NewPluginsManager()
  pm.InitPlugin("../js/plugins/wezterm")
  dump.P(pm.Plugins)
  // dump.P(meta.As)
}
 

func (pw *PluginWorker) EventLoop() {
	tick := time.NewTicker(time.Millisecond * 500)
	defer tick.Stop()
	for {
		select {
		case <-pw.done:
		case e := <-pw.events:
			zlog.Info(
				"plugin event",
				PluginLogFields.Watcher.With(
					zzlog.Object("plugin", pw.Plugin),
					zzlog.String("event", e.eventType.String()),
				)...)
			switch {
			case e.eventType&PluginBundled != 0:
				// pw.Plugin.LoadBundle(
					// func(p *Plugin, gov8 *plugin.GoV8Env, exports *v8go.Value) {
						pd := ErrorAssert(ModuleValueToExportsMap(gov8.Ctx, exports))
						p.PluginData = pd
					// 	zlog.Info(
					// 		"load bundle val",
					// 		PluginLogFields.Plugin.With(
					// 			zzlog.String("plugin", p.String()),
					// 			zzlog.Reflect("dump", pd),
					// 		)...)
					// },
				// )
			case e.eventType == PluginWindBlows:
				fallthrough
			case e.eventType&PluginModified != 0:
					// p.manager.bundler <- e
				err := pw.Plugin.LoadMeta()
				if err != nil {
					panic(err)
				}
				fallthrough
			case e.eventType&PluginInit != 0:
				// zlog.Info("Monitor", PluginLogFields.Watcher.With(zzlog.String("eventType", e.eventType.String()), zzlog.String("dir", "before"))...)
				// fmt.Printf("monitor_plugin_init: %s %s %+v\n", p.Name, p.Path, e)

				// pw.Plugin.manager.bundler <- e
			}
		case <-tick.C:
			// tick plugin handlers if need be
		}
	}
}
