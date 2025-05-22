#!/bin/bash

# copy license from Go
mkdir -p ui/src/license/

file=ui/src/license/go.txt

echo -e "GO VERSION: $(go version | awk '{print $3}')\n\nLICENSE:\n" > "$file"
cat $(go env GOROOT)/LICENSE >> "$file"

# copy license from zxing-cpp
file=ui/src/license/zxing-cpp.txt
version=$(cd zxing/zxing-cpp && git rev-parse HEAD)

cat <<EOF > "$file"
zxing-cpp

REPOSITORY: https://github.com/zxing/zxing-cpp
VERSION: $version

LICENSE:

EOF
cat zxing/zxing-cpp/LICENSE >> $file

# special thank you
cat <<EOF > "ui/src/license/thankyou.txt"
Thank you to mjeanroy on GitHub for the rollup-plugin-license project that helps
make this license file possible.
https://github.com/mjeanroy/rollup-plugin-license

EOF

# vault
cat <<EOF > "ui/src/license/vault.txt"
Vault /shamir module.

REPOSITORY: https://github.com/hashicorp/vault

LICENSE:

EOF
cat shamir/LICENSE >> ui/src/license/vault.txt

# npm dependencies
touch ui/src/license/THIRDPARTY
(cd ui && GENERATE_LICENSES=true npm run build)

# assembly

output=ui/src/license/THIRDPARTY
> "$output"

for file in "ui/src/license"/*; do
  [ "$file" = "$output" ] && continue
  cat "$file" >> "$output"
  echo -e "\n\n\n============================================================\n\n\n" >> "$output"
done
