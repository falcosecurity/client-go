falco_output_proto_path := pkg/api/output/falco_output.proto
falco_output_proto_pkg := pkg/api/output/falco_output
falco_output_proto_url ?= https://raw.githubusercontent.com/falcosecurity/falco/feat/grpc-server-poc/userspace/falco/falco_output.proto
falco_output_proto_sha256 ?= 242637535f207c242a7631e4c1364908ca2bf2d5632ccc673890e26ca7580aa6

.PHONY: protos clean_proto check_proto

protos: check_proto ${falco_output_proto_path} check_proto

check_proto: ${falco_output_proto_path}
	echo ${falco_output_proto_sha256} $< | sha256sum -c

${falco_output_proto_path}: clean_proto
	curl -L -o $@ ${falco_output_proto_url}
	echo ${falco_output_proto_sha256} $@ | sha256sum -c

clean_proto:
	rm ${falco_output_proto_path}

build:
	protoc-gen-go -I ${falco_output_proto_path} --go_out=plugins=grpc:${falco_output_proto_pkg}
