package coolor

import (
	"encoding/json"

	"rogchap.com/v8go"
)

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

func GoStructToV8Object(Context *v8go.Context, Struct any) (*v8go.Value, error) {
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
func V8ObjectToGoStringMap(Context *v8go.Context, Object v8go.Valuer) (map[string]string, error) {
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
func GoArrayToV8Object[T []string | []int](Context *v8go.Context, Array T) (*v8go.Value, error) {
	ArrayJson, err := json.Marshal(Array)
	if err != nil {
		return nil, err
	}

	ObjectValue, err := v8go.JSONParse(Context, string(ArrayJson))

	return ObjectValue, err
}
func V8ObjectToGoStringAnyMap(Context *v8go.Context, Object v8go.Valuer) (map[string]any, error) {
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
