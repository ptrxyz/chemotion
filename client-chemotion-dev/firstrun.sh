#!/bin/bash

git clone https://github.com/ComPlat/chemotion_ELN.git src
cd src && git checkout development ; cd ..

mkdir src/node_modules

echo "npm install ." | docker-compose run eln user-shell
docker-compose run eln init-dev

echo "Container should be ready for development, find the sources in [$(pwd)/src/]."
echo "Happy coding & welcome to the team! :)"
