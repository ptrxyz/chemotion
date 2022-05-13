#!/bin/bash

[[ ! -f /backup/backup.data.tar.gz ]] && {
    log "Can not restore. [/backup/backup.data.tar.gz] is not found."
    exit 1
}

[[ ! -f /backup/backup.sql.gz ]] && {
    log "Can not restore. [/backup/backup.data.tar.gz] is not found."
    exit 1
}

source /embed/lib/dbenv

test1() {
    PGDATABASE='' psql -tA -c 'select current_user' &>>${LOGFILE}
    return $?
}

test1 || {
    log "Database authentication failed. Restoration chancelled."
    exit 1
}


gunzip < "/backup/backup.sql.gz" | psql
RETSQL=$?
cat /backup/backup.data.tar.gz | tar xvz --exclude=".version" --directory /chemotion/data/
RETTAR=$?

[[ $RETSQL -eq 0 ]] && log "Database restoration finished successfully."
[[ $RETTAR -eq 0 ]] && log "File restoration finished successfully."

[[ $RETSQL -ne 0 || $RETTAR -ne 0 ]] && {
    log "Something went wrong. Your data is NOT restored successfully."
    log "Please check the output of what went wrong."
}