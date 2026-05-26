package excel

import (
	"reflect"

	fbsutils "github.com/arisu-archive/bluearchive-fbs-utils"
)

var fbs = map[string]reflect.Type{
}

func GetFlatDataByName(name string) fbsutils.FlatData {
	if data, ok := fbs[name]; ok {
		return reflect.New(data).Interface().(fbsutils.FlatData)
	}
	return nil
}
