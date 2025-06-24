package registry

import (
	"fmt"
	"reflect"

	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/responses"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/jinzhu/inflection"
)

// ModelInfo sekarang lebih sederhana
type ModelInfo struct {
	Instance any
	NewSlice func() any
}

var modelRegistry = make(map[string]ModelInfo)

func RegisterModel[T any](name string, model T) {
	modelType := reflect.TypeOf(model)
	modelRegistry[name] = ModelInfo{
		Instance: model,
		NewSlice: func() any {
			sliceType := reflect.SliceOf(modelType)
			return reflect.New(sliceType).Interface()
		},
	}
}

func RegisterAuto[T any](model T) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	modelName := modelType.Name()
	pluralSnake := inflection.Plural(utils.SnakeCase(modelName))

	// 1. Daftar model
	RegisterModel(pluralSnake, model)

	// 2. Otomatis cari dan daftar response struct jika ada
	responseType := tryFindResponseStruct(modelName)
	if responseType != nil {
		responses.Register(modelName, reflect.New(responseType).Elem().Interface())
	}
}

// Mencoba cari struct dengan nama <ModelName>Response di package responses
func tryFindResponseStruct(modelName string) reflect.Type {
	// Nama response struct yang dicari
	targetName := modelName + "Response"

	// Loop semua exported tipe di package responses (Go 1.21+ atau dengan known list)
	knownTypes := []any{
		responses.UserResponse{},
		// Tambah manual di sini jika perlu, atau generate otomatis
	}

	for _, t := range knownTypes {
		tType := reflect.TypeOf(t)
		if tType.Name() == targetName {
			return tType
		}
	}
	return nil
}

func GetModel(name string) (ModelInfo, error) {
	modelInfo, ok := modelRegistry[name]
	if !ok {
		return ModelInfo{}, fmt.Errorf("model not found: %s", name)
	}
	return modelInfo, nil
}

func init() {
	RegisterAuto(&models.User{})
}
