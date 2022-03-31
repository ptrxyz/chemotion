#!/bin/bash
secret=$(erb /chemotion/app/config/secrets.yml 2>>$LOGFILE | yq '.production.secret_key_base')
[[ -n "$secret" ]] && {    
    log "Secret set to [${secret:0:10}]"
    exit 0
}

log "Generating new secret."

secret=$(cd /chemotion/app && rake secret)
yq -i '.production.secret_key_base="'$secret'"' /chemotion/app/config/secrets.yml
log "Secret set to [${secret:0:10}]"
