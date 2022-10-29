package coolor

import (
	"container/ring"
	"sync"

	"rogchap.com/v8go"

	"github.com/digitallyserviced/coolors/coolor/plugin"
)

const (
	MAX_PLUGIN_WORKERS  = 16
	MAX_BUFFERED_EVENTS = 16
)

type PluginWorkers struct {
	Workers []*PluginWorker
}

type PluginWorker struct {
	idx uint
	*Plugin
	*PluginWorkerSandbox
	status        chan PluginWorkerStatus
	statusHistory *ring.Ring
	events        chan *PluginEvent
	done          chan struct{}
	once          *sync.Once
	onLoad        func(p *Plugin, psf *PluginSchemeFile)
}

type PluginMetaWorker struct {
	*PluginWorker
}

type PluginWorkerSandbox struct {
	gv8   *plugin.GoV8Env
	data  map[string]interface{}
	valid bool
}

type PluginFileHandler struct {
	Import *v8go.Function
	Export *v8go.Function
	TagMap *v8go.Object
	Misc   *v8go.Object
}

type PluginEventHandler interface {
	Handle(pe *PluginEvent, gv8 *plugin.GoV8Env, args ...interface{}) error
}

type PluginWorkerInitFunc func(p *Plugin, gv8 *plugin.GoV8Env, args ...interface{}) error
type PluginEventCallback func(pe *PluginEvent, gv8 *plugin.GoV8Env, args ...interface{}) error
type PluginWorkerResetFunc func(args ...interface{}) error

// type PluginEventHandlerFunc func(pe *PluginEvent, gv8 *plugin.GoV8Env, args ...interface{}) error

type PluginEvent struct {
	eventType PluginEventType
	plugin    *Plugin
	name      string
	refs      []interface{}
	callback  PluginEventCallback
	status    PluginEventStatus
}

// GetRef implements Referenced
func (pe *PluginEvent) GetRef() interface{} {
  return pe
}
