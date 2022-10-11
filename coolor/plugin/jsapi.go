package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/goutil/errorx"
	"github.com/knadh/koanf"

	"rogchap.com/v8go"
)
type ObjFnCallback func(info *v8go.FunctionCallbackInfo) *v8go.Value
type ObjFn struct {
	tpl  *v8go.ObjectTemplate
	name string
	fn   *ObjFnCallback
}

type V8CoolorObjTpl struct {
	*v8go.ObjectTemplate
}
type V8CoolorObj struct {
	*v8go.Object
}
//
// type V8CoolorColor struct {
// 	*V8CoolorObj
// 	Color *Coolor
// }

func (tpl *V8CoolorObjTpl) RegisterObjFn(name string, fn ObjFnCallback) {
	tpl.Set(name, fn, v8go.ReadOnly)
}

// func MakeFnTpl(iso *v8go.Isolate, ctx *v8go.Context, fn func(args... interface{}) interface{}) *v8go.FunctionTemplate {
//   return v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
//     rfn := reflect.TypeOf(fn)
//
//   })
// }

func GetFunctionName(f interface{}) string {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Func {
		if rf := runtime.FuncForPC(v.Pointer()); rf != nil {
			return rf.Name()
		}
	}
	return v.String()
	// strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	// return strs[len(strs)-1]
}

// func (v8c *V8ColorObject) StructToV8ObjTpl() *v8go.Value {
// 	rv := reflect.ValueOf(v8c.Color)
// 	iso := v8c.Ctx.Isolate()
// 	tpl := v8go.NewObjectTemplate(iso)
//
// 	if rv.Kind() == reflect.Struct {
// 		vf := reflect.VisibleFields(rv.Type())
// 		for _, v := range vf {
// 			n := v.Name
// 			valIdx := v.Index
// 			valField := rv.FieldByIndex(valIdx)
// 			tpl.Set(n, valField.Interface())
// 		}
// 	}
//
// 	for i := 0; i < rv.NumMethod(); i++ {
// 		m := rv.Type().Method(i)
// 		mf := rv.Method(i)
//
// 		// dump.P(m)
//
// 		tpl.Set(
// 			m.Name,
// 			v8go.NewFunctionTemplate(
// 				iso,
// 				func(info *v8go.FunctionCallbackInfo) *v8go.Value {
// 					args := make([]reflect.Value, len(info.Args()))
//
// 					for i := 0; i < len(args); i++ {
// 						args[i] = reflect.ValueOf(info.Args()[i])
// 					}
// 					oout := mf.Call(args)
// 					out := make([]*v8go.Value, len(oout))
// 					outObj, err := v8go.NewObjectTemplate(iso).NewInstance(v8c.Ctx)
// 					checkErr(err)
// 					outObj.Set("length", len(oout))
// 					// no := m.Type.NumOut()
// 					for i, v := range oout {
// 						// dump.P(i,v,fmt.Sprintf("%v %T", v, v),v.Kind())
// 						out[i] = errAss(AnyToValue(v8c.Ctx, v.Interface()))
// 						checkErr(err)
// 						// v.Kind()
// 						outObj.SetIdx(uint32(i), out[i])
// 					}
// 					// dump.P(out)
//
// 					// val, _ := v8go.NewValue(iso, out)
// 					return outObj.Value
// 				},
// 			))
// 	}
//
// 	itpl, err := tpl.NewInstance(v8c.Ctx)
// 	checkErr(err)
// 	return itpl.Value
// }

type GoV8Env struct {
	idx      int
	VM       *v8go.Isolate
	Ctx      *v8go.Context
	Snapshot *v8go.StartupData
	Gbl      *v8go.ObjectTemplate
  Creator *v8go.SnapshotCreator
}

//  description
func NewGoV8Env(
	iso *v8go.Isolate,
	ctx *v8go.Context,
	snap *v8go.StartupData,
) *GoV8Env {
	gov8 := &GoV8Env{
		VM:       iso,
		Ctx:      ctx,
		Snapshot: snap,
		// Gbl:      &v8go.ObjectTemplate{},
	}

	if snap == nil {
		gov8.InitForSnapshot(func(gv8 *GoV8Env) error {
      gov8.runScript("js/init.js")
      gov8.runScript("js/_.js")
      return nil
    })
	}

	return gov8
}
//  description
func NewGoV8(
  f func(gv8 *GoV8Env) error,
) *GoV8Env {
	gov8 := &GoV8Env{
		idx:      0,
		VM:       &v8go.Isolate{},
		Ctx:      &v8go.Context{},
		Snapshot: &v8go.StartupData{},
		Gbl:      &v8go.ObjectTemplate{},
		Creator:  &v8go.SnapshotCreator{},
	}

  if f == nil {
    gov8.InitForSnapshot(func(gv8 *GoV8Env) error {
      gv8.runScript("js/init.js")
      gv8.runScript("js/_.js")
      return nil
    })
  } else  {
    f(gov8)
  }

	return gov8
}


