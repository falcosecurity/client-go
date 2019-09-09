SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

PROTOS := pkg/api/output/falco_output.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/feat/grpc-server-poc/userspace/falco/falco_output.proto
PROTO_SHAS := 242637535f207c242a7631e4c1364908ca2bf2d5632ccc673890e26ca7580aa6

# $(1): the proto path
# $(2): the proto URL
# $(3): the proto SHA256
define download_rule
.PHONY: $(1)
$(1):
	@rm -f $(1)
	@mkdir -p $(dir $(1))
	@curl --silent -Lo $(1) $(2)
	@echo $(3) $(1) | sha256sum -c
	@${PROTOC} -I $(dir $(1)) $(1) --go_out=plugins=grpc,import_path=$(shell basename $(dir $(1))):$(dir $(1))
endef
$(foreach PROTO,$(PROTOS),\
	$(eval $(call download_rule,$(PROTO),$(firstword $(PROTO_URLS)),$(firstword $(PROTO_SHAS))))\
	$(eval PROTO_URLS := $(wordlist 2,$(words $(PROTO_URLS)),$(PROTO_URLS)))\
	$(eval PROTO_SHAS := $(wordlist 2,$(words $(PROTO_SHAS)),$(PROTO_SHAS)))\
)