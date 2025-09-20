build:
	./build-wasm.sh
	./build-zxing.sh
	./build-licenses.sh
	(cd ui && npm run build)
.PHONY: build

copy-to-e2e:
	cp ui/dist/index.html e2e/_index.html
.PHONY: copy-to-e2e

build-e2e: | build copy-to-e2e
.PHONY: build-e2e
