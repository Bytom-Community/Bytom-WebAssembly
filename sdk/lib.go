package sdk

import "syscall/js"

func endFunc(value js.Value) {
	value.Call("endFunc")
}

func isEmpty(value string) bool {
	if value == "" || value == "undefined" {
		return true
	}
	return false
}
