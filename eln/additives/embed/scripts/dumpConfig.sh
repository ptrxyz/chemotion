#!/bin/bash

files=(
    config/database.yml
    config/converter.yml
    config/datacollectors.yml
    config/editors.yml
    config/ketcher_service.yml
    config/shrine.yml
    config/spectra.yml
    config/storage.yml
    .env
)

mkdir -p /shared/dump/config
for i in "${files[@]}"; do
    cp "/chemotion/app/${i}" "/shared/dump/${i}" || true
done
