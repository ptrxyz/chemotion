#!/bin/bash

[[ ! -f /backup/backup.data.tar.gz ]] && {
    log "Can not restore. [/backup/backup.data.tar.gz] is not found."
    exit 1
}

[[ ! -f /backup/backup.sql.gz ]] && {
    log "Can not restore. [/backup/backup.sql.gz] is not found."
    exit 1
}

source /embed/lib/dbenv

test1() {
    PGDATABASE='' psql -tA -c 'select current_user' &>>"${LOGFILE}"
    return $?
}

test1 || {
    log "Database authentication failed. Restoration chancelled."
    exit 1
}
log "DB authentication successful."

test2() {
    psql -tA -c 'select current_database()' &>>"${LOGFILE}"
    return $?
}

if [[ ${FORCE_DB_RESET} -eq 1 || "$1" == "--force-db-reset" ]]; then
    if test2; then
        DB_TO_CREATE=$(echo ${PGDATABASE} | tr -d '"')
        log "Database reset enforced. [$DB_TO_CREATE] will be dropped!"
        sleep 10;
        PGDATABASE='' psql -tA -c 'DROP DATABASE "'${DB_TO_CREATE}'";' 2>>"${LOGFILE}"
        if [[ ! $? -eq 0 ]]; then
            log "Could not drop database [$PGDATABASE]. Maybe permissions are not set properly?"
            log "Trying to empty database manually."
            psql -tA -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public; CREATE EXTENSION IF NOT EXISTS "pg_trgm"; CREATE EXTENSION IF NOT EXISTS "hstore"; CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE EXTENSION IF NOT EXISTS "plpgsql";' 1>>${LOGFILE} && echo "Sucessfully dropped schema."
        fi
    else
        log "Old database not present. Nothing to delete."
    fi
fi

test2 || {
    log "DB not found. Trying to create it."
    DB_TO_CREATE=$(echo ${PGDATABASE} | tr -d '"')
    PGDATABASE='' psql -tA -c 'CREATE DATABASE "'${DB_TO_CREATE}'";' 1>>"${LOGFILE}"
    # bundle exec dotenv rake db:create 1>>"${LOGFILE}" || exit 2
    test2 || {
        echo "Could not create DB. Restoration cancelled."
    }
}
log "Database ready."

gunzip < "/backup/backup.sql.gz" | psql 1>>"${LOGFILE}"
RETSQL=$?
tar xvz --exclude=".version" --directory /chemotion/data/ < /backup/backup.data.tar.gz 1>>"${LOGFILE}"
RETTAR=$?

[[ $RETSQL -eq 0 ]] && log "Database restoration finished successfully."
[[ $RETTAR -eq 0 ]] && log "File restoration finished successfully."

if [[ $RETSQL -ne 0 || $RETTAR -ne 0 ]]; then
    log "Something went wrong. Your data is NOT restored successfully."
    log "Please check the output of what went wrong."
    exit 9
fi

exit 0

