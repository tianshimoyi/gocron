PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
PROTOC_FLAGS:=-I. -I$(PROTOC_INC_PATH)

GOPATH:=$(shell go env GOPATH)
PWD:=$(shell pwd)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  tools             to install depends tools"
	@echo "  help              to show this help informatrion"

generate: $(PROTO_FILES) Makefile
	mkdir -p ../pkg/pb
	rm -rf ../pkg/pb/*

	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	@echo "ok"