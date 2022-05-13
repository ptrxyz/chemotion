#!/bin/bash
log "Precompiling assets. This could take a while ..."
cd /chemotion/app && bundle exec rake assets:precompile &>>$LOGFILE
