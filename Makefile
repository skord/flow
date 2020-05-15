# Common-config is configuration that's common to the gazette/core repo,
# and also to consumer applications using Gazette's Make build system
# (as this application does).
include $(shell go list -f '{{ .Dir }}' -m go.gazette.dev/core)/mk/common-config.mk

# Next we declare project-specific configuration, and in particular configuration
# which is accessed by build rules in common-build.mk.

# 'protoc' invocations should also include the gazette/core repository.
PROTOC_INC_MODULES += "go.gazette.dev/core"

# protobuf-targets is the list of source files which are to be generated by protoc.
protobuf-targets = ./go/protocol/flow.pb.go

# ci-release-ping-pong lists target files or binaries which are to be packaged
# into a Docker image 'ping-pong:latest' by the 'ci-release-%' build rule.
#  ci-release-ping-pong-targets = \
#	${WORKDIR}/go-path/bin/ping-pong

include $(shell go list -f '{{ .Dir }}' -m go.gazette.dev/core)/mk/common-build.mk

# Push the ping-pong image to a specified private registry.
# Override the registry to use by passing a "registry=" flag to make.
registry=localhost:32000
push-to-registry:
	docker tag ping-pong:latest $(registry)/ping-pong:latest
	docker push $(registry)/ping-pong:latest