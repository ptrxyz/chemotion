#!/bin/bash

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

[[ -f ${INIT_BASE}/init.env ]] && source ${INIT_BASE}/init.env

while ! pg_isready -h "${DB_HOST}" 1>/dev/null 2>&1; do
    msg "Database not ready. Waiting ..."
    sleep 10
done

if [[ x"${INITIALIZE}" == x"yes" ]]; then
    [[ -d ${INIT_BASE}/init-scripts/ ]] || {
        error "Can not find init-directory."
        exit 1
    }
    declare -a init_files=($(ls ${INIT_BASE}/init-scripts/* | sort))
    for i in "${init_files[@]}"; do
        head -n2 "${i}" | grep 'RUNAS: root' 1>/dev/null && needs_root=true || needs_root=false
        if $needs_root; then
            note "--> running [$(basename ${i})] as [root]!"
            bash "${i}"
        else
            note "--> running [$(basename ${i})] as [${PROD}]!"
            sudo -E -H -u ${PROD} bash "${i}"
        fi
    done
    exit 0
fi

(echo "\q" | psql -d "${DB_NAME}" -h "${DB_HOST}" -U "${DB_ROLE}") || {
    error "Could not connect to database. Make sure to initialize it!"
    exit $?
}
info "Connection to database succeeded."

# make sure to use single quotes in the next line to prevent $HOME being resolved to /root by the invoking shell!
[[ -f "${INIT_BASE}/run-${ELN_ROLE}.sh" ]] && exec sudo -E -H -u ${PROD} ${INIT_BASE}/run-${ELN_ROLE}.sh
