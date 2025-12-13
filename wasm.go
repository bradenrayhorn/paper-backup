//go:build js && wasm

package main

import (
	"errors"
	"strings"
	"syscall/js"

	"github.com/bradenrayhorn/paper-backup/methods/filebackup"
	"github.com/bradenrayhorn/paper-backup/methods/shamirsecret"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("paperBackup", wasmHandler(paperBackup))
	js.Global().Set("paperBackupDecode", wasmHandler(paperBackupDecode))

	js.Global().Set("paperShamirSecretSplit", wasmHandler(paperShamirSecretSplit))
	js.Global().Set("paperShamirSecretCombineFromQR", wasmHandler(paperShamirSecretCombineFromQR))
	js.Global().Set("paperShamirSecretCombineFromText", wasmHandler(paperShamirSecretCombineFromText))

	<-c
}

func paperBackup(this js.Value, args []js.Value) (any, error) {
	dataUInt8Array := args[0]
	fileName := args[1].String()
	passphrase := args[2].String()

	if dataUInt8Array.Type() != js.TypeObject || dataUInt8Array.Get("constructor").Get("name").String() != "Uint8Array" {
		return nil, errors.New("data must be uint8array")
	}

	length := dataUInt8Array.Length()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataUInt8Array)

	qr, err := filebackup.Encode(data, fileName, passphrase)
	if err != nil {
		return nil, err
	}

	// assemble response
	array := js.Global().Get("Uint8Array").New(len(qr))
	js.CopyBytesToJS(array, qr)
	return array, nil
}

func paperBackupDecode(this js.Value, args []js.Value) (any, error) {
	key := args[0].String()
	dataUInt8Array := args[1]

	if dataUInt8Array.Type() != js.TypeObject || dataUInt8Array.Get("constructor").Get("name").String() != "Uint8Array" {
		return nil, errors.New("data must be uint8array")
	}

	length := dataUInt8Array.Length()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataUInt8Array)

	data, fileName, err := filebackup.Decode(data, key)
	if err != nil {
		return nil, err
	}

	obj := map[string]any{}
	array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(array, data)

	obj["data"] = array
	obj["fileName"] = fileName

	return js.ValueOf(obj), nil
}

func paperShamirSecretSplit(this js.Value, args []js.Value) (any, error) {
	secret := args[0].String()
	parts := args[1].Int()
	threshold := args[2].Int()

	passphrase, err := shamirsecret.RandomPassphrase()
	if err != nil {
		return nil, err
	}

	toPrint, err := shamirsecret.Encode(secret, standardizePassphrase(passphrase), parts, threshold)
	if err != nil {
		return nil, err
	}

	jsQRShares := js.Global().Get("Array").New(len(toPrint.QRShares))
	for i, share := range toPrint.QRShares {
		array := js.Global().Get("Uint8Array").New(len(share))
		js.CopyBytesToJS(array, share)
		jsQRShares.SetIndex(i, array)
	}

	jsTextShares := js.Global().Get("Array").New(len(toPrint.TextShares))
	for i, share := range toPrint.TextShares {
		jsTextShares.SetIndex(i, js.ValueOf(share))
	}

	obj := map[string]any{}
	obj["passphrase"] = passphrase
	obj["textShares"] = jsTextShares
	obj["qrShares"] = jsQRShares

	return js.ValueOf(obj), nil
}

func paperShamirSecretCombineFromQR(this js.Value, args []js.Value) (any, error) {
	passphrase := args[0].String()
	shares := [][]byte{}
	for i := 1; i < len(args); i++ {
		if args[i].Type() != js.TypeObject || args[i].Get("constructor").Get("name").String() != "Uint8Array" {
			return nil, errors.New("data must be uint8array")
		}

		data := make([]byte, args[i].Length())
		js.CopyBytesToGo(data, args[i])

		shares = append(shares, data)
	}

	result, err := shamirsecret.DecodeFromQR(shares, standardizePassphrase(passphrase))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func paperShamirSecretCombineFromText(this js.Value, args []js.Value) (any, error) {
	passphrase := args[0].String()
	shares := []string{}
	for i := 1; i < len(args); i++ {
		shares = append(shares, args[i].String())
	}

	result, err := shamirsecret.DecodeFromText(shares, standardizePassphrase(passphrase))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func standardizePassphrase(passphrase string) string {
	return strings.ToUpper(
		strings.TrimSpace(
			strings.ReplaceAll(passphrase, " ", ""),
		),
	)
}

func wasmHandler(handler func(this js.Value, args []js.Value) (any, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		promiseConstructor := js.Global().Get("Promise")
		promise := js.FuncOf(func(_this js.Value, _args []js.Value) any {
			resolve := _args[0]
			reject := _args[1]

			res, err := handler(this, args)
			if err != nil {
				obj := js.Global().Get("Object").New()
				obj.Set("error", err.Error())
				reject.Invoke(obj)
			} else {
				resolve.Invoke(res)
			}

			return nil
		})

		return promiseConstructor.New(promise)
	})
}
