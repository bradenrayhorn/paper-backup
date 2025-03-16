//go:build js && wasm

package main

import (
	"errors"
	"syscall/js"

	"github.com/bradenrayhorn/paper-backup/methods"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("paperBackup", js.FuncOf(paperBackup))
	js.Global().Set("paperBackupDecode", js.FuncOf(paperBackupDecode))

	<-c
}

func paperBackup(this js.Value, args []js.Value) any {
	dataUInt8Array := args[0]
	fileName := args[1].String()
	passphrase := args[2].String()

	if dataUInt8Array.Type() != js.TypeObject || dataUInt8Array.Get("constructor").Get("name").String() != "Uint8Array" {
		return makeJsError(errors.New("data must be uint8array"))
	}

	length := dataUInt8Array.Length()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataUInt8Array)

	qr, err := methods.FileBackupEncode(data, fileName, passphrase)
	if err != nil {
		return makeJsError(err)
	}

	// assemble response
	array := js.Global().Get("Uint8Array").New(len(qr))
	js.CopyBytesToJS(array, qr)
	return array
}

func paperBackupDecode(this js.Value, args []js.Value) any {
	key := args[0].String()
	dataUInt8Array := args[1]

	if dataUInt8Array.Type() != js.TypeObject || dataUInt8Array.Get("constructor").Get("name").String() != "Uint8Array" {
		return makeJsError(errors.New("data must be uint8array"))
	}

	length := dataUInt8Array.Length()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataUInt8Array)

	data, fileName, err := methods.FileBackupDecode(data, key)
	if err != nil {
		return makeJsError(err)
	}

	obj := map[string]any{}
	array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(array, data)

	obj["data"] = array
	obj["fileName"] = fileName

	return js.ValueOf(obj)
}

func makeJsError(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}
