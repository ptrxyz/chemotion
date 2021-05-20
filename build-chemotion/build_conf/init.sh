#!/bin/bash

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

waitForDB

if [[ x"${INITIALIZE}" == x"yes" ]]; then
    [[ -d ${INIT_BASE}/init-scripts/enabled/ ]] || {
        error "Can not find init-directory."
        exit 1
    }
    declare -a init_files=($(ls ${INIT_BASE}/init-scripts/enabled/* | sort))
    for i in "${init_files[@]}"; do
        execute "${i}"
    done
    exit 0
fi

DBConnect || exit $?

[[ -f "${INIT_BASE}/run-${ELN_ROLE}.sh" ]] && execute ${INIT_BASE}/run-${ELN_ROLE}.sh
