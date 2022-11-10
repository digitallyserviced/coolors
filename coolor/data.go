package coolor

import (
	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/gookit/goutil/dump"
	// "github.com/alecthomas/chroma/lexers/v"
	// "github.com/gookit/goutil/dump"
	msgpack "github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"

	"github.com/gdamore/tcell/v2"
	cl "github.com/ostafen/clover/v2"
	doc "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
	bh "github.com/timshannon/bolthold"

	"github.com/digitallyserviced/coolors/coolor/events"
	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/zzlog"
	"github.com/digitallyserviced/coolors/status"
)

type (
	CoolorColorOrigin   uint64
	CoolorPaletteOrigin uint64
)

const (
	// Origins - Color
	ColorOriginRandom CoolorColorOrigin = 1 << iota
	ColorOriginUser
	ColorOriginMixer
	ColorOriginShade
	ColorOriginEdit
	ColorOriginArgument
	ColorOriginFile
	ColorOriginClipboard

	// Origins - Palette
	PaletteOriginDefault CoolorPaletteOrigin = 1 << iota
	PaletteOriginRandomGenerated
	PaletteOriginUserSpecified
	PaletteOriginUserAdded
	PaletteOriginUserTagged

	RecentCoolorsMax = 80

	// Persistence - Status
	// RandomGenerated CoolorPaletteType = 1 << iota
)

type (
	coolorDataVal interface {
		CoolorColor | CoolorColorsPalette | PluginSchemeFile | Plugin
	}
	CoolorsCache struct {
		cache map[int32]*CoolorColor
	}
	Coolor struct {
		Color tcell.Color
	}
	// CoolorMeta struct {
	// 	*Coolor
	// 	CssName   string
	// 	AnsiName  string
	// 	XtermName string
	// 	UserNamed string
	// 	Hex       string
	// 	tags      []TagItem
	// 	Seent
	// 	ID uint64
	// 	Besty
	// }
	// CoolorsMeta []CoolorMeta
	// Coolors     struct {
	// 	Hash   uint64
	// 	Key    string
	// 	Colors []*Coolor
	// 	Saved  bool
	// }
	TagsKeys              []string
	CoolorPaletteTagsMeta struct {
		tagCount     uint
		TaggedColors map[string]*Coolor
	}

	CoolorColorsPaletteMeta struct {
		Created    int64             `clover:",omitempty"`
		Name      string                 `clover:",omitempty"`
		Author    string                 `clover:",omitempty"`
		ExtraData map[string]interface{} `clover:",omitempty"`
		Palette   *CoolorColorsPalette   `clover:",omitempty"`
	}

	// 	Current  *Coolors
	// 	Versions []*Coolors
	// 	Started  time.Time
	// 	Named    string `boltholdUnique:"Named"`
	// 	ID       uint64 `                       boltholdKey:"ID"`
	// 	// CoolorPaletteOrigin
	// 	Saved bool
	// }
	// Besty struct {
	// 	time.Time
	// 	Best bool
	// }
	// Seent struct {
	// 	time.Time
	// 	Used   uint64
	// 	Origin CoolorColorOrigin
	// }

	CoolorData struct {
		*cl.DB
		*bh.Store
		opts bh.Options
		*MetaService
	}

	MetaService struct {
		// *EventObserver
		*EventNotifier
		Cache   CoolorsCache
		Current *CoolorColorsPaletteMeta
		// PaletteMeta    []*CoolorColorsPaletteMeta
		// RecentColors   Coolors
		FavoriteColors CoolorsCache
	}

	DataCollection struct {
		CollectionKey string
		Indexes       []string
	}

	DataDoc[T coolorDataVal] struct {
		Value *T
		*DataCollection
		doc.Document
		Prefix string
	}

	DocReader[T coolorDataVal] interface {
		Load() *T
	}

	DocMaker[T coolorDataVal] interface {
		GetDoc() *DataDoc[T]
	}

)

var (
	PluginCollection = &DataCollection{
		CollectionKey: "plugins",
		Indexes:       []string{"Name"},
	}
	PaletteCollection = &DataCollection{
		CollectionKey: "palettes",
		Indexes:       []string{"Name"},
	}
	ColorsCollection = &DataCollection{
		CollectionKey: "colors",
		Indexes:       []string{"favorite"},
	}
	Store *CoolorData
	db    *cl.DB
)

type (
	ColorDataDoc           = DataDoc[CoolorColor]
	PaletteDataDoc         = DataDoc[CoolorColorsPalette]
	ColorSchemeFileDataDoc = DataDoc[PluginSchemeFile]
	PluginDataDoc          = DataDoc[Plugin]
)

func FromTcell(col tcell.Color) *Coolor {
	return MakeColorFromTcell(col).GetCC().Coolor()
}
func (ms *MetaService) Name() string {
	return "meta_service"
}
func (ms *MetaService) Service() {
	// GetStore().FavoriteColors = *GetStore().FavoriteColors.Load("MetaService_Favorites")
	// ccpms := GetStore().PaletteHistory(false)
	// fmt.Println(ccpms)
	// GetStore().ForEach(nil, func(r *){
	//
	// })
	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tick.C:
			// GetStore().MetaService.FavoriteColors.Save(false)
			// GetStore().MetaService.LoadFavorites()
			// ms.RecentColors = GetStore().MetaService.ColorHistory(-24 * time.Hour)
			// log.Println(msFavoriteColors)
		}
	}
}

