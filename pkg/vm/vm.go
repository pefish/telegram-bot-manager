package vm

import (
	"fmt"
	"github.com/dop251/goja"
)

func NewVm() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	//vm.Set("test", func(call goja.FunctionCall) goja.Value {
	//	result := vm.ToValue(2 + call.Argument(0).ToInteger())
	//	return result
	//})
	vm.Set("console", &struct{
		Log func (data interface{}) `json:"log"`
	}{
		Log: func(data interface{}) {
			fmt.Println(data)
		},
	})
	return vm
}

