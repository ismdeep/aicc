# `make help`
.PHONY: help
help:
	@cat Makefile | grep '# `' | grep -v '@cat Makefile'

# `make build`
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o ../aicc_linux_amd64  -mod vendor -trimpath -ldflags '-s -w' .
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -o ../aicc_linux_arm64  -mod vendor -trimpath -ldflags '-s -w' .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../aicc_darwin_arm64 -mod vendor -trimpath -ldflags '-s -w' .
