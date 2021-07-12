#!/bin/bash

export CC_RED='\033[0;31m'
export CC_GREEN='\033[0;32m'
export CC_YELLOW='\033[0;33m'
export CC_CYAN='\033[0;36m'
export CC_NC='\033[0m'

[[ -f /config/overwrite.env ]] && source /config/overwrite.env

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

function updateParam(){
    DB_ROLE="$1"
    DB_NAME="$2"
    DB_PW="$3"
    DB_HOST="$4"
    DB_PORT="$5"
    
    info "To use default values, just press Enter"
    
    read -p "(1 of 5) ::: Database role [default: chemotion]" DB_ROLLE
    DB_ROLLE=${DB_ROLLE:-chemotion}
   
    read -p "(2 of 5) ::: Database name [default: chemotion]" DB_NAME
    DB_NAME=${DB_NAME:-chemotion}
   
    read -p "(3 of 5) ::: database password [default: changeme]" DB_PW
    DB_PW=${DB_PW:-changeme}
    
    read -p "(4 of 5) ::: Database Host [default: db]" DB_HOST
    DB_HOST=${DB_HOST:-db}

    read -p "(5 of 5) ::: Database port [default:5432]" DB_PORT
    DB_PORT=${DB_PORT:-5432}

    info "setting new values for environment variables..."
    export DB_ROLE="$DB_ROLLE"
    export DB_NAME="$DB_NAME"
    export DB_PW="$DB_PW"
    export DB_HOST="$DB_HOST"
    export DB_PORT="$DB_PORT" 
    
    info "all environment variables has been reassigned with new values"
    echo "-----"
    echo "DB_ROLLE: "$DB_ROLLE
    echo "DB_NAME: "$DB_NAME
    echo "DB_PW: "$DB_PW
    echo "DB_HOST: "$DB_HOST
    echo "DB_PORT: "$DB_PORT
    echo "-----"

    export mustash=${APP_DIR}/app/mustash/mo
    export templates=${APP_DIR}/app/mustash/templates

    env | grep -e "^DB_NAME" -e "^DB_PW=" -e "^DB_HOST="  -e "^DB_ROLLE" -e "^DB_PORT" > /config/overwrite.env

    ${mustash} ${templates}/env_template.mo > ${APP_DIR}/app/.env
    ${mustash} ${templates}/env_template.mo > ${APP_DIR}/seed/.env 

    ${mustash} ${templates}/database_template.mo	> ${APP_DIR}/app/config/database.yml    
    ${mustash} ${templates}/datacollectors_template.mo	> ${APP_DIR}/app/config/datacollectors.yml
    #${mustash} ${templates}/editors_template.mo 	> ${APP_DIR}/app/config/editors.yml
    #${mustash} ${templates}/inference_template.mo 	> ${APP_DIR}/app/config/inference.yml
    ${mustash} ${templates}/secrets_template.mo 	> ${APP_DIR}/app/config/secrets.yml
    #${mustash} ${templates}/spectra_template.mo 	> ${APP_DIR}/app/config/spectra.yml
    ${mustash} ${templates}/storage_template.mo 	> ${APP_DIR}/app/config/storage.yml
    ${mustash} ${templates}/user_props_template.mo 	> ${APP_DIR}/app/config/user_props.yml
    #${APP_DIR}/seed/config/database.yml

    ${mustash} ${templates}/database_template.mo 	> ${APP_DIR}/seed/config/database.yml
    ${mustash} ${templates}/datacollectors_template.mo > ${APP_DIR}/seed/config/datacontrollrs.yml
    ${mustash} ${templates}/secrets_template.mo 	> ${APP_DIR}/seed/config/secrets.yml
    ${mustash} ${templates}/storage_template.mo 	> ${APP_DIR}/seed/config/storage.yml
    ${mustash} ${templates}/user_props_template.mo 	> ${APP_DIR}/seed/config/user_props.yml

    echo "-----"    
    echo "content of ${APP_DIR}/app/.env"
    cat ${APP_DIR}/app/.env
    echo "-----" 
    echo "content of ${APP_DIR}/seed/.env"
    cat ${APP_DIR}/seed/.env
    echo "-----" 
    echo "content of ${APP_DIR}/app/config/database.yml"
    cat ${APP_DIR}/app/config/database.yml
    echo "-----" 
    echo "content of ${APP_DIR}/seed/config/database.yml"
    cat ${APP_DIR}/seed/config/database.yml
    printf "\n" 
}

function confirm() {
    # ask for confirmation.
    # $1: text to show. Should indicate what's accepted 
    #     as confirmation (i.e. "Please enter 'yes' to continue")
    # $2: string considered as confirmation
    
    text="$1"
    confirm="$2"

    if [[ -z "$text" || -z "$confirm" ]]; then
        echo "Syntax error"
        return 1
    fi

    echo -e "$text"
    echo -n "> "
    read a

    if [[ x"$a" == x"$confirm" ]]; then
        return 0
    fi
    return 1
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

versionMatching() {
    # Check if we initalized with this container version
    thisVersion="${CHEMOTION_VERSION}-${FLAVOR}"
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
