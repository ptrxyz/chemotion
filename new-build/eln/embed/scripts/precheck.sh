#!/bin/bash

# Prepare log file
mkdir -p "$(dirname ${LOGFILE})"
truncate -s0 ${LOGFILE}

# Remove any pidfiles that might exist
log "PIDfile set to [$PIDFILE]"
[[ -e "${PIDFILE}" ]] && {
    log "Cleaning up old PID."
    rm -rf "${PIDFILE}"
}

# check version in volume vs version in container
this="$(grep '^RELEASE=' /.version | cut -d'=' -f2)"
theirs="$(grep '^RELEASE=' /chemotion/app/.version | cut -d'=' -f2)"
[[ "$this" == "$theirs" ]] && {
    log "Versions match."
} || {
    log "WARNING: Version mismatch [${this} (container)] vs. [${theirs} (runtime volume)]."
    log "Consider recreating the volume."
}

[[ ! -e /shared/.version ]] && cp /.version /shared/.version
[[ ! -e /chemotion/data/.version ]] && cp /.version /chemotion/data/.version

# check that uploads is mounted
(mount | grep -e 'on /chemotion/data/uploads type' -e 'on /chemotion/data type' &>/dev/null) || \
    log -e "\x1B[31mWARNING:\x1B[0m uploads folder is not mounted. All user data will be stored in the Chemotion volume."
