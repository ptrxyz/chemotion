#!/bin/bash

[[ ! -z "${INIT_BASE}" ]] && [[ -f ${INIT_BASE}/functions.sh ]] && source ${INIT_BASE}/functions.sh || {
    echo "Could not load base functions!"
    exit 1
}

containerInfo() {
    cores=$(cat /proc/cpuinfo | grep ^processor | wc -l)
    meminfo=($(free -h | awk '/Mem:/{print $2" "$4" "$7}'))
    storage=($(df -h /chemotion/app/ | tail -n1 | awk '{print $2" "$4}'; df -h /shared/ | tail -n1 | awk '{print $2" "$4}'))

    rubyVersion="$(odus 'ruby --version | head')"
    passengerVersion="$(odus 'passenger --version | head')"
    nodeVersion=$(odus 'echo "Node $(node --version)" | head')
    npmVersion=$(odus 'echo "NPM $(npm -v)" | head')

    chemotionVersion=${CHEMOTION_VERSION}-${FLAVOR}

    info "System information:"
    echo " - CPU Cores: $cores"
    echo " - Memory:"
    echo "    - ${meminfo[0]} (total)  ${meminfo[1]} (free)"
    echo " - Storage:"
    echo "    - root  : ${storage[0]} (total)  ${storage[1]} (free)"
    echo "    - hared : ${storage[2]} (total)  ${storage[2]} (free)"
    echo ""
    info "Used versions:"
    echo " - $rubyVersion"
    echo " - $passengerVersion"
    echo " - $nodeVersion"
    echo " - $npmVersion"
    echo " - $chemotionVersion" 
}

initalizeContainer() {
    export INITIALIZE=yes
    waitForDB
    execute "${INIT_BASE}/init-scripts/library/shared-storage-init.sh"
    execute "${INIT_BASE}/init-scripts/library/db-init.sh"
    execute "${INIT_BASE}/init-scripts/library/block-pubchem.sh"
    execute "${INIT_BASE}/init-scripts/library/eln-init.sh"
    execute "${INIT_BASE}/init-scripts/library/unblock-pubchem.sh"
}

upgradeContainer() {
    waitForDB
    execute "${INIT_BASE}/init-scripts/library/eln-upgrade.sh"
}

startELN() {
    export ELN_ROLE=app
    waitForDB
    execute "${INIT_BASE}/init-scripts/library/run.sh"
}

startBGWorker() {
    export ELN_ROLE=worker
    waitForDB
    execute "${INIT_BASE}/init-scripts/library/run.sh"
}

function usage()
{
    echo "Usage: test"
    echo -e "$(basename $0) COMMAND\n"
    echo "Commands:"
    echo "         info: display basic info." 
    echo "         init: initialize the shared storage and database.sh"
    echo "               WARNING: This will delete all data in your database"
    echo "                        and create a new one."
    echo "      upgrade: upgrade the container. Run this is the version changed."
    echo "    start-eln: starts container as ELN frontend"
    echo " start-worker: starts container as background worker"
    echo ""
}

#chown -R ${PROD}:${PROD} ./etc/init/reconfig.sh

case "$1" in
    reconfig)
        updateParam
	;;
    info)
        containerInfo
        ;;
    init)
        confirm "This is a destructive action and data will be lost! Type 'continue' to go on, anything else to abort" "continue" && ensureRoot && initalizeContainer
        ;;
    upgrade)
        ensureRoot && upgradeContainer
        ;;
    start-eln)
        startELN
        ;;
    start-worker)
        startBGWorker
        ;;
    cmd | bash)
        exec /bin/bash
        ;;
    *)
        usage
        ;;
esac
