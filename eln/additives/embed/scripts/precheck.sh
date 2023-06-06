#!/bin/bash

# Prepare log file
LOGDIR=$(dirname "${LOGFILE}")
mkdir -p "$LOGDIR"
truncate -s0 "${LOGFILE}"

# Prepare most important folders.
# THIS SHOULD NEVER BE NECESSARY OR FAIL.

mkdir -p /shared
mkdir -p /chemotion/app
mkdir -p /chemotion/data
touch /chemotion/app/.env

function volumePopulationDone() {
    count1=$(find /chemotion/ | wc -l)
    sleep 10
    count2=$(find /chemotion/ | wc -l)

    if [[ ${count1} -lt ${count2} ]]; then
        # volumes still populating ...
        echo "Volumes still populating...(currently ${count2} files)"
        return 1
    else
        # no new files appeared in the last 10 seconds.
        # we assume volumes are populated.
        echo "Volumes ready. Final file count: ${count2}"
        return 0
    fi
}

if [[ -n "${WAIT_FOR_VOLUME_POPULATION}" ]]; then
    echo "Waiting for volume population..."
    until volumePopulationDone; do
       sleep 3
    done
fi

# Remove any pidfiles that might exist
log "PIDfile set to [$PIDFILE]"
[[ -e "${PIDFILE}" ]] && {
    log "Cleaning up old PID."
    rm -rf "${PIDFILE}"
}

# check that uploads is mounted
if ! mount | grep -e 'on /chemotion/data/uploads type' -e 'on /chemotion/data type' &>/dev/null; then
    log -e "\x1B[31mWARNING: the upload folder is not mounted.\x1B[0m\nAll uploaded user data will be stored in a temporary volume and will be lost when the container is stopped."
fi

# check version in volume vs version in container
[[ ! -e /shared/.version ]] && cp /.version /shared/.version
[[ ! -e /chemotion/data/.version ]] && cp /.version /chemotion/data/.version

container_version="$(grep '^RELEASE=' /.version | cut -d'=' -f2)"
app_version="$(grep '^RELEASE=' /chemotion/app/.version | cut -d'=' -f2)"
data_version="$(grep '^RELEASE=' /chemotion/app/.version | cut -d'=' -f2)"

chemotion_ref="$(grep '^CHEMOTION_REF=' /chemotion/app/.version | cut -d'=' -f2)"
chemotion_tag="$(grep '^CHEMOTION_TAG=' /chemotion/app/.version | cut -d'=' -f2)"

log "Container Version: ${container_version}"
log "App Version: ${app_version}"
log "Data Version: ${data_version}"
log "Chemotion Ref: ${chemotion_ref}"
log "Chemotion Tag: ${chemotion_tag}"

if [[ "$container_version" == "$app_version" ]]; then
    log "Container and App Versions match."
else
    log "WARNING: Version mismatch [${container_version} (container)] vs. [${app_version} (application)]."
    log "Please recreate the application volume."
    exit 1
fi

if [[ "$app_version" == "$data_version" ]]; then
    log "App and Data Versions match."
else
    log "WARNING: Version mismatch [${app_version} (application)] vs. [${data_version} (user data)]."
    log "Data volume will be upgraded."
fi

# check /tmp
if ! ( touch /tmp/dummy && rm /tmp/dummy ); then
    log "[/tmp] does not seem to be writable! This will cause problems!"
fi
