.PHONY: proto 
proto:
	protoc -I=./proto --go_out=./pkg ./proto/*.proto