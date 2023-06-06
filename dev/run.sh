#!/bin/bash
[[ -n ${STARTUPDELAY} ]] && sleep "${STARTUPDELAY}"
[[ -f /embed/lib/env ]] && source /embed/lib/env
if [[ -f /.env ]]; then
    set -a
    source /.env
    set +a
fi

function cleanup() {
    rm -rf "${PIDFILE}"
    local pids
    mapfile -t pids < <(jobs -rp)
    if [[ ${#pids[@]} -gt 0 ]]; then
        kill "${pids[@]}" &>/dev/null
    fi
}

trap cleanup EXIT SIGHUP SIGQUIT SIGTERM SIGINT HUP QUIT TERM INT

cd /chemotion/app || exit 1
[[ -n ${DROP_UID} && -n ${DROP_GID} ]] && {
    DROP="/embed/bin/drop"
}

if [[ "${CONFIG_ROLE}" == "eln" ]]; then
    (cd /embed/ && make all-eln) || exit 1
    [[ -n ${DROP} ]] && export HOME=/chemotion/app
    exec ${DROP} bundle exec rails s -b 0.0.0.0 -p4000 --pid "${PIDFILE}"
    # exec passenger start -b 0.0.0.0 --pid-file "${PIDFILE}" --port 4000 --max-pool-size 5
elif [[ "${CONFIG_ROLE}" == "worker" ]]; then
    # Wait a bit. give the ELN some time to delete it's lock in case it's still present
    sleep 3
    (cd /embed/ && make all-worker) || exit 1
    [[ -n ${DROP} ]] && export HOME=/chemotion/app
    exec ${DROP} bundle exec bin/delayed_job run
elif [[ "${CONFIG_ROLE}" == "dev" ]]; then
    (set -o pipefail; cd /embed/ && make all-dev 2>&1 | grep --line-buffered -v -e "warning: already initialized constant" -e "warning: previous definition of" -e "warning: URI.escape is obsolete") || exit 1
    [[ -n ${DROP} ]] && export HOME=/chemotion/app
    mkdir -p log/
    (sleep 10 && ${DROP} bundle exec bin/delayed_job run | tee -a log/worker.log) &
    ${DROP} bundle exec rails s -b 0.0.0.0 -p4000 --pid "${PIDFILE}" | tee -a log/server.log
else
    echo "ERROR: Please specify CONFIG_ROLE ('eln'/'worker')."
fi