func GetStore() *CoolorData {
	if Store == nil {
		Store = &CoolorData{
			DB:          InitDB(),
			Store:       &bh.Store{},
			opts:        bh.Options{},
			MetaService: NewMetadataService(),
		}
		GetStore().MetaService.LoadFavorites()
    pals := GetStore().MetaService.LoadPalettes()
    for _, v := range pals {
      fmt.Println(v.Palette)
      
    }

		// Store.Store = openbolt()
		// seedbolt(Store.Store)
		// startBoltStats()
	}
	return Store
}

//	func (cc *Coolors) Add(c *CoolorColor) {
//		_, ok := cc.Contains(c)
//		// log.Println(col, ok)
//
//		if ok < 0 {
//			cc.Colors = append(cc.Colors, c.Coolor())
//		}
//	}
//
//	func (cc *Coolors) Remove(c *CoolorColor) {
//		_, ok := cc.Contains(c)
//		// colors := make([]Coolor, 0)
//		if ok >= 0 {
//			colors := make([]*Coolor, 0)
//			for _, v := range cc.Colors {
//				if v.Color.Hex() == c.Color.Hex() {
//					continue
//				}
//				colors = append(colors, v)
//			}
//			cc.Colors = colors
//		}
//	}
//
//	func (cc *Coolors) Contains(c *CoolorColor) (*Coolor, int) {
//		for i, v := range cc.Colors {
//			if v.Color.Hex() == c.Coolor().Color.Hex() {
//				return v, i
//			}
//		}
//		return nil, -1
//	}
func (oc CoolorColors) GetPalette() *CoolorColorsPalette {
	colors := NewCoolorColorsPalette()
	for _, v := range oc {
		colors.AddCoolorColor(v)
	}
	return colors
}

func (oc CoolorsCache) GetPalette() *CoolorColorsPalette {
	colors := NewCoolorColorsPalette()
	for _, v := range oc.cache {
		colors.AddCoolorColor(v)
	}
	return colors
}

func (cc CoolorsCache) Load(c *CoolorColor) *CoolorColor {
	_, ok := cc.Contains(c)
	if !ok {
		cc.cache[c.Color.Hex()] = c
	}
	return c
}

func (cc CoolorsCache) Add(c *CoolorColor) *CoolorColor {
	if c == nil || c.Color == nil {
		return nil
	}
	_, ok := cc.Contains(c)
	if !ok {
		cc.cache[c.Color.Hex()] = c
	}
	return c
}

func (cc CoolorsCache) Remove(c *CoolorColor) {
	_, ok := cc.cache[c.Color.Hex()]
	if ok {
		delete(cc.cache, c.Color.Hex())
	}
}

func (cc CoolorsCache) Contains(c *CoolorColor) (*CoolorColor, bool) {
	if c == nil || c.Color == nil {
		return nil, false
	}
	col, ok := cc.cache[c.Color.Hex()]
	if ok {
		return col, ok
	}
	return nil, false
}

func NewMetadataService() *MetaService {
	// recents := make(CoolorsCache)
	// favs := make(CoolorsCache)
	ms := &MetaService{
		// EventObserver: NewEventObserver("metaservice"),
		EventNotifier: NewEventNotifier("metaservice"),
		Cache:         CoolorsCache{cache: make(map[int32]*CoolorColor)},
		// RecentColors: Coolors{
		// 	Key:    "MetaService_Seent",
		// 	Colors: make([]*Coolor, 0),
		// },
		FavoriteColors: CoolorsCache{cache: make(map[int32]*CoolorColor)},
	}

	events.Global.Register(AllEvents, ms)
	// var o Observer = ms
	return ms
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Store = GetStore()
}

// func (cid CustomID) MarshalBinary() ([]byte, error) {
// 	return cid[:], nil
// }
//
// func (cid *CustomID) UnarshalBinary(b []byte) error {
// 	if len(b) != len(*cid) {
// 		return errors.New("wrong length")
// 	}
// 	copy((*cid)[:], b)
// 	return nil
// }

func (dc DataCollection) Save(docs ...*doc.Document) int {
	saved := 0
	for _, v := range docs {
		err := db.Save(dc.CollectionKey, v)
		// dump.P(v)
		if err != nil {
			zlog.Error(
				fmt.Sprintf("Failed saving document: %+v", err),
				zzlog.Any("error", err),
			)
			continue
		}
		saved += 1
	}
	return saved
}

func (dc DataCollection) Query() *q.Query {
	dc._exists()
	return q.NewQuery(dc.CollectionKey)
}

func (dc DataCollection) _exists() bool {
	if db == nil {
		panic(fmt.Errorf("No DB %s", ""))
	}
	collectionExists, _ := db.HasCollection(dc.CollectionKey)
	if !collectionExists {
		// Create a collection named 'todos'
		db.CreateCollection(dc.CollectionKey)
		return false
	}
	idxs, err := db.ListIndexes(dc.CollectionKey)
	if err != nil {
		panic(err)
	}
	// n, err := db.Count(dc.Query())
	if err != nil {
		panic(err)
	}
	// if n > 0 {
	if len(idxs) != len(dc.Indexes) {
		if len(dc.Indexes) > 0 {
			for _, v := range dc.Indexes {
				_ = v
				// db.CreateIndex(dc.CollectionKey, v)
			}
		}
		// }

	}

	return true
}

func DocNotExists(err error) bool {
	if err != nil {
		if err != cl.ErrDocumentNotExist {
			return true
		}
		panic(err)
	}
	return false
}

