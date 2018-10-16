// +build mini

package js

import (
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/base"
)

//RegisterFunc Register js func
type RegisterFunc func(args []js.Value)

var funcs map[string]RegisterFunc

func init() {
	funcs = make(map[string]RegisterFunc)

	funcs["createKey"] = base.CreateKey
	funcs["resetKeyPassword"] = base.ResetKeyPassword
	funcs["signTransaction1"] = base.SignTransaction1
}

//Register Register func
func Register() {
	jsFuncVal := js.Global().Get("AllFunc")
	for k, v := range funcs {
		call := js.NewCallback(v)
		jsFuncVal.Set(k, call)
	}
	setPrintMessage := js.Global().Get("setFuncOver")
	setPrintMessage.Invoke()
}
