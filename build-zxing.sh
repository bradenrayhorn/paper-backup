#!/bin/bash

ts_nocheck_file() {
    local file_path="$1"
    
    cat /dev/stdin "$file_path" <<EOI > "${file_path}.tmp"
// @ts-nocheck
EOI

    mv "${file_path}.tmp" "$file_path"
}

set -e

EMSCRIPTEN_VERSION=4.0.5
EMSDK_DIR="emsdk"

BASE=$(pwd)
WORKING=$(pwd)/zxing

cd $WORKING

# these steps are mostly based on zxing-cpp WASM README.


# 1. install emsdk
if [ ! -d "$EMSDK_DIR" ]; then
	git clone https://github.com/emscripten-core/emsdk.git
fi

cd $WORKING/$EMSDK_DIR

./emsdk install $EMSCRIPTEN_VERSION
./emsdk activate $EMSCRIPTEN_VERSION
source ./emsdk_env.sh

cd $WORKING

# 2. init cmake project
mkdir -p $WORKING/build
cd $WORKING/build
emcmake cmake ../zxing-cpp/wrappers/wasm

# 3. apply patch
cd $WORKING/zxing-cpp
git restore .
git apply $WORKING/patches/BarcodeWriter.patch


# 4. build for wasm
cd $WORKING/build
export EMCC_CFLAGS="-O2 -Oz -s SINGLE_FILE=1 -s EXPORT_ES6=1"
cmake --build .

# 5. copy to ui project
mkdir -p $BASE/ui/src/wasm/
ts_nocheck_file "$WORKING/build/zxing.js"
cp $WORKING/build/zxing.js $BASE/ui/src/wasm/

# add types
cat > $BASE/ui/src/wasm/zxing.d.ts <<EOF
interface ZXingInstance {
  generateQRCodeFromBinary(data: Uint8Array, encoding: string, margin: number, width: number, height: number, eccLevel: number): WriteResult;

  readBarcodeFromPixmap(buffer: number, width: number, height: number, tryHard: boolean, format: string): ReadResult

  HEAPU8: HeapU8
  _malloc(length: number): number
  _free(id: number): void
}

interface HeapU8 {
  set(data: Uint8Array | Uint8ClampedArray, buffer: number): void
}

interface WriteResult {
  svg: string
  error: string
  delete: () => void
}

interface ReadResult {
  bytes: Uint8Array|null
  format: string
  error: string
  text: string
  position: {
    bottomLeft: Coordinate
    bottomRight: Coordinate
    topLeft: Coordinate
    topRight: Coordinate
  }
}

interface Coordinate {
  x: number
  y: number
}

declare function ZXing(): Promise<ZXingInstance>;

export default ZXing;
EOF

# cleanup.. reset patch
cd $WORKING/zxing-cpp
git restore .

