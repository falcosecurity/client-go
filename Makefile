SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

PROTOS := pkg/api/schema/schema.proto pkg/api/output/output.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/feat/grpc-server-poc/userspace/falco/schema.proto https://raw.githubusercontent.com/falcosecurity/falco/feat/grpc-server-poc/userspace/falco/output.proto
PROTO_SHAS := 68a2bcf9c63c62b9cb1d6f8dff165f591241e09edaf034cf532513fa68809686 155f8376902ce6e792daa9574bce84c2b3863b6aa7c2b4a4c2261365b51de2e6

PROTO_DIRS := $(dir ${PROTOS})
PROTO_DIRS_INCLUDES := $(patsubst %/, -I %, ${PROTO_DIRS})

# $(1): the proto path
# $(2): the proto URL
# $(3): the proto SHA256
define download_rule
$(1):
	@rm -f $(1)
	@mkdir -p ${PROTO_DIRS}
	@curl --silent -Lo $(1) $(2)
	@echo $(3) $(1) | sha256sum -c
	@${PROTOC} ${PROTO_DIRS_INCLUDES} $(1) --go_out=plugins=grpc,paths=source_relative:$(dir $(1))
endef
$(foreach PROTO,$(PROTOS),\
	$(eval $(call download_rule,$(PROTO),$(firstword $(PROTO_URLS)),$(firstword $(PROTO_SHAS))))\
	$(eval PROTO_URLS := $(wordlist 2,$(words $(PROTO_URLS)),$(PROTO_URLS)))\
	$(eval PROTO_SHAS := $(wordlist 2,$(words $(PROTO_SHAS)),$(PROTO_SHAS)))\
)

.PHONY: protos
protos: ${PROTOS}

.PHONY: clean
clean: ${PROTO_DIRS}
	@rm -rf $^