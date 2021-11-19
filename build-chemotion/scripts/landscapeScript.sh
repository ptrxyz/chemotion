#!/bin/bash

checkFileExists(){
    if [[ ! -f $1 ]]; then
        echo "    File $1 does not exist. Please create it."
        return 1
    else
        echo "    Found file $1."
        return 0
    fi
}

checkFolderExists(){
    if [[ ! -d $1 ]]; then
        echo "    Folder $1 does not exist. Please create it."
        return 1
    else
        echo "    Found folder $1."
        return 0
    fi
}

checkFolderIsWritable(){
    if [[ ! -w $1 ]]; then
        echo "    Folder $1 has no write permission. Please grant it."
        return 1
    else
        echo "    Folder $1 has write permission."
        return 0
    fi
}

copyLandscape(){
    if ! checkFolderExists "/shared/landscapes/$1"        ; then exit 1; fi

    echo -e "    >>>> Copying configuration files from landscape [$1] to setup ...\n"
    if checkFolderExists "/shared/landscapes/$1/config" ; then cp -r /shared/landscapes/$1/config/*  /shared/eln/config/;  else echo "        Skipping..." ; fi
    if checkFolderExists "/shared/landscapes/$1/log"    ; then cp -r /shared/landscapes/$1/log/*     /shared/eln/log/;     else echo "        Skipping..." ; fi
    if checkFolderExists "/shared/landscapes/$1/public" ; then cp -r /shared/landscapes/$1/public/*  /shared/eln/public/;  else echo "        Skipping..." ; fi
    if checkFolderExists "/shared/landscapes/$1/tmp"    ; then cp -r /shared/landscapes/$1/tmp/*     /shared/eln/tmp/;     else echo "        Skipping..." ; fi
    if checkFolderExists "/shared/landscapes/$1/uploads"; then cp -r /shared/landscapes/$1/uploads/* /shared/eln/uploads/; else echo "        Skipping..." ; fi
    if checkFileExists   "/shared/landscapes/$1/.env"   ; then cp /shared/landscapes/$1/.env         /shared/eln/;         else echo "        Skipping..." ; fi

    return 0
}

copyDefaultLandscape(){
    if ! checkFolderExists "/template/defaultLandscape"        ; then exit 1; fi
    if ! checkFolderExists "/template/defaultLandscape/config" ; then exit 1; fi
    # if ! checkFolderExists "/template/defaultLandscape/log"    ; then exit 1; fi
    if ! checkFolderExists "/template/defaultLandscape/public" ; then exit 1; fi
    # if ! checkFolderExists "/template/defaultLandscape/tmp"    ; then exit 1; fi
    # if ! checkFolderExists "/template/defaultLandscape/uploads"; then exit 1; fi
    if ! checkFileExists "/template/defaultLandscape/.env"; then exit 1; fi
    
    echo -e "    >>>> Copying configuration files from default landscape to setup ...\n"

    SRC="/template/defaultLandscape"
    DEST="/shared/eln"
    tar c --directory ${SRC} . | tar xv --one-top-level=${DEST}
    find ${DEST}

    return 0
}

echo "---- Script deploys landscape [$1] ----"
if [[ $2 == "nodefault" ]]; then
    echo "---- not based on the default landscape ----" 
else
    echo "---- based on the default landscape ----"
fi

if ! mount | grep "on /shared" 2>&1 1>/dev/null; then
    echo "    The shared folder is not correctly connected as volume. Please make sure that a folder shared/ is available next to the docker-compose.yml file."
    exit 1
else
    echo "    Found folder /shared."
fi

if [[ $2 == "default" ]]; then
    copyDefaultLandscape
fi

if [[ $1 != "default" ]]; then
    copyLandscape $1
fi