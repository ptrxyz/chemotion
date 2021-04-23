#!/bin/bash

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

. ${HOME}/.profile;

info "Used versions:"
ruby --version | head
passenger --version | head
echo "Node $(node --version)" | head
echo "NPM $(npm -v)" | head

cd ${APP_DIR}/app
if [[ ${ELN_ROLE} == "app" ]]; then 
    # start ketcher background service if present
    if [ -f "${APP_DIR}/app/lib/node_service/nodeService.js" ]; then
        nohup node ${APP_DIR}/app/lib/node_service/nodeService.js production >>$directory/log/node.log 2>&1 &
    fi

    exec passenger start -e ${RAILS_ENV} --engine=builtin --address 0.0.0.0 --port ${PASSENGER_PORT}
else
    exec bundle exec bin/delayed_job run
fi