#!/bin/bash
LOGFILE=${LOGFILE-/secret.log}
export LOGFILE
cd /chemotion/app || exit 200
touch /chemotion/app/.env

readSecret() {
    bundle exec dotenv erb /chemotion/app/config/secrets.yml | yq '.production.secret_key_base'
}

secret=$(readSecret 2>>"${LOGFILE}")
if [[ -n "$secret" ]] && [[ "${secret}x" != "nullx" ]]; then
    echo "Secret already present in configuration."
else
    log "Generating new secret."
    secret=$(head -c32 < /dev/urandom | base64)
    log "Generated secret is: [${secret:0:10}]."
    yq -i '.production.secret_key_base="'"$secret"'"' /chemotion/app/config/secrets.yml
fi

secret=$(readSecret 2>>"${LOGFILE}")
if [[ -z "$secret" ]] || [[ "${secret}x" == "nullx" ]]; then
    log "Secret was not set correctly. Please make sure it is configured in [config/secrets.yml]."
    exit 1
fi
log "Secret set to [${secret:0:10}]"
