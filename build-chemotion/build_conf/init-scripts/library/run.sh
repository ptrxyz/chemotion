#!/bin/bash
# RUNAS: user

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

versionMatching || {
    error "Please upgrade."
    exit 100
}

. ${HOME}/.profile;

cd ${APP_DIR}/app

echo "Role: $ELN_ROLE"
if [[ ${ELN_ROLE} == "eln" || ${ELN_ROLE} == "app" ]]; then 
    # start ketcher background service if present
    if [ -f "${APP_DIR}/app/lib/node_service/nodeService.js" ]; then
        nohup node ${APP_DIR}/app/lib/node_service/nodeService.js production >>$directory/log/node.log 2>&1 &
    fi

    exec passenger start -e ${RAILS_ENV} --engine=builtin --address 0.0.0.0 --port ${PASSENGER_PORT}
else
    exec bundle exec bin/delayed_job run
fi