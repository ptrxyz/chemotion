.DEFAULT_GOAL := build

COMPOSE_PROJECT=chemotion-build
DOCKERCMD=COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT) DOCKER_BUILDKIT=1 docker
COMPOSECMD=$(DOCKERCMD)-compose -f docker-compose-test.yml

composefile:
	 sed 's/RELEASE/$(CHEMOTION_BUILD_RELEASE)/g' docker-compose-production.yml > docker-compose-$(CHEMOTION_BUILD_RELEASE).yml

clean-up: clean
	$(COMPOSECMD) up --renew-anon-volumes --remove-orphans --force-recreate
	$(COMPOSECMD) down -v --remove-orphans

up:
	$(COMPOSECMD) up

down:
	$(COMPOSECMD) down

clean: down clean-volumes clean-files

clean-volumes: 
	$(DOCKERCMD) volume rm -f chemotion $(COMPOSE_PROJECT)_spectra
	$(DOCKERCMD) volume rm -f chemotion $(COMPOSE_PROJECT)_chemotion
	$(DOCKERCMD) volume rm -f chemotion $(COMPOSE_PROJECT)_chemotion_data

clean-files:
	sudo rm -rf /tmp/db-data
	sudo rm -rf ./shared

shell: waitalive
	$(COMPOSECMD) exec eln /bin/bash

backup: waitalive
	$(COMPOSECMD) exec eln chemotion backup

restore:
	$(COMPOSECMD) run eln chemotion restore

waitalive:
	until $(COMPOSECMD) ps | egrep -e 'eln\s+running'; do sleep 1; done

