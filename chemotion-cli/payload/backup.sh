#!/bin/bash

[[ -d /backup ]] || {
    log "Backup skipped. [/backup] is not mounted."
    exit 1
}

source /embed/lib/dbenv

stamp=$(date +'%y%m%d-%H%M%S')

if [[ -z $1 || $1 == "data" ]]; then
    tar cvzH posix -f "/backup/backup-${stamp}.data.tar.gz" --directory=/chemotion/data/ . --directory=/ .version || {
        log "Could not backup user data!"
        exit 2
    }
fi

if [[ -z $1 || $1 == "db" ]]; then
    pg_dump --no-owner --clean --if-exists | gzip -c > "/backup/backup-${stamp}.sql.gz" || {
        log "Could not backup database!"
        exit 3
    }
fi

if [[ ! -e "/backup/backup.data.tar.gz" && ! -e "/backup/backup.sql.gz" ]] || \
   [[ -L "/backup/backup.data.tar.gz" && -L "/backup/backup.sql.gz" ]] ; then 
    log "Creating symlink to latest backup."
    if [[ -z $1 || $1 == "data" ]]; then
        rm -f /backup/backup.data.tar.gz
        ln -s "backup-${stamp}.data.tar.gz" "/backup/backup.data.tar.gz"
    fi
    if [[ -z $1 || $1 == "db" ]]; then
        rm -f /backup/backup.sql.gz
        ln -s "backup-${stamp}.sql.gz" "/backup/backup.sql.gz"
    fi
fi

log "Backup finished successfully. Timestamp: ${stamp}"
