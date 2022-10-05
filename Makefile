SHELL := ./.envshell
.DEFAULT_GOAL := build

include .env
export

ALL_STAGES=gather eln converter ketchersvc spectra ruby node msconvert base

COMPOSE_PROJECT=chemotion-build
DOCKERCMD=COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT) DOCKER_BUILDKIT=0 docker
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
	$(DOCKERCMD) image rm -f chemotion-build:ketchersvc	|| true
	$(DOCKERCMD) image rm -f chemotion-build:converter	|| true
	$(DOCKERCMD) image rm -f chemotion-build:eln    	|| true
	$(DOCKERCMD) image rm -f chemotion-build:gather 	|| true
	make -C eln clean
	make -C spectra clean

tag:
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion:eln-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) tag chemotion-build:converter ptrxyz/chemotion:converter-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) tag chemotion-build:ketchersvc ptrxyz/chemotion:ketchersvc-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion:spectra-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion:msconvert-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion:eln-latest
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion:spectra-latest
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion:msconvert-latest
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion:eln-$(CHEMOTION_SHORT_RELEASE)
	$(DOCKERCMD) tag chemotion-build:converter ptrxyz/chemotion:converter-$(CHEMOTION_SHORT_RELEASE)
	$(DOCKERCMD) tag chemotion-build:ketchersvc ptrxyz/chemotion:ketchersvc-$(CHEMOTION_SHORT_RELEASE)
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion:spectra-$(CHEMOTION_SHORT_RELEASE)
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion:msconvert-$(CHEMOTION_SHORT_RELEASE)
	make -C testenv composefile

upload: tag
	# $(DOCKERCMD) push ptrxyz/chemotion:eln-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:converter-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:ketchersvc-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:spectra-latest
	# $(DOCKERCMD) push ptrxyz/chemotion:msconvert-latest
	$(DOCKERCMD) push ptrxyz/chemotion:eln-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:converter-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:ketchersvc-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:spectra-$(CHEMOTION_RELEASE)
	$(DOCKERCMD) push ptrxyz/chemotion:msconvert-$(CHEMOTION_RELEASE)
	# $(DOCKERCMD) push ptrxyz/chemotion:eln-$(CHEMOTION_SHORT_RELEASE)
	# $(DOCKERCMD) push ptrxyz/chemotion:converter-$(CHEMOTION_SHORT_RELEASE)
	# $(DOCKERCMD) push ptrxyz/chemotion:ketchersvc-$(CHEMOTION_SHORT_RELEASE)
	# $(DOCKERCMD) push ptrxyz/chemotion:spectra-$(CHEMOTION_SHORT_RELEASE)
	# $(DOCKERCMD) push ptrxyz/chemotion:msconvert-$(CHEMOTION_SHORT_RELEASE)

upload-dev:
	$(DOCKERCMD) tag chemotion-build:eln ptrxyz/chemotion-build:eln
	$(DOCKERCMD) tag chemotion-build:converter ptrxyz/chemotion-build:converter
	$(DOCKERCMD) tag chemotion-build:ketchersvc ptrxyz/chemotion-build:ketchersvc
	$(DOCKERCMD) tag chemotion-build:spectra ptrxyz/chemotion-build:spectra
	$(DOCKERCMD) tag chemotion-build:msconvert ptrxyz/chemotion-build:msconvert
	$(DOCKERCMD) push ptrxyz/chemotion-build:eln
	$(DOCKERCMD) push ptrxyz/chemotion-build:converter
	$(DOCKERCMD) push ptrxyz/chemotion-build:ketchersvc
	$(DOCKERCMD) push ptrxyz/chemotion-build:spectra
	$(DOCKERCMD) push ptrxyz/chemotion-build:msconvert

$(ALL_STAGES): prepare
	$(DOCKERCMD) build -f .dockerfile 						\
		--build-arg BASE=base										\
		--build-arg TINI_VERSION=$${TINI_VERSION} 					\
		--build-arg SPECTRA_VERSION=$${SPECTRA_VERSION} 			\
		--build-arg RUBY_VERSION=$${RUBY_VERSION} 					\
		--build-arg NODE_VERSION=$${NODE_VERSION} 					\
		--build-arg SPECTRA_BUILD_TAG=$${SPECTRA_BUILD_TAG} 		\
		--build-arg KETCHERSVC_BUILD_HASH=$${KETCHERSVC_BUILD_HASH} \
		--build-arg CONVERTER_BUILD_TAG=$${CONVERTER_BUILD_TAG} 	\
		--build-arg ELNDIR=./eln 									\
		--build-arg NODEDIR=./eln 									\
		--build-arg RUBYDIR=./eln 									\
		--build-arg CONVERTDIR=./spectra 							\
		--build-arg SPECTRADIR=./spectra							\
		--target $@													\
		-t chemotion-build:$@										\
		-t chemotion-build/$@:$${CHEMOTION_SHORT_RELEASE}			\
		.

rmgather:
	$(DOCKERCMD) image rm -f chemotion-build:gather || true

composefile:
	mkdir -p release/$(CHEMOTION_RELEASE)
	make -C testenv composefile
	cp testenv/docker-compose-$(CHEMOTION_RELEASE).yml release/$(CHEMOTION_RELEASE)/docker-compose.yml
	rm release/latest; ln -s $(CHEMOTION_RELEASE) release/latest

release: clean-build composefile
soft-release: clean build composefile
build-nocache: clean-cache $(ALL_STAGES)
clean-build: clean build-nocache
build: $(ALL_STAGES)

# vim: set tabstop=4:softtabstop=4:shiftwidth=4:noexpandtab
