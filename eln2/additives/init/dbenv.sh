#!/bin/bash
firstOf() {
    # takes a list of arguments and returns the first non-empty one
    while [[ -z "$1" && -n "$@" ]]; do shift; done
    echo $1
}

DBCONFFILE="/chemotion/app/config/database.yml"
ENVIRONMENT=$(firstOf "${RAILS_ENV}" "development")

getenv() {
    # $1: env variable name
    # $2: yml key
    # $3: default value

    # return env variable if its set.
    [[ -n "${!1}" ]] && \
        echo ${!1} && return 0;

    # check yml file and return if it has a value != ''
    yml=$(yq '.'${ENVIRONMENT}'.'$2 $DBCONFFILE)
    [[ $? -eq 0 && "$yml" != "null" ]] && \
        echo $yml && return 0;

    # return default value
    echo "$3"
}

export PGHOST=$(getenv DB_HOST host db)
export PGPORT=$(getenv DB_PORT port 5432)
export PGUSER=$(getenv DB_USERNAME username chemotion)
export PGPASSWORD=$(getenv DB_PASSWORD password chemotion)
export PGDATABASE=$(getenv DB_DATABASE database chemotion)
