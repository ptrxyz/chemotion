#!/bin/bash

sudo rm -R db-data/ 

VERSION=$(date +%y.%m)-2
PID=$(pwd | md5sum | cut -c -16)
GITHASH=$(cd src; git log --pretty=format:'%h' -n 1)
GITTAG=$(cd src; git describe --exact-match --tags ${GITHASH} || echo ${GITHASH})

ok() {
	CC_GREEN='\033[0;32m'
	CC_NC='\033[0m'
	echo -ne "${CC_GREEN}"
	echo $@
	echo -ne "${CC_NC}"
}

error() {
	CC_RED='\033[0;31m'
	CC_NC='\033[0m'
	echo -ne "${CC_RED}"
	echo $@
	echo -ne "${CC_NC}"
}

info() {
	CC_CYAN='\033[0;36m'
	CC_NC='\033[0m'
	echo -ne "${CC_CYAN}"
	echo $@
	echo -ne "${CC_NC}"
}

echo "Building [$VERSION] based on [$GITTAG]. Press enter."
read

docker rm -f db 2>/dev/null
docker network rm ${PID}_net 2>/dev/null

docker network create ${PID}_net >/dev/null && info "Build network created."
docker run -d --name db --network "${PID}_net" -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=root postgres >/dev/null && info "Created database for the build process."
echo "Waiting for database to spin up ..."
sleep 10
time docker build --build-arg CHEMOTION_VERSION=${VERSION}@${GITTAG} --network "${PID}_net" -t ptrxyz/chemotion:${VERSION} $@ . && (
	ok "Image successfully built."
	docker tag ptrxyz/chemotion:${VERSION} ptrxyz/chemotion:latest-local
	echo -e "Version: \033[0;97m$VERSION\033[0m"
) || (
	error "Build process failed."	
)
docker rm -f db 2>&1 1>/dev/null && info "Build database removed."
docker network rm ${PID}_net 2>&1 1>/dev/null && info "Build network removed."

info "Done."
