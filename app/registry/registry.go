package registry

import (
	"fmt"
	"goravel/app/models"
	"reflect"
)

// ModelInfo sekarang lebih sederhana
type ModelInfo struct {
	Instance any
	NewSlice func() any
}

var modelRegistry = make(map[string]ModelInfo)

// RegisterModel sekarang tidak butuh parameter postProcess
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

func GetModel(name string) (ModelInfo, error) {
	modelInfo, ok := modelRegistry[name]
	if !ok {
		return ModelInfo{}, fmt.Errorf("model not found: %s", name)
	}
	return modelInfo, nil
}

// Inisialisasi sekarang jauh lebih bersih
func init() {
	RegisterModel("users", &models.User{})
	// RegisterModel("posts", &models.Post{}) // Contoh mendaftarkan model lain
	// RegisterModel("products", &models.Product{})
}