func NewCoolorColorsPaletteMetaFromDoc(ddc *doc.Document) *CoolorColorsPaletteMeta {
  ccpm := &CoolorColorsPaletteMeta{
  	Created:   ddc.Get("Created").(int64),
  	Name:      ddc.Get("Name").(string),
  	// Author:    IfElseStr(ddc.Has("Author")&& ddc.Get("Author") != nil, ,""),
  	ExtraData: make(map[string]interface{}),
  	Palette:   NewCoolorColorsPalette(),
  }
  if author, ok := ddc.Get("Author").(string); ok {
    ccpm.Author = author
  }

  pal := ddc.Get("Palette.Colors")
  for _, v := range pal.([]interface{}) {
    col := v.(map[string]interface{})
    ccpm.Palette.AddCoolorColor(NewIntCoolorColor(int32(col["Color"].(uint64))))
  }
	// c := NewDefaultCoolorColor()
	// err = ddc.Unmarshal(cc)
	// c.SetColorInt(int32(ddc.Get("Color").(uint64)))
	// c.Favorite = ddc.Get("favorite").(bool)
	// if err != nil {
	// 	panic(err)
	// }
	return ccpm
}
func NewCoolorColorFromDoc(ddc *doc.Document) *CoolorColor {
	c := NewDefaultCoolorColor()
	// err = ddc.Unmarshal(cc)
	c.SetColorInt(int32(ddc.Get("Color").(uint64)))
	c.Favorite = ddc.Get("favorite").(bool)
	// if err != nil {
	// 	panic(err)
	// }
	return c
}

func (cc *CoolorColor) Load() bool {
	var ddc *doc.Document
	ddc, err := db.FindFirst(
		ColorsCollection.Query().Where(q.Field("Color").Eq(uint64(*cc.Color))),
	)
	if err != nil {
		if !DocNotExists(err) {
			panic(err)
		}
	}
	cc = NewCoolorColorFromDoc(ddc)
	return true
}

func (cc *CoolorColor) Save() error {
	ColorsCollection._exists()
	nddc := doc.NewDocumentOf(cc)

	var ddc *doc.Document
	ddc, err := db.FindFirst(
		ColorsCollection.Query().Where(q.Field("Color").Eq(uint64(*cc.Color))),
	)
	if err != nil && !DocNotExists(err) {
		return err
	}
	if ddc != nil {
		nddc.Set(doc.ObjectIdField, ddc.ObjectId())
	}
	n := ColorsCollection.Save(nddc)
	if n > 0 {
		return nil
	}
	// dd := &DataDoc[CoolorColor]{
	//   Value: cc,
	//   Document: *ddc,
	//   DataCollection: ColorsCollection,
	// }
	// return dd
	return fmt.Errorf("Save didnt work %d", n)
}

func (cc *CoolorColor) GetDoc() *ColorDataDoc {
	dd := &DataDoc[CoolorColor]{
		Value:          cc,
		Document:       *doc.NewDocumentOf(cc),
		DataCollection: ColorsCollection,
	}
	return dd
}

func NewDataDoc[T coolorDataVal](
	name string,
	col DataCollection,
	v *T,
) *DataDoc[T] {
	return &DataDoc[T]{
		Value: v,
	}
}

// type ColorDoc struct {
// }

func InitDB() *cl.DB {
	ddb, err := cl.Open("clover-db")
	if err != nil {
		ddb, err = cl.Open("clover-dbs")
		if err != nil {
			panic(err)
		}
	}
	go handleSignals()
	db = ddb
	return ddb
	// defer db.Close()

	// q.CriteriaVisitor

	// cp := DefaultCoolorPalette()
	// Check if collection already exists
	// collectionExists, _ := db.HasCollection("todos")
	// if !collectionExists {
	// 	// Create a collection named 'todos'
	// 	db.CreateCollection("todos")
	// }
	//
	// cols := GetStore().MetaService.LoadFavorites()
	// dump.P(cols)

	// d := doc.NewDocumentOf(cp.CoolorColorsPalette)
	// err := db.Insert("todos", d)
	// if err != nil {
	//   panic(err)
	// }
	// cps, err := db.FindAll(q.NewQuery("todos"))
	// if err != nil {
	//   panic(err)
	// }
	// dump.P(cps)

	// Delete collection
	// err := db.DropCollection("todos")
	// if err != nil {
	// 	fmt.Printf("Failed to delete collection: %v", err)
	// }
	// fmt.Println("collection deleted")

	// base, _ := Hex("#1d344f")
	// done := make(chan struct{})
	// // cbp.colors = make(CoolorColors, 0)
	// defer close(done)
	// colors := RandomShadesStream(base, 0.5)
	// colors.Status.SetProgressHandler(NewProgressHandler(func(u uint32) {
	//    fmt.Println(fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()))
	// 	// status.NewStatusUpdate(
	// 	// 	"action_str",
	// 	// 	fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()),
	// 	// )
	// }, func(i uint32) {
	// 	// status.NewStatusUpdate(
	// 	// 	"action_str",
	// 	// 	fmt.Sprintf("Found Shades (%d / %d)", colors.Status.GetValid(), i),
	// 	// )
	// }))
	// colors.Run(done)
	// for _, v := range TakeNColors(done, colors.OutColors, 64) {
	//    fmt.Println(v.GetCC().TerminalPreview())
	//    fmt.Println(v.Tcell().TrueColor().Hex())
	// 	// SeentColor("stream_random_shade", newcc, newcc.pallette)
	// }
}

