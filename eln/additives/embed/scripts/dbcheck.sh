#!/bin/bash

source /embed/lib/dbenv

LOGFILE=${LOGFILE-/dbcheck.log}
export LOGFILE
cd /chemotion/app || exit 200

test1() {
    PGDATABASE='' psql -tA -c 'select current_user' &>>"${LOGFILE}"
    return $?
}

test2() {
    psql -tA -c 'select current_database()' &>>"${LOGFILE}"
    return $?
}

test3() {
    psql -tA -c "select 'non-empty' from molecule_names limit 1;"
    return $?
}

test4() {
    psql -tA -c "select 'non-empty' from ketcherails_common_templates limit 1;"
    return $?
}


test1 || exit 1
log "Authentication successful"

test2 || {
    log "DB will be created."
    bundle exec dotenv rake db:create 1>>"${LOGFILE}" || exit 2
}
log "Database exists"

bundle exec dotenv rake db:migrate 1>>"${LOGFILE}" || exit 3
log "Migrations completed."

out=$(test3) || {
    log "Can not reliably detect if the database is fresh or not. Seeding skipped."
    exit 4
}

if [[ "$out" =~ "non-empty" ]]; then
    log "Seeding not needed."
else
    log "Needs seeding...(this will take a while)"
    bundle exec rake db:seed 1>>"${LOGFILE}" || exit 5
    log "Seeding done."
fi

# add ketcherails templates if needed
out=$(test4) || {
    log "Can not reliably detect if common templates for Ketcher are needed. Skipping."
    log "Recreation triggered..."
    bundle exec rails r 'MakeKetcherailsSprites.perform_now'
    log "Spritesheet created."
    exit 5
}

if [[ "$out" =~ "non-empty" ]]; then
    log "No need to seed common templates for Ketcher."
else
    log "Seeding common template for Ketcher...(this will take a while)"
    bundle exec rake ketcherails:import:common_templates
    bundle exec rails r 'MakeKetcherailsSprites.perform_now'
fi

# Recreate Ketcher sprite sheet if it does not exist.
if [[ ! -f ./public/images/sprites/ketcherails.png ]]; then
    log "Spritesheet for Ketcher is missing. Recreation triggered..."
    bundle exec rails r 'MakeKetcherailsSprites.perform_now'
    log "Spritesheet created."
fi