type GoV8Class struct {
	o   interface{}
	VM  *v8go.Isolate
	Ctx *v8go.Context
	Gbl *v8go.ObjectTemplate
  Creator *v8go.SnapshotCreator
}

// func ColorConstructor(ctx *v8go.Context) v8go.FunctionCallback {
// 	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
// 		var col Color
// 		if len(info.Args()) > 0 {
// 			cstr := info.Args()[0].String()
// 			// log.Println(cstr)
// 			col, _ = Hex(cstr)
// 		}
// 		v8cc := &V8ColorObject{
// 			Color: &col,
// 			VM:    ctx.Isolate(),
// 			Ctx:   ctx,
// 			// Obj:   col,
// 		}
// 		ot := v8cc.StructToV8ObjTpl()
// 		return ot
// 	}
// }


func Watcher() (*fsnotify.Watcher, <-chan interface{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	// defer watcher.Close()

	done := make(chan bool)
	changed := make(chan interface{})
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("error: NOT OK")
					return
				}
					log.Printf("%s %s\n", event.Name, event.Op)
				// if event.Op == fsnotify.Write {
				// 	log.Printf("%s %s\n", event.Name, event.Op)
					changed <- event
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("error: NOT OK")
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// err = watcher.Add("./js")
	// if err != nil {
	// 	log.Fatal("Add failed:", err)
	// }
	return watcher, changed
}

func getMainScript(iso *v8go.Isolate, ctx *v8go.Context) (*v8go.Value, error) {
	main, err := ioutil.ReadFile("js/main.js")
	checkErrX(err)
	val, err := ctx.RunScript(string(main), "main.js")
	checkErrX(err)
	return val, err
	// if val.IsFunction() {
	//
	//   return val.AsFunction()
	// }
	return nil, fmt.Errorf("no function returned %v ", err)
}
func init() {
  // ansitags := tview.TranslateANSI(testTxt)
  // fmt.Println(ansitags)
}
func (gov8 *GoV8Env) InitForSnapshot(f func(gv8 *GoV8Env) error) {
  // tview.TranslateANSI
	snapshotCreator := v8go.NewSnapshotCreator()
  gov8.Creator = snapshotCreator
	snapshotCreatorIso, err := snapshotCreator.GetIsolate()
	checkErr(err)
	gov8.VM = snapshotCreatorIso

	snapshotCreatorCtx := v8go.NewContext(snapshotCreatorIso)
	gov8.Ctx = snapshotCreatorCtx

  err = f(gov8)
  checkErr(err)
}

func (gov8 *GoV8Env) Snap(def *v8go.Context) *v8go.StartupData {
	defer gov8.Creator.Dispose()
	defer gov8.Ctx.Close()
  defer gov8.VM.Dispose()
  err := gov8.Creator.SetDefaultContext(def)
	checkErr(err)
	data, err := gov8.Creator.Create(v8go.FunctionCodeHandlingClear)
	checkErr(err)
	// gov8.idx = ind
  // gov8.Creator.Dispose()
  gov8.Creator = nil
	gov8.Snapshot = data
	gov8.VM = nil
	gov8.Ctx = nil

  return data
}

func (gov8 *GoV8Env) RunFromSnapshot(path string, args... func(gv8 *GoV8Env)) {
	iso := v8go.NewIsolate(v8go.WithStartupData(gov8.Snapshot))
	defer iso.Dispose()
	gov8.VM = iso

	ctx := v8go.NewContext(iso)
	// checkErr(err)
	gov8.Ctx = ctx
	defer ctx.Close()

	gov8.DoBindings()

  if len(args) > 0 {
    args[0](gov8)
  }

	if gov8.runScript(path) {
		log.Printf("Script ran successfully%s", "")
	}
	// val, err := ctx.RunScript(source string, origin string)

	// runVal, err := ctx.Global().Get("run")
	// checkErr(err)
	//
	// fn, err := runVal.AsFunction()
	// checkErr(err)
	//
	// val, err := fn.Call(v8go.Undefined(iso))
	// checkErr(err)
	//
}

