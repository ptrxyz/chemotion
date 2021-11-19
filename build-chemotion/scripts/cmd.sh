#!/bin/bash

case "$1" in
    bash)
        exec /bin/bash
        ;;
    *)
        python3 /etc/scripts/CLI.py $@
        ;;
esac
