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

mkdir -p /chemotion/data/public/images
mkdir -p /chemotion/data/public/images/ghs
mkdir -p /chemotion/data/public/images/molecules
mkdir -p /chemotion/data/public/images/qr
mkdir -p /chemotion/data/public/images/reactions
mkdir -p /chemotion/data/public/images/research_plans
mkdir -p /chemotion/data/public/images/samples
mkdir -p /chemotion/data/public/images/thumbnail
mkdir -p /chemotion/data/public/images/wild_card

mkdir -p /chemotion/data/uploads

mkdir -p /chemotion/data/public/safety_sheets