// func NewCoolorMeta(c *CoolorColor) CoolorMeta {
// 	ccm := &CoolorMeta{
// 		Coolor: &Coolor{
// 			Color: *c.Color,
// 		},
// 		AnsiName:  "",
// 		XtermName: "",
// 		UserNamed: "",
// 		Hex:       fmt.Sprintf("#%06X", c.Color.Hex()),
// 		CssName:   GetColorName(*c.Color),
// 		// Tags:      []TagItem{},
// 		Seent: Seent{
// 			Time:   time.Time{},
// 			Used:   1,
// 			Origin: ColorOriginRandom,
// 		},
// 		Besty: Besty{
// 			Time: time.Now(),
// 			Best: false,
// 		},
// 	}
// 	ccm.UpdateSeent(time.Now())
// 	ccm.Update(false)
//
// 	return *ccm
// }

func GenColors() {
	// colors := GenerateRandomCoolors(200)
	// a := &Coolor{
	// 	Color: tcell.ColorBlack,
	// }
	// ccm := NewCoolorMeta(a.Escalate())
	// err := Store.Insert(bh.NextSequence(), &ccm)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, c := range colors {
	// 	ccm := NewCoolorMeta(c.Escalate())
	// 	err := Store.Insert(bh.NextSequence(), &ccm)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// fmt.Println(c.Escalate().TerminalPreview())
	// }
}

func TrimSeentCoolors(n int) {
	return
	Store.Store.Bolt().Update(func(tx *bbolt.Tx) error {
		var b, c *bbolt.Bucket
		c = tx.Bucket([]byte("Coolors"))
		b = c.Bucket([]byte("Seent"))
		count := b.Stats().KeyN
		if count > n {
			rem := count - n
			log.Printf("Seent colors - trimming %d of %d", rem, count)
			c := b.Cursor()
			for k, _ := c.First(); k != nil && rem > 0; k, _ = c.Next() {
				c.Delete()
				rem--
			}
		}
		if true {
			// log.Printf("Checking dupe seents")
			uniq := make(map[int32]int)
			uc := 0
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				// fmt.Println(k,v)
				var ve int32
				e := dec(v, &ve)
				if checkErrX(e) {
					// return e
				}
				// fmt.Printf("%s", NewIntCoolorColor(int32(ve)).TerminalPreview())
				_, ok := uniq[ve]
				if ok {
					c.Delete()
					uniq[ve] += 1
					uc++
					fmt.Printf(
						"Dupe %s x%d of %d dupes",
						NewIntCoolorColor(int32(ve)).TerminalPreview(),
						uniq[ve],
						uc,
					)
				} else {
					uniq[ve] = 1
				}
			}
		}
		return nil
	})
}

// func (ms *MetaService) GetPaletteMeta(
// 	ccp *CoolorColorsPalette,
// ) *CoolorColorsPaletteMeta {
// 	if ms.Current == nil {
// 		ccm := NewCoolorColorsPaletteMeta("", ccp)
// 		ms.Current = &ccm
// 	} else {
// 		ms.Current.Current = ccp.Coolors()
// 	}
//
// 	return ms.Current
// }

func (ms *MetaService) ToggleFavorite(cc *CoolorColor) {
	cc.Favorite = !cc.Favorite
	err := cc.Save()
	if err != nil {
		panic(err)
	}
	_, ok := ms.FavoriteColors.Contains(cc)
	if ok && !cc.Favorite {
		ms.FavoriteColors.Remove(cc)
	} else if cc.Favorite {
		ms.FavoriteColors.Add(cc)
	}
	// ms.FavoriteColors.Save(false)
}

//	Store.Store.Bolt().View(func(tx *bbolt.Tx) error {
//		var b, c *bbolt.Bucket
//		c = tx.Bucket(bh.Ro)
//		b = c.Bucket([]byte("Seent"))
//		err := b.ForEach(func(k, v []byte) error {
//			var ve uint64
//			// var st uint64
//			// e := dec(k, &st)
//			// checkErr(e)
//			e := dec(v, &ve)
//			checkErr(e)
//			// colors = append(colors, NewIntCoolorColor(int32(ve)))
//			return nil
//		})
//		checkErr(err)
//		return nil
//	})

// func (oc *Coolors) Load(key string) *Coolors {
// 	var c = Coolors{
// 		Key:    key,
// 		Colors: make([]*Coolor, 100),
// 		Saved:  false,
// 	}
// 	// var c Coolors
// 	// var ccs []Coolors
// 	// err := Store.FindOne(&c, bh.Where(bh.Key).Eq(key))
// 	// fmt.Println("FUCK", c)
// 	// if err != nil {
// 	// 	if err == bh.ErrNotFound {
// 	// oc.Saved = true
// 	// oc.Save(false)
// 	// fmt.Print("shit", err, c)
// 	// 	}
// 	// 		checkErrX(err)
// 	// }
// 	c.Saved = true
// 	// fmt.Println("foundshit", c)
// 	return &c
// }

// func (ms *MetaService) Load(retry bool) {
// 	// mms := NewMetadataService()
// 	var c = Coolors{
// 		Key:    "MetaService_Favorites",
// 		Colors: make([]Coolor, 0),
// 	}
// 	//  var c []Coolors
// 	err := Store.FindOne(&c, bh.Where("Coolors.Key").Eq("MetaService_Favorites"))
// 	if err != nil {
// 		if err == bh.ErrNotFound && !retry {
// 			ms.Save(true)
// 			ms.Load(true)
// 			return
// 		}
// 		panic(err)
// 	}
// 	fmt.Println(c)
// 	// ms.FavoriteColors = mms.FavoriteColors
// 	// ms.RecentColors = mms.RecentColors
// 	ms.FavoriteColors = c
// }
//
// func (cs *Coolors) Save(insert bool) {
// 	// var err error
// 	if len(cs.Colors) > 0 {
// 		// err = Store.Upsert("MetaService_Favorites", cs)
// 		// checkErr(err)
//
// 	}
// }

