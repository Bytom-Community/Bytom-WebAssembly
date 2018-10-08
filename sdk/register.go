package sdk

import "syscall/js"

//RegisterFunc Register js func
type RegisterFunc func(args []js.Value)

var funcs map[string]RegisterFunc

func init() {
	funcs = make(map[string]RegisterFunc)

	funcs["scMulBase"] = scMulBase
	funcs["createKey"] = createKey
	funcs["resetKeyPassword"] = resetKeyPassword
	funcs["createAccount"] = createAccount
	funcs["createAccountReceiver"] = createAccountReceiver
	funcs["signTransaction"] = signTransaction
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
