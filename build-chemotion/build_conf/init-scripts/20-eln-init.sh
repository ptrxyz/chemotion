#!/bin/bash
# RUNAS: user
# Base app init

# # Make sure we run this as non-root ($PROD)
# if [[ "$EUID" -eq 0 ]]; then
#     echo "Forking [$0] away as [${PROD}]."
#     sudo -E -H -u "${PROD}" bash "$0" && exit 0
#     exit $?
# fi

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

# source profile to get proper ruby environment
. ${HOME}/.profile

cd ${APP_DIR}/app
info "Initializing database schemas..."
bundle exec rake db:create
info "Database created."
bundle exec rake db:migrate >/dev/null
info "Database migrated."
bundle exec rake db:seed
info "Database seeded."

info "Creating sprites..."
bundle exec rake ketcherails:import:common_templates
rm -rf ${APP_DIR}/app/public/images/ketcherails/icons/original/*
bundle exec rails r 'MakeKetcherailsSprites.perform_now'
