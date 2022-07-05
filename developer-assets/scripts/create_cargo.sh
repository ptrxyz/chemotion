#!/bin/bash

chown -R 0:0 /cargo
mkdir -p $GEM_HOME

# get to the specified commit
[ ! -d $CHEMOTION_DIR ] && git clone https://github.com/ComPlat/chemotion_ELN.git $CHEMOTION_DIR
cd $CHEMOTION_DIR
git stash
git fetch --all
git -c advice.detachedHead=false checkout $BRANCH

# make node modules
yarn install --production=true --modules-folder $NODE_MODULES --cache-folder $YARN_CACHE

# install ruby
gem install solargraph
bundle install --jobs=$(getconf _NPROCESSORS_ONLN)
if [[ $(grep -L passenger Gemfile) ]]; then bundle add passenger; fi

chown -R $UID:$GID /cargo