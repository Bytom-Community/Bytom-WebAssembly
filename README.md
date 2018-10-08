# Bytom-WebAssembly
It is a project for Bytom WebAssembly

## build
```sh
govendor sync
GOOS=js GOARCH=wasm go build -o main.wasm
```
