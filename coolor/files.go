package coolor

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/dump"

	// "github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/samber/lo"

	// "github.com/gookit/goutil/maputil"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

type PaletteMetaData struct {
	Created  time.Time `koanf:"created"`
	Version  uint64    `koanf:"version"`
	Name     string    `koanf:"name"`
	Palettes []string  `koanf:"palettes"`
}

type PaletteData struct {
	Name   string   `koanf:"name"`
	Names  []string `koanf:"names"`
	Hash   uint64   `koanf:"hash"`
	Colors []string `koanf:"colors"`
}

func (pd PaletteData) GetPalette() *CoolorPalette {
	hash := HashCssColors(pd.Colors)
	if hash != pd.Hash {
		dump.P(fmt.Sprintf("Hashes do not match... %d != %d", hash, pd.Hash))
	}
	pairs := lo.Zip2(pd.Names, pd.Colors)
	entries := lo.Map(pairs, TupleToEntry)
	mapper := lo.FromEntries(entries)
	return NewCoolorPaletteFromMap(mapper).GetPalette()
}

type WezPaletteData struct {
	Foreground    string   `koanf:"foreground"`
	Background    string   `koanf:"background"`
	Cursor_bg     string   `koanf:"cursor_bg"`
	Cursor_border string   `koanf:"cursor_border"`
	Cursor_fg     string   `koanf:"cursor_fg"`
	Selection_bg  string   `koanf:"selection_bg"`
	Selection_fg  string   `koanf:"selection_fg"`
	Ansi          []string `koanf:"ansi"`
	Brights       []string `koanf:"brights"`
}

type CoolorPaletteData struct {
	Metadata PaletteMetaData `koanf:"metadata"`
	Palettes []PaletteData   `koanf:"palettes"`
}

type PaletteFile struct {
	tmp     bool
	version uint64
	path    string
	name    string
	ref     *os.File
}

type HistoryDataConfig struct {
	*PaletteFile
	*koanf.Koanf
	data *CoolorPaletteData
}

var (
	APPNAME string = "coolor"
	k              = koanf.New(".")
)

func (pdc *HistoryDataConfig) GetPalettes() []string {
	if pdc != nil && pdc.data != nil && pdc.data.Palettes != nil {
		return pdc.data.Metadata.Palettes
	}
	return make([]string, 0)
}

func TupleToEntry(item lo.Tuple2[string, string], i int) lo.Entry[string, string] {
	var a string = item.A
	var b string = item.B
	var e lo.Entry[string, string]
	e.Key = a
	e.Value = b
	return e
}

func HashCssColors(colors []string) uint64 {
	hash := lo.SumBy[string, uint64](colors, func(s string) uint64 {
		cc := tcell.GetColor(s)
		return uint64(cc.Hex())
	})
	return hash
}

func (pdc *HistoryDataConfig) LoadPalette(s string) Palette {
	if arrutil.Contains(pdc.data.Metadata.Palettes, s) {
		for _, v := range pdc.data.Palettes {
			if v.Name == s {
				v.GetPalette()
			}
		}
	}
	return nil
}

func LoadConfigFrom(path string) *HistoryDataConfig {
	if fsutil.IsDir(path) {
		return nil
	}
	pdc := NewPaletteHistoryData()
	err := pdc.LoadConfigFromFile(path, true)
	if err != nil {
		panic(err)
	}
	// dump.P(pdc.data)
	return pdc
}

func (pdc *HistoryDataConfig) FixFileVersion() {
	// pdc.BumpVersion()
	pdc.UpdateFileVersion(pdc.version)
	if pdc.NeedsSave() {
		pdc.Save()
	}
}

func (pdc *HistoryDataConfig) Save() {
	if pdc.version != 0 && pdc.version <= pdc.GetFileVersion() {
		panic(errorx.New("version too low"))
	}
	pdc.UpdateVersion(pdc.version)
	// err := pdc.Koanf.Load(structs.Provider(pdc.data, "koanf"), nil)
	// if err != nil {
	// 	panic(err)
	// }
	//
	b, err := pdc.Koanf.Marshal(toml.Parser())
	// dump.P(b)
	if err != nil {
		panic(err)
	}
	f, err := fsutil.QuickOpenFile(pdc.PaletteFile.path)
	if err != nil {
		panic(err)
	}
	pdc.PaletteFile.ref = f
	pdc.PaletteFile.ref.Truncate(0)
	_, err = pdc.PaletteFile.ref.Write(b)
	if err != nil {
		panic(err)
	}
	pdc.PaletteFile.ref.Close()
	pdc.PaletteFile.ref = nil
}

