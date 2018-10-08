# Bytom-WebAssembly
It is a project for Bytom WebAssembly

## Prepare
```sh
git clone https://github.com/Bytom-Community/Bytom-WebAssembly.git $GOPATH/src/github.com/bytom-community/wasm
```

## Build

Need Go version 1.11

```sh
cd $GOPATH/src/github.com/bytom-community/wasm
GOOS=js GOARCH=wasm go build -o main.wasm
```
