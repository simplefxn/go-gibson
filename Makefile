BINARY=go-gibson
BINARY_DIR=bin
GOARCH=amd64

all: linux darwin windows

linux: mod
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-linux-${GOARCH} main.go

darwin: mod
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-darwin-${GOARCH} main.go

windows: mod
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-windows-${GOARCH}.exe main.go

mod: 
	go mod tidy

.PHONY: proto 
proto:
	protoc -I=./proto --go_out=./pkg ./proto/*.proto

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

.PHONY: dev 
dev: update-go-deps push