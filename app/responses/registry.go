package responses

import "reflect"

var Registry = map[string]reflect.Type{}

func Register(modelName string, responseStruct any) {
	Registry[modelName] = reflect.TypeOf(responseStruct)
}

func Get(modelName string) (reflect.Type, bool) {
	t, ok := Registry[modelName]
	return t, ok
}
