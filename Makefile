.RECIPEPREFIX = >
GO:=$(shell which go)
PLUGINLIST:=$(shell find processor -name plugin -type d | cut -f2 -d/)

all:
> mkdir -p build/plugins build/bin
> ${GO} build -o build/bin/sluus
> $(foreach plug,$(PLUGINLIST), ${GO} build -buildmode=plugin -o build/plugins/$(plug).so ${PWD}/processor/$(plug)/plugin; echo;)

clean:
> rm -rvf build