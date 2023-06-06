#!/bin/bash


if [[ -n "${DROP_UID}" && -n "${DROP_GID}" ]]; then
    log "Setting permissions to [${DROP_UID}:${DROP_GID}]. This might take a while ..."
    useradd -u "${DROP_UID}" -d /chemotion/app -r chemotion
    [[ "${CONFIG_ROLE}" == "eln" ]] && chown -R "${DROP_UID}:${DROP_GID}" /chemotion
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/log
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/tmp
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/uploads
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/public
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/config
    # chown -R ${DROP_UID}:${DROP_GID} /chemotion/app/app/packs/src/components/
    chown -R "${DROP_UID}:${DROP_GID}" /asdf /cache
    log "Permissions set."
fi
