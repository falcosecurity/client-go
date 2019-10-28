SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

PROTOS := pkg/api/schema/schema.proto pkg/api/output/output.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/dev/userspace/falco/schema.proto https://raw.githubusercontent.com/falcosecurity/falco/dev/userspace/falco/output.proto
PROTO_SHAS := a1f427c114b945d0880b55058862b74015d036aa722985ca6e5474ab4ed19f69 283a717cf891edb3247f0afe640dc625c421cee31a66c1763873c9f32352dee0

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
