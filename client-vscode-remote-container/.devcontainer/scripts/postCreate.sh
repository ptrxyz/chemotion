#!/bin/bash
set -o xtrace

yarn install

[[ -f .devcontainer/scripts/dbinit.sh ]] && (
	bash .devcontainer/scripts/dbinit.sh
)

# precompiledGems=/workspace/chemotion/.devcontainer/gems.tar.gz
# [[ -f "${precompiledGems}" ]] && (
# 	cd $HOME
# 	tar xfvz "${precompiledGems}" && rm "${precompiledGems}"
# )

bundle install

bundle exec rake db:create
bundle exec rake db:migrate
bundle exec rake db:seed
