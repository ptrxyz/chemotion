#!/bin/bash
set -a

log() {
    TAG="\x1B[35m[${TASKNAME}]\x1B[0m"
    echo -en "$TAG "
    echo $@
}
