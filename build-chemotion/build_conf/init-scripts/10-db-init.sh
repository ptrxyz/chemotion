#!/bin/bash
# RUNAS: root
# Database init

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

whoami
info "Creating database..."

if ! (echo "\q" | psql -d "${DB_NAME}" -h "${DB_HOST}" -U "${DB_ROLE}" 2>/dev/null) || [[ x"${INITIALIZE}" == x"yes" ]]; then
    info "Can not connect to database or database needs to be initialized."
    if [[ "x${CREATE_USER}" == x"yes" || x"${INITIALIZE}" == x"yes" ]]; then
        psql --host="${DB_HOST}" --username 'postgres' -c "
            DROP DATABASE IF EXISTS ${DB_NAME};"
        psql --host="${DB_HOST}" --username 'postgres' -c "
            DROP ROLE IF EXISTS ${DB_ROLE};
            CREATE ROLE ${DB_ROLE} LOGIN CREATEDB NOSUPERUSER PASSWORD '${DB_PW}';"
        psql --host="${DB_HOST}" --username 'postgres' -c "            
            CREATE DATABASE ${DB_NAME} OWNER ${DB_ROLE};
        " || {
            error "Could not create database. PSQL returned [$?]."
            exit 1
        }
        psql --host="${DB_HOST}" --username="${DB_ROLE}" -c "
            CREATE EXTENSION IF NOT EXISTS \"pg_trgm\";
            CREATE EXTENSION IF NOT EXISTS \"hstore\";
            CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
            ALTER USER ${DB_ROLE} PASSWORD '${DB_PW}';
        " || {
            error "Failed to set password for database user. PSQL returned [$?]."
            exit 1
        }
    else
        error "Could not connect to database. Make sure to specify connection parameters using DB_HOST, DB_ROLE, DB_NAME, DB_PW."
        exit 1
    fi
fi

# At this point we can be sure that the following command succeeds:
(echo "\q" | psql -d "${DB_NAME}" -h "${DB_HOST}" -U "${DB_ROLE}") || exit $?
info "Connection to database succeeded."

# (re-)configure the app to match our database variables
sed -i -e '/^DB_[A-Z]\+\=/d' -e '/^\s*$/d' /chemotion/app/.env
cat >>/chemotion/app/.env <<EOF
DB_NAME='${DB_NAME}'
DB_ROLE='${DB_ROLE}'
DB_PW='${DB_PW}'
DB_HOST='${DB_HOST}'
DB_PORT=${DB_PORT}
EOF
