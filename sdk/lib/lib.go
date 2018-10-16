package lib

import "syscall/js"

//EndFunc end call go func
func EndFunc(value js.Value) {
	value.Call("endFunc")
}

//IsEmpty js value check empty
func IsEmpty(value string) bool {
	if value == "" || value == "undefined" {
		return true
	}
	return false
}
