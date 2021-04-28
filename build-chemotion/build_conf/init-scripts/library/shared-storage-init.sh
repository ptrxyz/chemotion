#!/bin/bash
# RUNAS: user

[[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

# make sure we do this before DB config...
info "Setting up shared folders..."
cd /shared
mkdir -p /shared/uploads /shared/tmp_uploads
cp -R ${APP_DIR}/seed/* /shared || {
    error "Could not copy to shared storage. Make sure the permissions are correct."
    msg "The storage needs to be read and writeable by $(id)"
    exit 1
}
cp -R ${APP_DIR}/seed/.env /shared  || {
    error "Could not copy to shared storage. Make sure the permissions are correct."
    msg "The storage needs to be read and writeable by $(id)"
    exit 1
}

setVersion