package errorPack

import (
	"reflect"
)

func memberwiseClone(src any) any {
	v := reflect.ValueOf(src).Elem()
	vp2 := reflect.New(v.Type())
	vp2.Elem().Set(v)
	return vp2.Interface()
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
