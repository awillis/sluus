.RECIPEPREFIX = >
GO:=$(shell which go)
VERSION:=0.1.0
PLUGINLIST:=$(shell find processor -name plugin -type d | cut -f2 -d/)
LDFLAGS:=-ldflags "-s -w -X github.com/awillis/sluus/core.VERSION=${VERSION}"

all:
> mkdir -p build/plugins build/bin
> ${GO} build ${LDFLAGS} -buildmode=pie -o build/bin/sluus
> $(foreach plug,$(PLUGINLIST), ${GO} build -buildmode=plugin -o build/plugins/$(plug).so ${PWD}/processor/$(plug)/plugin; echo;)

clean:
> rm -rvf build
> go clean