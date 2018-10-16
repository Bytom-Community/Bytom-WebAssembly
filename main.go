package main

import "github.com/bytom-community/wasm/sdk/js"

func main() {
	done := make(chan struct{}, 0)
	js.Register()
	<-done
}
