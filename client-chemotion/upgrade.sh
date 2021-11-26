#!/bin/bash

echo "You are about to upgrade your existing Chemotion ELN setup."
echo "Please make sure that you put the new version of docker-compose.yml parallel to this upgrade script."
echo "Be aware that this script might modifiy existing data and make sure that you have created a backup of your existing data."
read -e -p "
Do you wish to proceed ? [yes/N] " YN
[[ $YN != "yes" ]] && exit 1

sharedTmp=$(date +"%FT%H%M")_shared
if [ -d shared/eln ]; then
   echo "Folder shared/eln already exists. Has an upgrade already been perfromed for this Chemotion ELN setup?"
   exit 1
fi
if [ -d shared ]; then
    mv shared/ $sharedTmp && mkdir shared && mv $sharedTmp shared/eln/
else
   echo "Folder shared/ does not exist. Are you in the right folder of a <1.0.3 Chemotion ELN setup?"
   exit 1
fi

./setup.sh

if [ -f shared/eln/config/database.yml ]; then
    oldLandscape=$(date +"%FT%H%M")_old
    mkdir -p shared/landscapes/$oldLandscape/config
    read -e -p "
    database.yml exists in your old setup. Do you want to transfer it to the new setup ? [yes/N] " YN
    [[ $YN == "yes" ]] &&  cp shared/eln/config/database.yml shared/landscapes/$oldLandscape/config/
    read -e -p "
    Overwrite configuration files with default configuration files ? [yes/N] " YN
    [[ $YN == "yes" ]] && docker-compose run eln landscape deploy --name $oldLandscape
else
    read -e -p "
    Overwrite configuration files with default files ? [yes/N] " YN
    [[ $YN == "yes" ]] && docker-compose run eln landscape deploy
fi

secret="production:\n   secret_key_base: $(date +"%FT%H%M")"

read -e -p "
    Do you want to update the secrets.yml file ? [yes/N] " YN
[[ $YN == "yes" ]] && echo -e $secret > shared/eln/config/secrets.yml

read -e -p "
    Run upgrade script (generating secret, migrating database, generating sprites, compiling assets) ? [yes/N] " YN
[[ $YN == "yes" ]] && docker-compose run eln upgrade
