#!/bin/bash

ts_nocheck_file() {
    local file_path="$1"
    
    cat /dev/stdin "$file_path" <<EOI > "${file_path}.tmp"
// @ts-nocheck
EOI

    mv "${file_path}.tmp" "$file_path"
}

# copy wasm_exec.js from Go
mkdir -p ui/src/wasm/
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ui/src/wasm/

# Get base64 encoding of main.wasm
GOOS=js GOARCH=wasm go build -o main.wasm .
base64_wasm=$(cat main.wasm | base64 | tr -d '\n')

# Create the wasm_load.js file
cat > ui/src/wasm/load.js <<EOF
const wasmBinary = Uint8Array.from(atob('${base64_wasm}'), c => c.charCodeAt(0));
const go = new Go();
WebAssembly.instantiate(wasmBinary.buffer, go.importObject).then((result) => {
  go.run(result.instance);
});
EOF

# Create the wasm_load_worker.js file
cat > ui/src/wasm/load_worker.js <<EOF
import "./wasm_exec"
const wasmBinary = Uint8Array.from(atob('${base64_wasm}'), c => c.charCodeAt(0));
const go = new Go();
const wasmPromise = WebAssembly.instantiate(wasmBinary.buffer, go.importObject).then((result) => {
  go.run(result.instance);
});
export default wasmPromise;
EOF

cat > ui/src/wasm/load_worker.d.ts <<EOF
declare const wasmPromise: Promise<void>;
export default wasmPromise;
EOF

ts_nocheck_file "ui/src/wasm/wasm_exec.js"
ts_nocheck_file "ui/src/wasm/load.js"
ts_nocheck_file "ui/src/wasm/load_worker.js"

