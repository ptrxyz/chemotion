#!/bin/bash

source /embed/lib/dbenv

until pg_isready &>>"${LOGFILE}"; do
    log "Waiting for database [${PGHOST}:${PGPORT}]..."
    sleep 3;
done
log "$(pg_isready)"