func (ms *MetaService) Save(insert bool) {
	// var err error
	// if insert {
	//    err := Store.Insert("MetaService_Favorites", &ms.FavoriteColors)
	//    checkErr(err)
	// } else {
	//    err := Store.Update("MetaService_Favorites", &ms.FavoriteColors)
	//    checkErr(err)
	// }
	// log.Println(err)
}

// func (ms *MetaService) PaletteHistory(saved bool) []CoolorColorsPaletteMeta {
// 	var ccpms []CoolorColorsPaletteMeta
// 	//  savedQ := bh.Where("Started").Le(time.Now())
// 	// err := Store.Find(&ccpms, savedQ.SortBy("Started").Reverse())
// 	//  if err != nil && err != bh.ErrNotFound {
// 	//    panic(err)
// 	//  }
// 	return ccpms
// }

func (ms *MetaService) LoadPalettes() []*CoolorColorsPaletteMeta {
  err := db.Delete(PaletteCollection.Query().Where(q.Field("Name").NotExists()))
	if err != nil {
		panic(err)
	}
	pals, err := db.FindAll(
		PaletteCollection.Query(),
	)
	if err != nil {
		panic(err)
	}

  nccpm := make([]*CoolorColorsPaletteMeta, 0)
  for _, p := range pals {
    nccpm = append(nccpm, NewCoolorColorsPaletteMetaFromDoc(p))
 //    fmt.Println(p.AsMap())
 //    err = p.Unmarshal(nccpm)
	// if err != nil {
 //      fmt.Println(err)
	// 	// panic(err)
	// }

    // dump.P(nccpm)

	}
	// palss, err := db.FindAll(
	// 	PaletteCollection.Query().Where(q.Field("Name").Exists()),
	// )
	// if err != nil && !DocNotExists(err) {
	//
	// }
	// for _, v := range palss {
	// 	fmt.Println(v.Get("Name"))
	// }

	return nccpm 

}
func (ms *MetaService) LoadFavorites() *CoolorColors {
	colors := make(CoolorColors, 0)

	// for _, v := range cp.CoolorColorsPalette.Colors {
	//   d := v.GetDoc()
	//   db.InsertOne(ColorsCollection.CollectionKey, &d.Document)
	// }
	cols, err := db.FindAll(
		ColorsCollection.Query().Where(q.Field("favorite").IsTrue()),
	)
	if err != nil {
		panic(err)
	}

	for _, v := range cols {
		// var c CoolorColor
		//   // var u json.Unmarshaler = c
		// err := v.Unmarshal(&c)
		// if err != nil {
		// 	panic(err)
		// }
		c := NewCoolorColorFromDoc(v)
		ms.FavoriteColors.Add(c)
		// fmt.Println(c.TerminalPreview())
		colors = append(colors, c)
	}

	return &colors
}

func (ms *MetaService) ColorHistory(t time.Duration) *CoolorColors {
	colors := make(CoolorColors, 0)
	// Store.Store.Bolt().View(func(tx *bbolt.Tx) error {
	// 	var b, c *bbolt.Bucket
	// 	c = tx.Bucket([]byte("Coolors"))
	// 	b = c.Bucket([]byte("Seent"))
	// 	err := b.ForEach(func(k, v []byte) error {
	// 		var ve uint64
	// 		// var st uint64
	// 		// e := dec(k, &st)
	// 		// checkErr(e)
	// 		e := dec(v, &ve)
	// 		checkErr(e)
	// 		colors = append(colors, NewIntCoolorColor(int32(ve)))
	// 		return nil
	// 	})
	// 	checkErr(err)
	// 	return nil
	// })
	return &colors
}

func (ms *MetaService) HandleEvent(o ObservableEvent) bool {
	// if ObservableEventType(
	// 	ColorSeentEvent | ColorEvent | SelectedEvent,
	// ).Is(o.Type) {
	// 	// col, ok := o.Ref.(*CoolorColor)
	// 	// if !ok {
	// 	// 	return true
	// 	// }
	// 	// cm, ok := ms.Cache.Contains(col)
	// 	// if !ok {
	// 	// 	cmm := ms.Cache.Add(col)
	// 	// 	cm = &cmm
	// 	// }
	// 	// cm.Update(false)
	// }
	if ObservableEventType(PaletteCreatedEvent).Is(o.Type) {
		// fmt.Printf("DICK %+v %T", o.Ref, o.Ref)
		cp := ms.Current.ToPalette(o.Ref)
		// cp, ok := o.Ref.(*CoolorPaletteMainView)
		// if !ok {
		// 	return false
		// }
		if cp == nil {
			return false
		}
		ms.UpdateCurrent(cp)

		if ms.Current.Save() {
			events.Global.Notify(
				*events.Global.NewObservableEvent(events.PaletteSavedEvent, "saved_palette", ms.Current.Palette, ms.Current),
			)
		}
	}
	if ObservableEventType(
		PaletteMetaUpdatedEvent,
	).Is(o.Type) {
    status.NewStatusUpdate("title", ms.Current.Name)
	}
	// fmt.Printf("*** Data Observed %s %s received: %T  %T\n", o.Note,o.Type.String(), o.Ref, o.Src)
	// Store.MetaService.ColorHistory(-24 * time.Hour)
	return true
}

