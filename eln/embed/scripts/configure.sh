#!/bin/bash
log "Loading configuration."
tar cH posix --exclude='*backup*' --exclude='ignore' --exclude='import' --exclude='.version' --directory=/shared . | tar xh --one-top-level=/chemotion/app
