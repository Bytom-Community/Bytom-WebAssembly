package main

import (
	"github.com/bytom-community/wasm/sdk"
)

func main() {
	done := make(chan struct{}, 0)
	sdk.Register()
	<-done
}