func (ms *MetaService) UpdateCurrent(cp *CoolorColorsPalette) {
	if !ms.Current.IsCurrent(cp) {
		ms.Current = ms.Current.ToMeta(cp)
			events.Global.Notify(
				*events.Global.NewObservableEvent(events.PaletteMetaUpdatedEvent, "current_palette", ms.Current.Palette, ms.Current),
			)
	}
}
// func (ms *MetaService) UpdateCurrent() bool {
// }

func (ccpm *CoolorColorsPaletteMeta) ToPalette(
	cp interface{},
) *CoolorColorsPalette {
	switch cpp := cp.(type) {
	case *CoolorColorsPaletteMeta:
		return cpp.Palette
	case *CoolorColorsPalette:
		return cpp
	case *CoolorPaletteMainView:
		return cpp.CoolorColorsPalette
	case Palette:
		return cpp.GetPalette()
	}
	return nil
}

func (ccpm *CoolorColorsPaletteMeta) ToMeta(
	cp interface{},
) (nccpm *CoolorColorsPaletteMeta) {
	// var nccpm *CoolorColorsPaletteMeta
	switch cpp := cp.(type) {
	case *CoolorColorsPaletteMeta:
		nccpm = cpp
	case *CoolorColorsPalette:
		nccpm = NewCoolorColorsPaletteMeta(cpp)
	case Palette:
		nccpm = NewCoolorColorsPaletteMeta(cpp.GetPalette())
	}
	return
}

func (ccpm *CoolorColorsPaletteMeta) IsCurrent(cpi interface{}) bool {
	cp := ccpm.ToPalette(cpi)
	if ccpm == nil || cp == nil {
		return false
	}
	if cp == ccpm.Palette {
		return true
	}
	if cp.Name() == ccpm.Name {
		return true
	}
	if cp.HashColors() == ccpm.Palette.HashColors() {
		return true
	}
	return false
}

func (ccpm *CoolorColorsPaletteMeta) Save() bool {
	fmt.Println(ccpm)
	if ccpm == nil {
		return false
	}
	dccpm := doc.NewDocumentOf(ccpm)
	pal, err := db.FindFirst(
		PaletteCollection.Query().Where(q.Field("Name").Eq(ccpm.Name)),
	)
	if err != nil {
		if !DocNotExists(err) {
			panic(err)
		}
	}
	if pal != nil {
		dccpm.Set(doc.ObjectIdField, pal.ObjectId())
	}
	n := PaletteCollection.Save(dccpm)
	return n > 0
}

var _ msgpack.CustomEncoder = (*CoolorColor)(nil)
var _ msgpack.CustomDecoder = (*CoolorColor)(nil)

func (s *CoolorColor) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(s.Color)
}

func (s *CoolorColor) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&s.Color)
}

func (cc *Coolor) Escalate() *CoolorColor {
	return NewIntCoolorColor(int32(cc.Color))
}

func (cc *CoolorData) Close() {
	db.Close()
}

func (cc *CoolorData) FindNamedPalette(name string) *CoolorColorsPaletteMeta {
	var ccpm CoolorColorsPaletteMeta
	ddc, err := GetStore().DB.FindFirst(
		PaletteCollection.Query().Where(q.Field("Name").Like(fmt.Sprintf("%s", name))),
	)
	if err != nil && DocNotExists(err) {
		return nil
	}
	if ddc != nil {
		err = ddc.Unmarshal(&ccpm)
		if err != nil {
			panic(err)
		}
	}
	// fmt.Println(ddc, ccpm)
	return &ccpm
	// err := GetStore().FindOne(&ccm, bh.Where("Named").Eq(name))
	// if err != nil {
	// 	if err == bh.ErrNotFound {
	// 		// fmt.Println("not found", err, name)
	// 		return nil
	// 	}
	// }
}

// func (cc *CoolorMeta) UpdateSeent(t time.Time) {
// 	return
// 	if cc.Seent.Time.IsZero() {
// 		cc.Seent.Time = time.Now()
// 	}
// 	seent := []byte(cc.Seent.Time.Format(time.RFC3339))
// 	Store.Bolt().Update(func(tx *bbolt.Tx) error {
// 		if t.IsZero() {
// 			t = time.Now()
// 		}
// 		var b, c *bbolt.Bucket
// 		c = tx.Bucket([]byte("Coolors"))
// 		b = c.Bucket([]byte("Seent"))
// 		v := b.Get(seent)
// 		if v != nil {
// 			var hex uint64
// 			e := dec(v, &hex)
// 			if e != nil {
// 				checkErr(e)
// 			}
// 			// fmt.Println(NewIntCoolorColor(int32(hex)).TerminalPreview())
// 			b.Delete(seent)
// 		}
// 		// t.UnixMicro()
// 		// newseent := t.Format(time.RFC3339)
// 		newseent := ErrorAssert[[]byte](enc(t.UnixMicro()))
// 		be, _ := enc(cc.Color.Hex())
// 		err := b.Put(newseent, be)
// 		checkErr(err)
// 		return nil
// 	})
//
// 	cc.Seent.Time = t
// 	cc.Seent.Used += 1
// }

// func (cc *CoolorMeta) Update(clean bool) {
// 	// k := uint64(cc.Color.TrueColor())
// 	// err := Store.Upsert(k, &cc)
// 	// checkErr(err)
// }

func (cc CoolorColors) Contains(c *CoolorColor) bool {
	for _, v := range cc {
		if v.Html() == c.Html() {
			return true
		}
	}
	return false
}

