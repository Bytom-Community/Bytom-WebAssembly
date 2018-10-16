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
#default build
GOOS=js GOARCH=wasm go build -o main.wasm
#mini build
GOOS=js GOARCH=wasm go build -tags=mini -o main.wasm 
```


## WebAssembly JS Function
##### mini build
>createKey\
resetKeyPassword \
signTransaction1

##### default build
>createKey \
resetKeyPassword \
createAccount \
createAccountReceiver \
signTransaction \
signTransaction1