func (gov8 *GoV8Env) runScript(scriptPath string) bool {
	info, err := os.Stat(scriptPath)
	if checkErrX(err, info) {
		libFile, err := ioutil.ReadFile(scriptPath)
		checkErr(err)
		_, err = gov8.Ctx.RunScript(string(libFile), scriptPath)
		checkErrX(err)
		return true
	}
	return false
}

func (gov8 *GoV8Env) runEnvironmentLibraries() {
	gov8.runScript("js/init.js")
	gov8.runScript("js/_.js")
	// lodash, err := ioutil.ReadFile("js/_.js")
	// checkErr(err)
	// _, err = ctx.RunScript(string(lodash), "lodash.js")
	// checkErr(err)
	// err = ioutil.WriteFile("code.cache", data.CreateCodeCache().Bytes, fs.ModePerm)
	// checkErr(err)
	// data.Run(ctx)
}

func (gov8 *GoV8Env) DoBindings() {
	// cc := ColorConstructor(gov8.Ctx)
	// hsluvfn := v8go.NewFunctionTemplate(
	// 	gov8.VM,
	// 	func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	// 		if len(info.Args()) != 3 {
	// 			return nil
	// 		}
	//
	// 		hue := info.Args()[0].Int32()
	// 		s, err := strconv.ParseFloat(info.Args()[1].DetailString(), 64)
	// 		checkErr(err)
	// 		l, err := strconv.ParseFloat(info.Args()[2].DetailString(), 64)
	// 		checkErr(err)
	//
	// 		sf := float64(s / 100.0)
	// 		sl := float64(l / 100.0)
	//
	// 		hslvu := HSLuv(float64(hue), sf, sl)
	// 		fmt.Printf("%s", hslvu.GetCC().TerminalPreview())
	// 		// r,g,b := hslvu.RGB255()
	// 		// dump.P(hue,sf,sl)
	// 		v8cc := &V8ColorObject{
	// 			Color: &hslvu,
	// 			VM:    gov8.VM,
	// 			Ctx:   gov8.Ctx,
	// 			// Obj:   col,
	// 		}
	// 		ot := v8cc.StructToV8ObjTpl()
	// 		return ot
	// 		// return errAss[Color](GoStructToV8Object(ctx, hslvu))
	// 	},
	// )
	fetchfn := v8go.NewFunctionTemplate(gov8.VM, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
    
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			val, _ := v8go.NewValue(info.Context().Isolate(), string(body))
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value
	})
	dumpfn := v8go.NewFunctionTemplate(
		gov8.VM,
		func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			// arr := ValuesToArray(info.Args())
			arr := make([]string, len(info.Args()))
			for i, v := range info.Args() {
				arr[i] = v.DetailString()
			}
			// arr = append(arr, info.This().DetailString())
      fmt.Printf("%s", strings.Join(arr, " "))
			// log.Printf("%T %v", arr, arr)
			return nil
		},
	)
  // console := v8go.NewObjectTemplate(gov8.VM)
  // console.Set("log", dumpfn.GetFunction(gov8.Ctx), v8go.ReadOnly)
	// ccfn := v8go.NewFunctionTemplate(
	// 	gov8.VM,
	// 	func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	// 		return cc(info)
	// 	},
	// )
	// err := gov8.Ctx.Global().Set("Hsluv", hsluvfn.GetFunction(gov8.Ctx))
	// checkErr(err)
  err := gov8.Ctx.Global().Set("fetch", fetchfn.GetFunction(gov8.Ctx))
	checkErr(err)
  err = gov8.Ctx.Global().Set("dump", dumpfn.GetFunction(gov8.Ctx))
	checkErr(err)
  val, err := gov8.Ctx.RunScript(`
    const console={log:dump,warn:dump,error:dump};
    `, "console.js")
  _ = val
	checkErr(err)
	// err = gov8.Ctx.Global().Set("color", ccfn.GetFunction(gov8.Ctx))
}
func eajs[R any](v R, e error) R {
	iserr := func(v interface{}) {
		// if v == nil {
		// 	return
		// }
		e, ok := v.(*v8go.JSError)
		if ok {
			// zlog.Error(
			// 	fmt.Sprintf("%T %v", v, e),
			// 	zzlog.String("msg", e.Message),
			// 	zzlog.String("loc", e.Location),
			// )
			// doLog(errorx.Newf("%s %s %s %v", e.StackTrace, e.Message, e.Location, v))
			panic(e)
		}
	}
	iserr(e)
	return v
}

