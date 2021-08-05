#!/bin/bash
set -o xtrace

[[ -d /workspace/chemotion/rdkit_chem.tar.gz ]] && ln -s /workspace/chemotion/rdkit_chem.tar.gz /precompiled/

yarn install
bash dbinit.sh
bundle install

bundle exec rake db:create
bundle exec rake db:migrate
bundle exec rake db:seed
