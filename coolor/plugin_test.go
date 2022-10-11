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
 