func doLog(args ...interface{}) {
	log.Printf("%v", args)
}

func checkErrX(err error, vars ...interface{}) bool {
	if err != nil {
		doLog(errorx.WithPrevf(errorx.Traced(err), "%T %v", err, vars))
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		doLog(err)
		panic(err)
	}
}


func RunJSForColorScheme(conf *koanf.Koanf){
  mapd := conf.All()
  keys := conf.Keys()

  _, result := Bundle("template")
  if len(result.Errors) > 0 {
    for _, v := range result.Errors {
      fmt.Printf("%s %v", v)
    }
    // fmt.Println(result.Errors)
  }
  if len(result.Warnings) > 0 {
    for _, v := range result.Warnings {
      fmt.Printf("%s %v", v)
    }
    // fmt.Println(result.Errors)
  }
  // fmt.Println(script, result)
	// iso := v8go.NewIsolate()
	// ctx := v8go.NewContext(iso)
	//  getCachedCode(iso, ctx)
	gov8env := NewGoV8(func(gv8 *GoV8Env) error {
    // gv8.runScript("js/init.js")
    // gv8.runScript("js/_.js")
    return nil
  })

  gov8env.InitForSnapshot(func(gv8 *GoV8Env) error {
    // gv8.runScript("js/init.js")
    // gv8.runScript("js/_.js")
    return nil
  })

  data := gov8env.Snap(gov8env.Ctx)
  // fmt.Println(data)
  v8env := NewGoV8(func(gv8 *GoV8Env) error {
    gv8.Snapshot=data
    return nil
  })
  _ = v8env

  gov8env.RunFromSnapshot("js/.temp/template.js", func(gv8 *GoV8Env) {
    str, _ := GoStructToV8Object(gv8.Ctx, mapd)
    kstr, _ := GoStructToV8Object(gv8.Ctx, keys)
    // ansiNames, _ := GoStructToV8Object(gv8.Ctx, baseXtermAnsiColorNames)
    // gv8.Ctx.Global().Set("xtermAnsiNames", ansiNames)
    gv8.Ctx.Global().Set("mapd", str)
    gv8.Ctx.Global().Set("keys", kstr)
  })

}

func jsapi() {
  script, result := Bundle("template")
  fmt.Println(script, result)
	// iso := v8go.NewIsolate()
	// ctx := v8go.NewContext(iso)
	//  getCachedCode(iso, ctx)
	gov8env := NewGoV8(func(gv8 *GoV8Env) error {
    // gv8.runScript("js/init.js")
    // gv8.runScript("js/_.js")
    return nil
  })

  gov8env.InitForSnapshot(func(gv8 *GoV8Env) error {
    gv8.runScript("js/init.js")
    gv8.runScript("js/_.js")
    return nil
  })

  data := gov8env.Snap(gov8env.Ctx)
  // fmt.Println(data)
  v8env := NewGoV8(func(gv8 *GoV8Env) error {
    gv8.Snapshot=data
    return nil
  })
  _ = v8env

	watcher, changed := Watcher()
	watcher.Add("js")
	watcher.Add("js/lib")
	watcher.Add("js/lib/schemes")
	defer watcher.Close()
  
  // ════════────── ── ─··  ·  ·
  debouncedChan := debounce(50*time.Millisecond, 200*time.Millisecond, changed)

	for {
		select {
		case <-debouncedChan:
      script, result := Bundle("template")
      fmt.Sprintf(script, result)
			gov8env.RunFromSnapshot("js/.temp/template.js")
			// checkErrX(err,fn)
		}
		// val, err := fn.Call(nil)
		// val, err := ctx.RunScript("shuffle()", "main.js")
		// checkErrX(err,val)
	}
}

