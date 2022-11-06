package coolor

import (
	"os"
	"time"
)

//
//
func NewCoolorColorsPaletteMeta(
	name string,
	ccp *CoolorColorsPalette,
) *CoolorColorsPaletteMeta {
  if ccp == nil {
    return nil
  }
	now := time.Now()
	ccm := &CoolorColorsPaletteMeta{
		Created: now,
		Name:    name,
		Author:  os.ExpandEnv("$USER"),
		ExtraData:    make(map[string]interface{}),
		Palette: ccp,
	}
	return ccm
}
// func (ccs CoolorColorsPaletteMeta) String() string {
// 	str := ""
// 	str = fmt.Sprintf(
// 		"%d %s %s",
// 		ccs.ID,
// 		ccs.Named,
// 		ccs.Started.Format(time.RFC3339),
// 	)
// 	if ccs.Current == nil {
// 		return "nil"
// 	}
// 	for _, v := range ccs.Current.Colors {
// 		str = fmt.Sprintf("%s %s", str, v.Escalate().TerminalPreview())
// 	}
// 	return str
// 	// return fmt.Sprintf("%s %s", )
// }
// func (cc *CoolorColorsPaletteMeta) Update() {
// 	// var err error
// 	//  if cc == nil || cc.Current == nil || len(cc.Current.Colors) == 0 {
// 	//    return
// 	//  }
// 	// if cc.ID == 0 {
// 	//    var ccpm CoolorColorsPaletteMeta
// 	//    q := bh.Where("Named").Eq(cc.Named)
// 	//    err = Store.FindOne(&ccpm, q)
// 	//    if err != nil && err == bh.ErrNotFound {
// 	//      err = Store.Insert(bh.NextSequence(), cc)
// 	//      if checkErrX(err) {
// 	//        // fmt.Println(cc.ID)
// 	//        // cc = &ccpm
// 	//      }
// 	//    } else if err != nil {
// 	//      checkErrX(err)
// 	//    } else {
// 	//      cc = &ccpm
// 	//    }
// 	// }
// 	// 	err = Store.Upsert(cc.ID, cc)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// }
// func (cc CoolorColorsPaletteMeta) GetMeta() interface{} {
// 	return cc
// }
