SHELL := /bin/bash

PROTOC ?= $(shell which protoc)

PROTOS := pkg/api/schema/schema.proto pkg/api/output/output.proto pkg/api/version/version.proto
PROTO_URLS := https://raw.githubusercontent.com/falcosecurity/falco/dev/userspace/falco/schema.proto https://raw.githubusercontent.com/falcosecurity/falco/dev/userspace/falco/output.proto https://raw.githubusercontent.com/falcosecurity/falco/e1a753f792b93903fff9bdaaf02b336f672cb8c7/userspace/falco/version.proto
PROTO_SHAS := a1f427c114b945d0880b55058862b74015d036aa722985ca6e5474ab4ed19f69 4ce2f3e6d6ebc07a74535c4f21da73e44c6ef848ab83627b1ac987058be5ece9 0ac29f477e146a14504a34f91509668e44f4ee107d37b0be57e82847e8e89510

PROTO_DIRS := $(dir ${PROTOS})
PROTO_DIRS_INCLUDES := $(patsubst %/, -I %, ${PROTO_DIRS})

.PHONY: protos
protos: ${PROTOS}

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

.PHONY: clean
clean: ${PROTO_DIRS}
	@rm -rf $^