func debounce[T any](min time.Duration, max time.Duration, input <-chan T) chan T {
	output := make(chan T)

	go func() {
		var (
			buffer   T
			ok       bool
			minTimer <-chan time.Time
			maxTimer <-chan time.Time
		)

		// Start debouncing
		for {
			select {
			case buffer, ok = <-input:
				if !ok {
					return
				}
				minTimer = time.After(min)
				if maxTimer == nil {
					maxTimer = time.After(max)
				}
			case <-minTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			case <-maxTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			}
		}
	}()

	return output
}
// func NewSetObjectFunctions(
// 	iso *v8go.Isolate,
// 	objTpl *v8go.ObjectTemplate,
// ) SetObjectFunctions {
// 	setObjectFunction := func(name string, fnCallback v8go.FunctionCallback) error {
// 		errorFnTpl, err := v8go.NewFunctionTemplate(iso, fnCallback)
// 		if err != nil {
// 			return fmt.Errorf("failed to create %s FunctionTemplate: %v", name, err)
// 		}
//
// 		if err := objTpl.Set(name, errorFnTpl, v8go.ReadOnly); err != nil {
// 			return fmt.Errorf("failed to set %s function: %v", name, err)
// 		}
//
// 		return nil
// 	}
// _, err = ctx.RunScript(`
// 	var c = color('#abcdef');
//    var data = [];
//    var h = _.random(0,360);
//    var H = _.map([0,45,90,135,180,225,270,315,360], function(offset){
//      return (h+offset) % 360;
//    });
//    var backS = _.random(5, 40);
//    var darkL = _.random(0, 20);
//    var rangeL = 90 - darkL;
//    var cols = {};
//    cols.bg = Hsluv(H[0], backS, darkL / 2).Hex();
//    cols.fg = Hsluv(H[0], backS, rangeL).Hex();
//    cols.cursor = Hsluv(H[0], backS, (darkL + rangeL)/2).Hex();
//    for (var i = 1; i < 8; i++) {
//      // var h = H[_.random(0, H.length-1)];
//    // var d = (H[i] + _.random(-11,11)) % 360;
//    var backS = _.random(5, 40);
//    var darkL = _.random(0, 20);
//    var rangeL = 90 - darkL;
//      var s = [H[i], backS, darkL + rangeL * Math.pow(i/8, 1.3)];
//      var uv = Hsluv(s[0],s[1],s[2]);
//      var hex = uv.Hex();
//      cols[` + "`${i}color`" + `] = hex;
//      data.push(hex)
//    }
//    // 8 Random shades
//    var minS = _.random(30, 70);
//    var maxS = minS + 30;
//    var minL = _.random(50,70);
//    var maxL = minL + 20;
//      var sr = _.random(minS, maxS);
//      var l = _.random(minL, maxL);
//    for (var j = 8; j < 16; j++) {
//    // var h = (H[j % H.length] + _.random(-11,11)) % 360;
//      // var h = H[_.random(0, H.length-1)];
//      var s = [H[(j % H.length)], sr, l ];
//      var uv = Hsluv(s[0],s[1],s[2]);
//      // var uv = new Hsluv(h,s,l);
//      var hex = uv.Hex();
//      cols[` + "`${j}color`" + `] = hex;
//    }
//    // Object.keys(cols).map(x => {dump(x);dump(cols[x]);return cols[x]})
//  `,"console.js")
// checkErr(err)
//
// 	return func(fnMap map[string]v8go.FunctionCallback) error {
// 		for name, fn := range fnMap {
// 			if err := setObjectFunction(name, fn); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// }
// func shite() {
//     // jsfnctx, err := js.NewContext()
//   // checkErr(err)
//   // iso,err := js.NewIsolate()
//   defer iso.Dispose()
//   checkErr(err)
//   // poly, err := js.NewPolyfill(iso)
//   // checkErr(err)
//   global, err := v8go.NewObjectTemplate(iso)
//   checkErr(err)
//   fn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
//     var val *v8go.Value
//     val, err := v8go.NewValue(iso, "SHIT")
//     checkErr(err)
//     return val.Object().Value
//     // strs := lo.Map[*v8go.Value,string](info.Args(), func(v *v8go.Value, i int) string {
//     //   return v.String()
//     // })
//     // // i, err := info.Context().Isolate()
//     // checkErr(err)
//     // val, err := v8go.NewValue(iso, strs)
//     // checkErr(err)
//     // return val
//   })
//   checkErr(err )
//   // v8go.NewContext(v8go.ContextOption{})
//   err = global.Set("tester", fn)
//   checkErr(err)
//   // pf, err := js.NewPolyfill(iso)
//   // checkErr(err)
//   nctx, err := v8go.NewContext(global)// js.NewContext(global)
//   checkErr(err)
//   dump.P(iso.GetHeapStatistics())
//   val, err := nctx.RunScript("tester('ass','hole')", "script.js")
//   checkErr(err)
//   dump.P(val.IsString(), iso.GetHeapStatistics())
//
//
// }
// "#060306", "#2b1a29", "#492e45", "#694463", "#8a5b82", "#ad73a4", "#c494bc", "#d7b8d1", "#a5a767", "#c3c66b", "#e18d7c", "#efb3a8", "#de8aa9", "#eeb0c5", "#79b167", "#86d46a",
