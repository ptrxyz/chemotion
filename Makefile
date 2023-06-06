MAKEFLAGS += --always-make
MAKEFLAGS += -s
MAKEFLAGS += --no-print-directory

SHELL := /bin/bash

.PHONY: *

default: all

_%:
	localName="$*"; \
	localArgs="${$*_args}"; \
	localOpts="${$*_opts}"; \
	localPath="${$*_path}"; \
	localImageTag="${$*_imagetag}"; \
	localAdditionalTags="${$*_additionaltags}"; \
	localHash="${$*_hash}"; \
	localDockerfile="${$*_dockerfile}"; \
	localVersion="${$*_version}"; \
	buildContainer

base:
	$(MAKE) _$@

eln converter ketcher ketchersc spectra msconvert: base
	$(MAKE) _$@

all: base eln converter ketcher ketchersc spectra msconvert

devify:
	# this is experimental!
	@echo "Devifying [${REPO}/eln:latest] as [${REPO}/eln:devified]."
	cd dev; docker build --build-arg BUILDBASE="${REPO}/eln:latest" --tag "${REPO}/eln:devified" -f Dockerfile .
