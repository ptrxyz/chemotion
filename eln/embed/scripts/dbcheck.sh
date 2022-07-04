#!/bin/bash

source /embed/lib/dbenv

test1() {
    PGDATABASE='' psql -tA -c 'select current_user' &>>${LOGFILE}
    return $?
}

test2() {
    psql -tA -c 'select current_database()' &>>${LOGFILE}
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
    (cd /chemotion/app && bundle exec rake db:create 1>>${LOGFILE}) || exit 2
}
log "Database exists"

(cd /chemotion/app && bundle exec rake db:migrate 1>>${LOGFILE}) || exit 3
log "Migrations completed."

out=$(test3)
[[ $? -eq 0 ]] || {
    log "Can not reliably detect if the database is fresh or not. Seeding skipped."
    exit 4
}

[[ "$out" =~ "non-empty" ]] && {
    log "Seeding not needed."
} || {
    log "Needs seeding...(this will take a while)"
    (cd /chemotion/app && bundle exec rake db:seed 1>>${LOGFILE}) || exit 5
    log "Seeding done."
}

# add ketcherails templates if needed
out=$(test4)
[[ $? -eq 0 ]] || {
    log "Can not reliably detect if common templates for Ketcher are needed. Skipping."
    exit 5
}

[[ "$out" =~ "non-empty" ]] && {
    log "No need to seed common templates for Ketcher."
} || {
    log "Seeding common template for Ketcher...(this will take a while)"
    bundle exec rake ketcherails:import:common_templates
    bundle exec rails r 'MakeKetcherailsSprites.perform_now'
}
