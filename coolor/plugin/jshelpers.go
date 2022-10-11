package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path"

	"github.com/evanw/esbuild/pkg/api"
	jsoniter "github.com/json-iterator/go"
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

// ToInterface JS->GO Convert *v8go.Value to Interface
func ToInterface(value *v8go.Value) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	var v interface{} = nil
	if value.IsNull() || value.IsUndefined() {
		return nil, nil
	} else if value.IsBigInt() {
		return value.BigInt(), nil
	} else if value.IsBoolean() {
		return value.Boolean(), nil
	} else if value.IsString() {
		return value.String(), nil
	}

	content, err := value.MarshalJSON()
	if err != nil {
		log.Printf("ToInterface MarshalJSON: %#v Error: %s", value, err.Error())
		return nil, err
	}

	err = jsoniter.Unmarshal([]byte(content), &v)
	if err != nil {
		log.Printf("ToInterface Unmarshal Value: %#v Content: %#v Error: %s", value, content, err.Error())
		return nil, err
	}
	return v, nil
}

// MustAnyToValue GO->JS Convert any to *v8go.Value
func MustAnyToValue(ctx *v8go.Context, value interface{}) *v8go.Value {
	v, err := AnyToValue(ctx, value)
	if err != nil {
		checkErr(err)
	}
	return v
}

// AnyToValue JS->GO Convert data to *v8go.Value
func AnyToValue(ctx *v8go.Context, value interface{}) (*v8go.Value, error) {
	switch value.(type) {
	case *v8go.Value:
		return value.(*v8go.Value), nil
	case []byte:
		// Todo: []byte to Uint8Array
		return v8go.NewValue(ctx.Isolate(), string(value.([]byte)))
	case string, int32, uint32, int64, uint64, bool, float64, *big.Int:
		return v8go.NewValue(ctx.Isolate(), value)
	case uint:
		return v8go.NewValue(ctx.Isolate(), uint32(value.(uint)))
	case uint16:
		return v8go.NewValue(ctx.Isolate(), uint32(value.(uint16)))
	case uint8:
		return v8go.NewValue(ctx.Isolate(), uint32(value.(uint8)))
	case int8:
		return v8go.NewValue(ctx.Isolate(), uint32(value.(int8)))
	case int16:
		return v8go.NewValue(ctx.Isolate(), uint32(value.(int16)))
	case int:
		return v8go.NewValue(ctx.Isolate(), int32(value.(int)))
	}
	err := fmt.Errorf("Unable to coerce %v (%T) to valid *v8go.Value ", value, value)
	log.Printf("AnyToValue error: %s", err)
	return nil, err
}

// ArrayToValuers GO->JS Convert []inteface to []v8go.Valuer
func ArrayToValuers(ctx *v8go.Context, values []interface{}) ([]v8go.Valuer, error) {
	res := []v8go.Valuer{}
	if ctx == nil {
		return res, fmt.Errorf("Context is nil")
	}

	for i := range values {
		value, err := AnyToValue(ctx, values[i])
		if err != nil {
			log.Printf("AnyToValue error: %s", err)
			value, _ = v8go.NewValue(ctx.Isolate(), nil)
		}
		res = append(res, value)
	}
	return res, nil
}

// ValuesToArray JS->GO Convert []*v8go.Value to []interface{}
func ValuesToArray(values []*v8go.Value) []interface{} {
	res := []interface{}{}
	for i := range values {
		var v interface{} = nil
		if values[i].IsNull() || values[i].IsUndefined() {
			res = append(res, nil)
			continue
		}

		v, err := ToInterface(values[i])
		if err != nil {
			log.Printf("ValuesToArray Value: %v Error: %s", err.Error(), values[i])
			res = append(res, nil)
			continue
		}

		res = append(res, v)

		// res = append(res, ToInterface(values[i]))
		// content, _ := values[i].MarshalJSON()
		// jsoniter.Unmarshal([]byte(content), &v)
		// res = append(res, v)
	}
	return res
}
func V8StringArrayToGoStringArray(Array *v8go.Value) ([]string, error) {
	// 将数组转换为对象
	ArrayObject, err := Array.AsObject()
	if err != nil {
		return nil, err
	}
	// 获取数组长度
	ArrayLength, err := ArrayObject.Get("length")
	if err != nil {
		return nil, err
	}
	var Arrays []string
	for i := 0; i < int(ArrayLength.Integer()); i++ {
		Array, err := ArrayObject.GetIdx(uint32(i))
		if err != nil {
			return nil, err
		}
		Arrays = append(Arrays, Array.String())
	}
	return Arrays, nil
}

func GoStructToV8Object(
	Context *v8go.Context,
	Struct any,
) (*v8go.Value, error) {
	StructJson, err := json.Marshal(Struct)
	if err != nil {
		return nil, err
	}

	ObjectValue, err := v8go.JSONParse(Context, string(StructJson))
	if err != nil {
		return nil, err
	}

	return ObjectValue, nil
}

func V8ObjectToGoStringMap(
	Context *v8go.Context,
	Object v8go.Valuer,
) (map[string]string, error) {
	ObjectValue, err := v8go.JSONStringify(Context, Object)
	if err != nil {
		return nil, err
	}

	var Map map[string]string
	err = json.Unmarshal([]byte(ObjectValue), &Map)
	if err != nil {
		return nil, err
	}
	return Map, nil
}

func GoArrayToV8Object[T []string | []int](
	Context *v8go.Context,
	Array T,
) (*v8go.Value, error) {
	ArrayJson, err := json.Marshal(Array)
	if err != nil {
		return nil, err
	}

	ObjectValue, err := v8go.JSONParse(Context, string(ArrayJson))

	return ObjectValue, err
}

func V8ObjectToGoStringAnyMap(
	Context *v8go.Context,
	Object v8go.Valuer,
) (map[string]any, error) {
	ObjectValue, err := v8go.JSONStringify(Context, Object)
	if err != nil {
		return nil, err
	}

	var Map map[string]any
	err = json.Unmarshal([]byte(ObjectValue), &Map)
	if err != nil {
		return nil, err
	}
	return Map, nil
}


func Bundle(jsfile string) (string, api.BuildResult) {
	outfile := fmt.Sprintf(".temp/%s.js", jsfile)
	defer os.Remove(outfile)
	wd, err := os.Getwd()
	checkErr(err)

	result := api.Build(api.BuildOptions{
		AbsWorkingDir: path.Join(wd, "js"),
		// NodePaths:         []string{"js/node_modules"},
		EntryPoints:       []string{fmt.Sprintf("%s.js", jsfile)},
		Bundle:            true,
		Format:            api.FormatIIFE,
		Outfile:           outfile,
		Write:             true,
		GlobalName:        "global",
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})

	if len(result.Errors) != 0 {
		return "", result
	}

	if bytes, err := ioutil.ReadFile(outfile); err != nil {
		return "", result
	} else {
		return string(bytes), result
	}
}
