package coolor

import (
	"fmt"

	"github.com/bfanger/enhanced-v8go/js"
	"github.com/gookit/goutil/dump"
	"rogchap.com/v8go"
)



func AddObjectToGloabl(ctx *v8go.Context, tpl *v8go.ObjectTemplate, name string) error {
	consoleObj, err := tpl.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("failed to create %s instance: %v", name, err)
	}

	global := ctx.Global()
	if err := global.Set(name, consoleObj); err != nil {
		return fmt.Errorf("failed to set %s object to global: %v", name, err)
	}

	return nil
}

type SetObjectFunctions func(fnMap map[string]v8go.FunctionCallback) error

func NewSetObjectFunctions(
	iso *v8go.Isolate,
	objTpl *v8go.ObjectTemplate,
) SetObjectFunctions {
	setObjectFunction := func(name string, fnCallback v8go.FunctionCallback) error {
		errorFnTpl, err := v8go.NewFunctionTemplate(iso, fnCallback)
		if err != nil {
			return fmt.Errorf("failed to create %s FunctionTemplate: %v", name, err)
		}

		if err := objTpl.Set(name, errorFnTpl, v8go.ReadOnly); err != nil {
			return fmt.Errorf("failed to set %s function: %v", name, err)
		}

		return nil
	}

	return func(fnMap map[string]v8go.FunctionCallback) error {
		for name, fn := range fnMap {
			if err := setObjectFunction(name, fn); err != nil {
				return err
			}
		}
		return nil
	}
}
func shite() {
    // jsfnctx, err := js.NewContext()
  // checkErr(err)
  iso,err := js.NewIsolate()
  defer iso.Dispose()
  checkErr(err)
  // poly, err := js.NewPolyfill(iso)
  // checkErr(err)
  global, err := v8go.NewObjectTemplate(iso)
  checkErr(err)
  fn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
    var val *v8go.Value
    val, err := v8go.NewValue(iso, "SHIT")
    checkErr(err)
    return val.Object().Value
    // strs := lo.Map[*v8go.Value,string](info.Args(), func(v *v8go.Value, i int) string {
    //   return v.String()
    // })
    // // i, err := info.Context().Isolate()
    // checkErr(err)
    // val, err := v8go.NewValue(iso, strs)
    // checkErr(err)
    // return val
  })
  checkErr(err )
  // v8go.NewContext(v8go.ContextOption{})
  err = global.Set("tester", fn)
  checkErr(err)
  // pf, err := js.NewPolyfill(iso)
  // checkErr(err)
  nctx, err := v8go.NewContext(global)// js.NewContext(global)
  checkErr(err)
  dump.P(iso.GetHeapStatistics())
  val, err := nctx.RunScript("tester('ass','hole')", "script.js")
  checkErr(err)
  dump.P(val.IsString(), iso.GetHeapStatistics())
  

}
