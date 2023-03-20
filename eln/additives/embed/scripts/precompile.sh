#!/bin/bash
set -euo pipefail

LOGFILE=${LOGFILE-/precompile.log}
export LOGFILE
cd /chemotion/app || exit 200

log "Precompiling assets. This could take a while ..."
bundle exec rake assets:precompile &>>"${LOGFILE}"
