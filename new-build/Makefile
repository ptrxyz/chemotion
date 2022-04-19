SHELL := ./.envshell
.DEFAULT_GOAL := build

include .env
export

ALL_STAGES=gather eln spectra ruby node msconvert base

COMPOSE_PROJECT=chemotion-build
DOCKERCMD=COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT) DOCKER_BUILDKIT=1 docker
COMPOSECMD=$(DOCKERCMD)-compose -f docker-compose-test.yml

prepare: 
	(cat dockerfiles/??-dckr-* dockerfiles/99-gather > .dockerfile)
	make -C eln prepare
	make -C spectra prepare

clean-cache:
	$(DOCKERCMD) buildx  prune -a -f --verbose || true
	$(DOCKERCMD) builder prune -a -f --filter type=exec.cachemount || true

clean:
	rm -f .dockerfile
	$(DOCKERCMD) image rm -f chemotion-build:base		|| true
	$(DOCKERCMD) image rm -f chemotion-build:node		|| true
	$(DOCKERCMD) image rm -f chemotion-build:ruby		|| true
	$(DOCKERCMD) image rm -f chemotion-build:msconvert	|| true
	$(DOCKERCMD) image rm -f chemotion-build:spectra	|| true
	$(DOCKERCMD) image rm -f chemotion-build:eln    	|| true
	$(DOCKERCMD) image rm -f chemotion-build:gather 	|| true

tag:
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion:eln-$(CHEMOTION_BUILD_RELEASE)
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion:spectra-$(CHEMOTION_BUILD_RELEASE)
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion:msconvert-$(CHEMOTION_BUILD_RELEASE)
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion:eln-latest
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion:spectra-latest
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion:msconvert-latest
	make -C testenv composefile

upload: tag
	# $(DOCKERCMD) push ptrxyz/chemotion:eln-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:spectra-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:msconvert-latest
	$(DOCKERCMD) push ptrxyz/chemotion:eln-$(CHEMOTION_BUILD_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:spectra-$(CHEMOTION_BUILD_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:msconvert-$(CHEMOTION_BUILD_RELEASE)
	

upload-dev:
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion-build:eln
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion-build:spectra
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion-build:msconvert
	$(DOCKERCMD) push ptrxyz/chemotion-build:eln
	$(DOCKERCMD) push ptrxyz/chemotion-build:spectra
	$(DOCKERCMD) push ptrxyz/chemotion-build:msconvert

$(ALL_STAGES): prepare
	$(DOCKERCMD) buildx build -f .dockerfile \
		--build-arg TINI_VERSION=$${TINI_VERSION} \
		--build-arg ELNDIR=./eln 			\
		--build-arg NODEDIR=./eln 			\
		--build-arg RUBYDIR=./eln 			\
		--build-arg CONVERTDIR=./spectra 	\
		--build-arg SPECTRADIR=./spectra	\
		--target $@							\
		-t chemotion-build:$@				\
		.
	@[[ "$@" == "gather" ]] && $(DOCKERCMD) image rm -f chemotion-build:gather || true

build-nocache: clean-cache $(ALL_STAGES)
clean-build: clean build-nocache
build: $(ALL_STAGES)
