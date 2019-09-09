SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

falco_output_proto_path := pkg/api/output/falco_output.proto
falco_output_proto_pkg := $(dir ${falco_output_proto_path})
falco_output_go_pkg := $(shell basename ${falco_output_proto_pkg})
falco_output_proto_url ?= https://raw.githubusercontent.com/falcosecurity/falco/feat/grpc-server-poc/userspace/falco/falco_output.proto
falco_output_proto_sha256 ?= 242637535f207c242a7631e4c1364908ca2bf2d5632ccc673890e26ca7580aa6

.PHONY: protos clean_protos

protos: ${falco_output_proto_path}
	echo ${falco_output_proto_sha256} $< | sha256sum -c

${falco_output_proto_path}: clean_protos
	mkdir -p $(dir $@)
	curl -L -o $@ ${falco_output_proto_url}

clean_protos:
	rm -f ${falco_output_proto_path}

build:
	${PROTOC} -I ${falco_output_proto_pkg} ${falco_output_proto_path} --go_out=plugins=grpc,import_path=${falco_output_go_pkg}:${falco_output_proto_pkg}
