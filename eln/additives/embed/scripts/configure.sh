#!/bin/bash
log "Loading configuration."
tar cH posix \
    --exclude='*backup*'    \
    --exclude='ignore'      \
    --exclude='import'      \
    --exclude='.version'    \
    --exclude='lost+found'  \
    --directory=/shared .   \
| tar xh --one-top-level=/chemotion/app

touch /chemotion/app/.env