func (pdc *HistoryDataConfig) SetConfigData(k *koanf.Koanf) {
	if k == nil && pdc.Koanf == nil {
		pdc.NewTempConfigFile(pdc.name)
		if pdc.Koanf == nil {
			panic(errorx.New("No proper config setup"))
		}
	}
	if pdc.Koanf == nil {
		pdc.Koanf = k
	}
	pdc.Koanf.Delete("")
	pdc.data.Metadata.Version = pdc.version
	err := pdc.Koanf.Load(structs.Provider(*pdc.data, "koanf"), nil)
	if err != nil {
		// dump.P(err)
		panic(err)
	}
}

func (pdc *HistoryDataConfig) NewTempConfigFile(name string) *koanf.Koanf {
	k := koanf.New(".")
	path := GetTempFile(name)
	f, err := TempPalettesFile(name)
	if err != nil {
		panic(err)
	}

	pf := &PaletteFile{
		tmp:  true,
		path: path,
		name: name,
		ref:  f,
	}

	pdc.PaletteFile = pf
	pdc.Koanf = k
	pdc.Save()
	return k
}

func NewPaletteHistoryFile() *HistoryDataConfig {
	pdc := NewPaletteHistoryData()

	name := fmt.Sprintf("palette_%x", time.Now().Unix())
	pdc.NewTempConfigFile(name)
	return pdc
}

func NewPaletteHistoryData() *HistoryDataConfig {
	now := time.Now()
	pdc := &HistoryDataConfig{
		PaletteFile: &PaletteFile{
			tmp:  true,
			path: "",
			name: "",
			ref:  nil,
		},
		// Config: c,
		data: &CoolorPaletteData{
			Metadata: PaletteMetaData{
				Created:  now,
				Name:     "untitled",
				Palettes: []string{},
			},
			Palettes: make([]PaletteData, 0),
		},
	}
	return pdc
}

func (pdc *HistoryDataConfig) Dirty() bool {
	return pdc.version != pdc.GetFileVersion()
}

func (pdc *HistoryDataConfig) Status() int {
  return 1
	if pdc.Dirty() {
		kv := pdc.GetFileVersion()
		if pdc.version > kv {
			return 1
		} else if pdc.version < kv {
			return -1
		} else {
			return 0
		}
	}
	return 0
}

func (pdc *HistoryDataConfig) NeedsSave() bool {
	if pdc.Status() >= 0 {
		return true
	}
	return false
}

func (pdc *HistoryDataConfig) GetFileVersion() uint64 {
	version := k.Int("metadata.version")
	if version != 0 {
		return uint64(version)
	}
	return uint64(0)
}

func (pdc *HistoryDataConfig) UpdateFileVersion(i uint64) {
	// dump.P("pdc.version = %d and pdc.data.metadata.versionn = %d", pdc.version, pdc.data.Metadata.Version)
	pdc.data.Metadata.Version = i
}

func (pdc *HistoryDataConfig) UpdateVersion(i uint64) {
	// dump.P("pdc.version = %d and pdc.data.metadata.versionn = %d", pdc.version, pdc.data.Metadata.Version)
	pdc.version = i
}

func (pdc *HistoryDataConfig) BumpVersion() {
	// dump.P("pdc.version = %d and pdc.data.metadata.versionn = %d", pdc.version, pdc.data.Metadata.Version)
	pdc.UpdateVersion(uint64(time.Now().Unix()))
}