func (cc *CoolorColorsPalette) Update(clean bool) {
	if cc == nil || cc.Colors == nil {
		return
	}

	if cc.Hash == 0 && len(cc.Colors) > 0 {
		cc.UpdateHash()
	}

	// err := Store.Upsert(cc.Hash, &cc)
	// if err != nil {
	// 	panic(err)
	// }
}

func (ccs CoolorPaletteTagsMeta) String() string {
	str := ""
	for k, col := range ccs.TaggedColors {
		str = fmt.Sprintf(
			"%s %s",
			str,
			fmt.Sprintf("%s %s", k, col.Escalate().TerminalPreview()),
		)

	}
	return str
}

func (cc CoolorColor) FavString() string {
	str := fmt.Sprintf(
		" %s %s ",
		IfElseStr(cc.Favorite, "  ", "  "),
		cc.TVPreview(),
		// cc.Seent.Used,
	)
	return str
}

func handleSignals() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
			case syscall.SIGINT:
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
			// GetStore().Close()
		}
	}()

	code := <-exit_chan
	fmt.Printf("Closing DB... %d", code)
	db.Close()
	AppModel.app.Stop()
	// os.Exit(code)
}

// func (cc *CoolorColor) IsFavorite() bool {
//   return false
// }

// func (cc *CoolorColor) GetMeta() *CoolorMeta {
// 	// GetStore().MetaService.RecentColors.Add(cc)
// 	var cm CoolorMeta
// 	// k := uint64(cc.Color.TrueColor())
// 	// err := Store.FindOne(&cm, bh.Where(bh.Key).Eq(k))
// 	// if err != nil {
// 	// 	if err != bh.ErrNotFound {
// 	// 		checkErr(err)
// 	// 	}
// 	// 	cm = NewCoolorMeta(cc)
// 	// }
// 	// if _, ok := GetStore().FavoriteColors.Contains(cc); ok >= 0 {
// 	// 	cm.Best = true
// 	// } else {
// 	// 	cm.Best = false
// 	// }
// 	// cm.Update(false)
// 	return &cm
// }

// func (cc *CoolorColorsPalette) GetMeta() *CoolorColorsPaletteMeta {
// 	// cc.UpdateHash()
// 	// var pals []CoolorColorsPaletteMeta
// 	// ccm := NewCoolorColorsPaletteMeta(cc)
// 	var ccm CoolorColorsPaletteMeta
// 	// ccm.Current = cc
// 	// current := bh.Where("Current.Hash").Eq(cc.Hash)
// 	// version := bh.Where("Versions").Contains(cc.Hash)
// 	// err := Store.FindOne(&ccm, current)
// 	// if err != nil {
// 	// 	// fmt.Println("not found", err)
// 	//    // doCallers()
// 	// 	if err == bh.ErrNotFound {
// 	// 		var ccms []CoolorColorsPaletteMeta
// 	// 		err := Store.Find(&ccms, current)
// 	// 		if err != nil {
// 	// 			// fmt.Println(err)
// 	// 			panic(err)
// 	// 		}
// 	// 		for _, p := range ccms {
// 	//        // fmt.Println("pals:", p)
// 	// 			if p.Current.HashColors() == ccm.Current.HashColors() {
// 	// 				return &p
// 	// 			}
// 	// 		}
// 	// 		cc.UpdateHash()
// 	// 		ccm = NewCoolorColorsPaletteMeta(
// 	// 			Generator().WithSeed(int64(cc.Hash)).GenerateName(2),
// 	// 			cc,
// 	// 		)
// 	// 	} else {
// 	// 		panic(err)
// 	// 	}
// 	// }
// 	// fmt.Println(ccm)
// 	return &ccm
// }

// func (ccs CoolorsMeta) String() string {
// 	str := ""
// 	for _, _ = range ccs {
// 		// str = fmt.Sprintf("%s %s", str, v.String())
// 	}
// 	return str
// 	// return fmt.Sprintf("%s %s", )
// }

// func seedbolt(store *bh.Store) {
// err := store.Bolt().Update(func(tx *bbolt.Tx) error {
// 	if tx.Cursor().Bucket().Stats().KeyN == 0 {
// 		pals := errAss[*bbolt.Bucket](
// 			tx.CreateBucketIfNotExists([]byte("Palettes")),
// 		)
// 		colors := errAss[*bbolt.Bucket](
// 			tx.CreateBucketIfNotExists([]byte("Coolors")),
// 		)
// 		seent := errAss[*bbolt.Bucket](
// 			colors.CreateBucketIfNotExists([]byte("Seent")),
// 		)
// 		favs := errAss[*bbolt.Bucket](
// 			colors.CreateBucketIfNotExists([]byte("Favorites")),
// 		)
// 		recents := errAss[*bbolt.Bucket](
// 			pals.CreateBucketIfNotExists([]byte("Recents")),
// 		)
// 		anon := errAss[*bbolt.Bucket](
// 			pals.CreateBucketIfNotExists([]byte("Anonymous")),
// 		)
// 		user := errAss[*bbolt.Bucket](
// 			pals.CreateBucketIfNotExists([]byte("User")),
// 		)
// 		_, _, _, _, _ = seent, anon, user, favs, recents
// 	}
// 	return nil
// })
// checkErr(err)
// }

// func init(){
// handleSignals()
// }

// func openbolt() *bh.Store {
// 	store, err := bh.Open("testfile", 0o666, &bh.Options{
// 		// Encoder: bh.EncodeFunc(enc),
// 		// Decoder: bh.DecodeFunc(dec),
// 		Options: &bbolt.Options{
// 			FreelistType: bbolt.FreelistMapType,
// 		},
// 	})
// 	checkErr(err)
// 	// go Store.MetaService.Service()
// 	// go handleSignals()
//
// 	return store
// }

