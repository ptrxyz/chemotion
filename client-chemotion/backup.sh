#!/bin/bash

read -e -p "
Stop a running ELN docker-container setup ? [yes/N] " YN
[[ $YN == "yes" ]] && docker-compose down --remove-orphans

echo "Looking for folders [ config shared db-data ] ..."

CONFIG=""
[[ -d config ]] && CONFIG="config"

SHARED=""
[[ -d shared ]] && SHARED="shared"

DBDATA=""
[[ -d db-data ]] && DBDATA="db-data"

echo "Adding [ ${CONFIG} ${SHARED} ${DBDATA} ] to the backup"

tar cfz backup.tar.gz ${CONFIG} ${SHARED} ${DBDATA}
