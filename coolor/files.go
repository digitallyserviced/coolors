package coolor

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/arrutil"

	// "github.com/gookit/goutil/dump"

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
	Name     string    `koanf:"name"`
	Palettes []string  `koanf:"palettes"`
	Version  uint64    `koanf:"version"`
}

type PaletteData struct {
	Name   string   `koanf:"name"`
	Names  []string `koanf:"names"`
	Colors []string `koanf:"colors"`
	Hash   uint64   `koanf:"hash"`
}

func (pd PaletteData) GetPalette() *CoolorColorsPalette {
	// hash := HashCssColors(pd.Colors)
	// if hash != pd.Hash {
	// 	dump.P(fmt.Sprintf("Hashes do not match... %d != %d", hash, pd.Hash))
	// }
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
	Palettes []PaletteData   `koanf:"palettes"`
	Metadata PaletteMetaData `koanf:"metadata"`
}

type PaletteFile struct {
	ref     *os.File
	path    string
	name    string
	version uint64
	tmp     bool
}

type MetaData interface {
  GetMeta() interface{}
}

type HistoryDataConfig struct {
	*PaletteFile
	*koanf.Koanf
	data *CoolorPaletteData
  Meta []MetaData
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

func TupleToEntry(item lo.Tuple2[string, string], _ int) lo.Entry[string, string] {
	var a  = item.A
	var b  = item.B
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
				return v.GetPalette()
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
	// if pdc.version != 0 && pdc.version <= pdc.GetFileVersion() {
	// 	panic(errorx.New("version too low"))
	// }
	pdc.UpdateVersion(pdc.version)
	// err := pdc.Koanf.Load(structs.Provider(pdc.data, "koanf"), nil)
	// if err != nil {
	// 	panic(err)
	// }
	//
	b, err := pdc.Marshal(toml.Parser())
	// dump.P(b)
	if err != nil {
		panic(err)
	}
	f, err := fsutil.QuickOpenFile(pdc.path)
	if err != nil {
		panic(err)
	}
	pdc.ref = f
	pdc.ref.Truncate(0)
	_, err = pdc.ref.Write(b)
	if err != nil {
		panic(err)
	}
	pdc.ref.Close()
	pdc.ref = nil
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
	pdc.Delete("")
	pdc.data.Metadata.Version = pdc.version
	err := pdc.Load(structs.Provider(*pdc.data, "koanf"), nil)
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
  pdc.Meta = make([]MetaData, 0)
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
	
	
}

func (pdc *HistoryDataConfig) NeedsSave() bool {
	return pdc.Status() >= 0 
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
  cp.UpdateHash()
  ccm := cp.GetMeta()
  fmt.Println(Generator().WithSeed(int64(cp.UpdateHash())).GenerateName(2),ccm.Current.Name, ccm.Named, ccm.String())
  ccm.Update(false)
  // fmt.Println(ccm)
}

  // if cp.Hash != hash {
  //
  // }

func (cp *CoolorColorsPalette) UpdateHash() uint64 {
  cp.Hash = cp.HashColors()
  return cp.Hash
}

func (cp *CoolorColorsPalette) HashColors() uint64 {
	// var hash uint64 = 0
  hashed := lo.Reduce[*CoolorColor, uint64](cp.Colors, func(h uint64, c *CoolorColor, i int) uint64 {
    return h + uint64(c.Color.Hex())
  }, 0)
	// for _, v := range cp.Colors {
	// 	hash += uint64(v.Color.Hex())
	// }
	return hashed
}

func (cp *CoolorColorsPalette) ToMap() map[string]string {
	outcols := make(map[string]string)
	for i, v := range cp.Colors {
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
// func (pdc *HistoryDataConfig) AddPalette(name string, p Palette) {
// 	cp := p.GetPalette()
//   cp.UpdateHash()
//   ccm := cp.GetMeta()
//   fmt.Println(ccm.Named, ccm.String())
//   fmt.Println(ccm)
	// if cp == nil {
	// 	panic(errorx.Errorf("Unable to save %d %s to %s", cp.GetItemCount(), name, pdc.path))
	// }
	// name = fmt.Sprintf("%s.%d", name, len(pdc.data.Metadata.Palettes))
	// flat := cp.ToMap()
	// colors := make([]string, 0)
	// names := make([]string, 0)
	// for x, v := range flat {
	// 	// k := fmt.Sprintf("%s", x)
	// 	names = append(names, x)
	// 	colors = append(colors, v)
	// 	// colors[k] = v
	// }
	// pdc.data.Palettes = append(pdc.data.Palettes, PaletteData{
	// 	Names:  names,
	// 	Name:   name,
	// 	Colors: colors,
	// 	Hash:   cp.Hash(),
	// })
	// pdc.data.Metadata.Palettes = append(pdc.data.Metadata.Palettes, name)
	// pdc.UpdateVersion(cp.Hash())
	// pdc.SetConfigData(nil)
	// if pdc.NeedsSave() {
	// 	pdc.Save()
	// }
// }
