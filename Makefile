.RECIPEPREFIX = >
GO:=$(shell which go)
VERSION:=0.1.0
PLUGINLIST:=$(shell find plugin -mindepth 1 -name plugin -type d | cut -f2 -d/)
LDFLAGS:=-ldflags "-s -w -X github.com/awillis/sluus/core.VERSION=${VERSION}"
PKGLIST:=$(go list --deps | grep sluus)

build: protoc
> mkdir -p build/plugin build/bin
# add -buildmode=pie in at a later date
> ${GO} build ${LDFLAGS} -race -o build/bin/sluus
> $(foreach plug,$(PLUGINLIST), ${GO} build -race -buildmode=plugin -o build/plugin/$(plug).so ${PWD}/plugin/$(plug)/plugin;)
protoc:
> protoc -I protobufs -I message --go_out message message.proto

test:
> go test ${PKGLIST}

clean:
> rm -rvf build
> go clean