//
// func (i *Item) MarshalMsgpack() ([]byte, error) {
// 	v := item{
// 		Age: i.age,
// 		Str: i.str,
// 	}
//
// 	return msgpack.Marshal(v)
// }
//
// func (i *Item) UnmarshalMsgpack(b []byte) error {
// 	var result item
// 	err := msgpack.Unmarshal(b, &result)
//
// 	i.age = result.Age
// 	i.str = result.Str
//
// 	return err
// }

// GenColors()
// var m CoolorsMeta
// q := &bh.Query{}
// err := Store.Find(&m, q)
// fmt.Printf("shit %s %v", m, err)
// for _, v := range m {
//   fmt.Println(v.Coolor.Escalate().TerminalPreview())
// fmt.Println(v.ID, v.Coolor.Color, v.Hex, time.Now().Sub(v.Seent.Time).Seconds(), v.AnsiName)
// }

//
// func (i *Coolors) Type() string { return "Coolors" }
// func (i *Coolors) Indexes() map[string]bh.Index {
// 	return map[string]bh.Index{
// 		"Key": {
// 			IndexFunc: func(n string, value interface{}) ([]byte, error) {
// 				// If the upsert wants to delete an existing value first,
// 				// value could be a **Item instead of *Item
// 				// panic: interface conversion: interface {} is **Item, not *Item
// 				v := value.(*Coolors).Key
// 				return bh.DefaultEncode(v)
// 			},
// 			Unique: true,
// 		},
// 	}
// }
// func (i *Coolors) SliceIndexes() map[string]bh.SliceIndex {
// 	return map[string]bh.SliceIndex{
// 		// "Colors": func(name string, value interface{}) ([][]byte, error) {
// 		// 	cols, ok := value.(*Coolors)
// 		// 	keys := make([][]byte, len(cols.Colors))
// 		// 	if ok {
// 		// 		for i, v := range cols.Colors {
// 		// 			keys[i] = errAss[[]byte](enc(v.Color.Hex()))
// 		// 		}
// 		// 		return keys, nil
// 		// 	}
// 		// 	return keys, fmt.Errorf(
// 		// 		"Error casting to proper slice %v type %v",
// 		// 		cols,
// 		// 		ok,
// 		// 	)
// 		// },
// 	}
// }

// func (i *CoolorColorsPaletteMeta) Type() string { return "CoolorColorsPaletteMeta" }
// func (i *CoolorColorsPaletteMeta) Indexes() map[string]bh.Index {
// 	return map[string]bh.Index{
// 		"CoolorColorsPalette": {
// 			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
// 				// If the upsert wants to delete an existing value first,
// 				// value could be a **Item instead of *Item
// 				// panic: interface conversion: interface {} is **Item, not *Item
// 				v := fmt.Sprintf("%x", value.(*CoolorColorsPalette).Hash)
// 				return []byte(v), nil
// 			},
// 			Unique: false,
// 		},
// 	}
// }
//
// func (i *CoolorColorsPaletteMeta) SliceIndexes() map[string]bh.SliceIndex {
// 	return map[string]bh.SliceIndex{}
// }
// func (s *Coolors) EncodeMsgpack(enc *msgpack.Encoder) error {
// 	return enc.EncodeMulti(s.Key, s.Colors)
// }
//
// func (s *Coolors) DecodeMsgpack(dec *msgpack.Decoder) error {
// 	// var err error
// 	// fmt.Println(dec)
// 	// buf := make([]byte, 1024)
// 	// dec.ReadFull(buf)
// 	// // var k string
// 	var b []interface{}
// 	// // b, err := dec.DecodeBytes()
// 	// // err = dec.DecodeMulti(&b)
// 	// checkErr(err)
// 	// // slice, err := dec.DecodeSlice()
// 	// // checkErr(err)
// 	// return nil
// 	err := dec.DecodeMulti(&s.Key, &b)
// 	colors := make([]*Coolor, 0)
// 	for _, v := range b {
// 		// fmt.Printf("\n\n******%v ", v)
// 		// fmt.Printf("\n\n******%T", v)
// 		for _, vv := range v.(map[string]interface{}) {
// 			// fmt.Printf("\n\n******%T %T %v %v", kk, vv, kk, vv)
// 			c := &Coolor{
// 				Color: tcell.Color(vv.(uint64)),
// 			}
// 			colors = append(colors, c)
// 		}
// 	}
// 	s.Colors = colors
// 	// fmt.Println(colors)
// 	// fmt.Printf("\n***** %T %v", s, s)
// 	return err
// }

//
// func (s *Coolor) EncodeMsgpack(enc *msgpack.Encoder) error {
// 	return enc.EncodeUint64(uint64(s.Color))
// }
//
// func (s *Coolor) DecodeMsgpack(dec *msgpack.Decoder) error {
//   u64, err := dec.DecodeUint64()
//   if err != nil {
//     return err
//   }
// 	s.Color = tcell.Color(u64)
//   return nil
// }
//

// func (v Vector) MarshalBinary() ([]byte, error) {
// 	// A simple encoding: plain text.
// 	var b bytes.Buffer
// 	fmt.Fprintln(&b, v.x, v.y, v.z)
// 	return b.Bytes(), nil
// }

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
//
//	func (v *Vector) UnmarshalBinary(data []byte) error {
//		// A simple encoding: plain text.
//		b := bytes.NewBuffer(data)
//		_, err := fmt.Fscanln(b, &v.x, &v.y, &v.z)
//		return err
//	}
