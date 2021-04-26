#!/bin/bash

error() {
    RED='\033[0;31m'
    NC='\033[0m'
    echo -ne "${RED}"
    echo $@
    echo -ne "${NC}"
}

note() {
    YELLOW='\033[0;33m'
    NC='\033[0m'
    echo -ne "${YELLOW}"
    echo $@
    echo -ne "${NC}"
}

info() {
    CYAN='\033[0;36m'
    NC='\033[0m'
    echo -ne "${CYAN}"
    echo $@
    echo -ne "${NC}"
}

msg() {
    echo $@
}

export -f error
export -f note
export -f info
export -f msg