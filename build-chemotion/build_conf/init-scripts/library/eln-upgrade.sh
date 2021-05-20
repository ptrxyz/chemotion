#!/bin/bash
# RUNAS: user

# This script should be executed to upgrade the container config and the database
# to a new version.
# Steps to take:
#   - run rake db:migrate (as $PROD)
#   - (TODO!) rewrite the config files

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

# source profile to get proper ruby environment
. ${HOME}/.profile

cd ${APP_DIR}/app
info "Executing migrations..."
bundle exec rake db:migrate >/dev/null
info "Database migrated."

info "Creating sprites..."
bundle exec rake ketcherails:import:common_templates
rm -rf ${APP_DIR}/app/public/images/ketcherails/icons/original/*
bundle exec rails r 'MakeKetcherailsSprites.perform_now'
info "Sprite generation done."

setVersion