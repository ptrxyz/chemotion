#!/bin/bash

export CC_RED='\033[0;31m'
export CC_GREEN='\033[0;32m'
export CC_YELLOW='\033[0;33m'
export CC_CYAN='\033[0;36m'
export CC_NC='\033[0m'

[[ -f ${INIT_BASE}/overwrite.env ]] && source ${INIT_BASE}/overwrite.env

error() {
    echo -ne "${CC_RED}"
    echo $@
    echo -ne "${CC_NC}"
}

warn() {
    echo -ne "${CC_YELLOW}"
    echo $@
    echo -ne "${CC_NC}"
}

info() {
    echo -ne "${CC_CYAN}"
    echo $@
    echo -ne "${CC_NC}"
}

msg() {
    echo $@
}

ensureRoot() {
    # Returns 0 if we are root, 1 otherwise

    if [[ "$EUID" -eq 0 ]]; then
        return 0
    else        
        echo -e "${CC_RED}Please make sure to run this as root.${CC_NC}"
        return 1
    fi
}

odus() {
    # 'Opposite' of sudo ...: run a single command as user ${PROD}
    # (uses a shell, sources profile)

    sudo -E -H -u ${PROD} bash -c '. $HOME/.profile; '"$@"
}

execute() {
    # executes all files passed as parameters
    # looks for "RUNAS: root|user" in the first two lines and sets the 
    # execution context accordingly

    ff="$1"
    while [[ ! -z ${ff} ]]; do
        if [[ ! -f "${ff}" ]]; then
            echo -e "${CC_RED}File not found:${CC_NC} ${ff}"
        else
            head -n2 "${ff}" | grep 'RUNAS: root' 1>/dev/null && needs_root=true || needs_root=false
            if $needs_root; then
                echo -e "${CC_CYAN}Executing [${CC_NC}$(basename ${ff})${CC_CYAN}] as [${CC_NC}root${CC_CYAN}].${CC_NC}"
                bash "${ff}"
            else
                echo -e "${CC_CYAN}Executing [${CC_NC}$(basename ${ff})${CC_CYAN}] as [${CC_NC}${PROD}${CC_CYAN}].${CC_NC}"
                sudo -E -H -u ${PROD} bash "${ff}"
            fi            
        fi
        shift
        ff="$1"
    done
}

waitForDB() {
    # simply waits for the DB to be up and done booting

    while ! pg_isready -h "${DB_HOST}" 1>/dev/null 2>&1; do
        msg "Database not ready. Waiting ..."
        sleep 10
    done
}

setVersion() {    
    [[ -d "$(dirname ${VERSION_FILE})" ]] || {
        warn "Could not write version file!"
        msg "Make sure [$VERSION_FILE] is writeable for UID $(id -u ${PROD})"
    }
    echo ${CHEMOTION_VERSION}-${FLAVOR} > ${VERSION_FILE}

    [[ "${EUID}" -eq 0 ]] && {
        chown ${PROD}:${PROD} ${VERSION_FILE}
        chmod a+wr ${VERSION_FILE}
    }
}

DBConnect() {
    # This makes sure that we can connect as the 'proper' DB user to
    # the DB used by the actual app

    (echo "\q" | psql -d "${DB_NAME}" -h "${DB_HOST}" -U "${DB_ROLE}") || {
        error "Could not connect to database. Make sure to initialize it!"
        return $?
    }
    info "Connection to database succeeded."
    return 0
}

versionMatching() {
    # Check if we initalized with this container version
    thisVersion="$(${CHEMOTION_VERSION}-${FLAVOR})"
    theirVersion="$(cat ${VERSION_FILE} 2>/dev/null || echo 'unknown')"

    if [[ ! -f "${VERSION_FILE}" ]]; then
        error "Can not find container version info."
        msg "Please make sure the container is properly initalized."
        msg "Expected version is: $thisVersion"        
        return 1
    fi

    if [[ ${thisVersion} != ${theirVersion} ]]; then
        error "Container version mismatch."
        msg "This container version is: ${thisVersion}"
        msg "Your version is: ${theirVersion}"
        return 2
    fi

    return 0
}