func (pdc *HistoryDataConfig) AddPalette(name string, p Palette) {
	cp := p.GetPalette()
	if cp == nil {
		panic(errorx.Errorf("Unable to save %d %s to %s", cp.GetItemCount(), name, pdc.PaletteFile.path))
	}
	name = fmt.Sprintf("%s.%d", name, len(pdc.data.Metadata.Palettes))
	flat := cp.ToMap()
	colors := make([]string, 0)
	names := make([]string, 0)
	for x, v := range flat {
		// k := fmt.Sprintf("%s", x)
		names = append(names, x)
		colors = append(colors, v)
		// colors[k] = v
	}
	pdc.data.Palettes = append(pdc.data.Palettes, PaletteData{
		Names:  names,
		Name:   name,
		Colors: colors,
		Hash:   cp.Hash(),
	})
	pdc.data.Metadata.Palettes = append(pdc.data.Metadata.Palettes, name)
	pdc.UpdateVersion(cp.Hash())
	pdc.SetConfigData(nil)
	if pdc.NeedsSave() {
		pdc.Save()
	}
}

func (cp *CoolorPalette) Hash() uint64 {
	var hash uint64 = 0
	for _, v := range cp.colors {
		hash += uint64(v.color.Hex())
	}
	return hash
}

func (cp *CoolorPalette) ToMap() map[string]string {
	outcols := make(map[string]string)
	for i, v := range cp.colors {
		k := fmt.Sprintf("color%d", i)
		outcols[k] = v.Html()
	}
	return outcols
}

func GetDataDirs() (string, string, string, error) {
	p, err := os.UserConfigDir()
	if err != nil {
		return "", "", "", err
	}
	configPath := strings.Join([]string{p, "coolor"}, "/")
	historyPath := path.Join(configPath, "history")
	palettePath := path.Join(configPath, "palettes")

	for _, path := range []string{historyPath, palettePath} {
		if !fsutil.PathExists(path) {
			err := fsutil.Mkdir(path, fs.ModePerm)
			if err != nil {
				fmt.Println(err)
				return "", "", "", err
			}
		}
	}
	return configPath, historyPath, palettePath, nil
}

func TempPalettesFile(name string) (*os.File, error) {
	path := GetTempFile(name)
	return OpenPalettesFile(path)
}

func OpenPalettesFile(path string) (*os.File, error) {
	f, err := fsutil.CreateFile(path, 0o664, 0o666)
	// f, err := fsutil.QuickOpenFile(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetTempFile(name string) string {
	_, historyPath, _, err := GetDataDirs()
	if err != nil {
		panic(err)
	}
	path := strings.Join([]string{historyPath, fmt.Sprintf("%s.toml", name)}, string(os.PathSeparator))
	return path
}

func (pdc *HistoryDataConfig) LoadConfigFromFile(path string, overwrite bool) error {
	pdc.Koanf = koanf.New(".")
	err := k.Load(file.Provider(path), toml.Parser())
	if err != nil {
		return errorx.Newf("error loading config: %s err: %v", path, err)
	}
	err = k.Unmarshal("", pdc.data)
	if err != nil {
		return errorx.Stacked(errorx.Newf("error unmarshalling config: %s err: %v", path, err))
	}
	// dump.P(pdc.data)
	// if pdc.GetFileVersion() == 0 {
	// 	pdc.FixFileVersion()
	// 	if pdc.GetFileVersion() == 0 {
	// 		return errorx.Newf("No version found in template file and could not fix: %s err: %v", path)
	// 	}
	// }

	// pdc.UpdateVersion(pdc.GetPalettes())

	return nil
}

// func (pdc *PaletteDataConfig) InitConfigData(k *koanf.Koanf) {
// }

func Colorizer(s string) string {
	for _, v := range colorRegs {
		CheckForReg(v, s)
	}
	return ""
}

func CheckForReg(reg string, c string) {
	if match := regexp.MustCompile(reg).FindAllStringSubmatch(c, -1); match != nil {
		colors := make([]string, 0)
		for _, c := range match {
			if len(c) == 2 {
				colors = append(colors, c[1])
			}
		}
	}
	// regexp.MustCompile(reg).FindAllSubmatch()
	if matchIdxs := regexp.MustCompile(reg).FindAllStringSubmatchIndex(c, -1); matchIdxs != nil {
		// dump.P(matchIdxs)
	}
}
