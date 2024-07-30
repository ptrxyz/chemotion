#!/bin/bash
[[ -n ${STARTUPDELAY} ]] && sleep "${STARTUPDELAY}"
[[ -f /embed/lib/env ]] && source /embed/lib/env

function cleanup() {
    rm -rf "${PIDFILE}"
    local pids
    mapfile -t pids < <(jobs -rp)
    if [[ ${#pids[@]} -gt 0 ]]; then
        kill "${pids[@]}" &>/dev/null
    fi
}

trap cleanup EXIT
trap cleanup SIGTERM
trap cleanup SIGINT

cd /chemotion/app || exit 1
[[ -n ${DROP_UID} && -n ${DROP_GID} ]] && {
    DROP="/embed/bin/drop"
}

if [[ "${CONFIG_ROLE}" == "eln" ]]; then
    (cd /embed/ && make all-eln)
    [[ -n ${DROP} ]] && export HOME=/chemotion/app
    exec ${DROP} bundle exec rails s -b 0.0.0.0 -p4000 --pid "${PIDFILE}"
    # exec passenger start -b 0.0.0.0 --pid-file "${PIDFILE}" --port 4000 --max-pool-size 5
elif [[ "${CONFIG_ROLE}" == "worker" ]]; then
    # Wait a bit. give the ELN some time to delete it's lock in case it's still present
    sleep 3
    (cd /embed/ && make all-worker)
    [[ -n ${DROP} ]] && export HOME=/chemotion/app
    exec ${DROP} bundle exec bin/delayed_job ${DELAYED_JOB_ARGS} run
else
    echo "ERROR: Please specify CONFIG_ROLE ('eln'/'worker')."
fi
