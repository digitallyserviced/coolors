module github.com/digitallyserviced/coolors

go 1.18

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/charmbracelet/harmonica v0.2.0
	github.com/digitallyserviced/tview v0.0.0-00010101000000-000000000000
	github.com/dmarkham/enumer v1.5.5
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gdamore/tcell/v2 v2.5.2
	github.com/gookit/goutil v0.5.13
	github.com/jphsd/graphics2d v0.0.0-20220717174954-0a0ff2476d4b
	github.com/json-iterator/go v1.1.12
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mattn/go-runewidth v0.0.14
	github.com/mazznoer/colorgrad v0.9.0
	github.com/pgavlin/femto v0.0.0-20201224065653-0c9d20f9cac4
	github.com/samber/lo v1.26.0
	github.com/timshannon/bolthold v0.0.0-20210913165410-232392fc8a6a
	github.com/vmihailenco/msgpack/v5 v5.3.5
	github.com/zyedidia/micro v1.4.1
	go.etcd.io/bbolt v1.3.6
	go.uber.org/zap v1.23.0
	golang.org/x/term v0.0.0-20220919170432-7a66f970e087
	rogchap.com/v8go v0.7.0
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/atotto/clipboard v0.1.2 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/creack/pty v1.1.18
	github.com/evanw/esbuild v0.15.10
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gookit/color v1.5.2
	github.com/knadh/koanf v1.4.3
	github.com/mazznoer/csscolorparser v0.1.2 // indirect
	github.com/pascaldekloe/name v1.0.0 // indirect
	github.com/rakyll/autopprof v0.1.0
	github.com/rivo/uniseg v0.4.2
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sys v0.0.0-20220926163933-8cfa568d3c25
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/digitallyserviced/tview => /home/chris/Documents/coolors/pkg/tview

replace github.com/pgavlin/femto => /home/chris/Documents/coolors/pkg/femto

// replace github.com/rivo/tview@v0.0.0-20201204190810-5406288b8e4e => /home/chris/Documents/coolors/pkg/tview // indrect

// replace github.com/josa42/term-finder/tree => /home/chris/Documents/term-finder/tree

replace github.com/josa42/term-finder/tree => /home/chris/Documents/coolors/tree

replace rogchap.com/v8go => ../v8go

replace github.com/gdamore/tcell/v2 => /home/chris/Documents/coolors/pkg/tcell/v2

replace github.com/mattn/go-runewidth => github.com/mattn/go-runewidth v0.0.9
