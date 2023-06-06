#!/bin/bash

#
# run me like this: cat build.sh | docker run -i --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd)/add:/add ubuntu
#

TAG=$1
TAG=${TAG-0.3.0}

echo "Tag to use: $TAG"

set -e

apt-get -y update && apt-get install -y docker.io
git clone https://github.com/NFDI4Chem/nmrium-react-wrapper.git build
cd build
git checkout "v${TAG}" # 1f1530cad7886b9b1d45f2650629dd70f899ef42  # development?
# cp /patches/*.patch /build

# for patch in *.patch; do
# 	echo "Patching: ${patch}..."
# 	git apply "${patch}" || true
# done

cp /add/allowed-origins.json /build/src/

docker build -f Dockerfile.prod -t nmrium:"${TAG}" -t registry.hydrogen.scc.kit.edu/pk01/nmrium:"${TAG}" .
