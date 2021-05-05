#!/bin/bash

VERSION=$(date +%+2y.%+2m)-3
PID=$(pwd | md5sum | cut -c -16)
GITHASH=$(cd src; git log --pretty=format:'%h' -n 1)
GITTAG=$(cd src; git describe --exact-match --tags ${GITHASH})

echo "Building [$VERSION] based on [$GITHASH ($GITTAG)]. Press enter."
read

docker rm -f db
docker network rm ${PID}_net

docker network create ${PID}_net
docker run -d --name db --network "${PID}_net" -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=root postgres
sleep 10
time docker build --build-arg CHEMOTION_VERSION=${VERSION}@$(git log -1 --pretty='%h') --network "${PID}_net" -t ptrxyz/chemotion:${VERSION} $@ .
docker rm -f db
docker network rm ${PID}_net

docker tag ptrxyz/chemotion:${VERSION} ptrxyz/chemotion:latest-local
echo "Done. Version: $VERSION"
