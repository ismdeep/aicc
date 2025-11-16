# `make help`
.PHONY: help
help:
	@cat Makefile | grep '# `' | grep -v '@cat Makefile'

# `make build`
.PHONY: build
build:
	GOOS=linux  GOARCH=amd64 go build -o ../aicc_linux_amd64  -mod vendor -trimpath -ldflags '-s -w' .
	GOOS=linux  GOARCH=arm64 go build -o ../aicc_linux_arm64  -mod vendor -trimpath -ldflags '-s -w' .
	GOOS=darwin GOARCH=arm64 go build -o ../aicc_darwin_arm64 -mod vendor -trimpath -ldflags '-s -w' .
