#!/bin/bash
echo Running the upgrade script 

cd /chemotion/app/

echo "    Executing migrations..."
bundle exec rake db:migrate
echo "    Database migrated."

echo "    Creating sprites..."
bundle exec rake ketcherails:import:common_templates
rm -rf /chemotion/app/app/public/images/ketcherails/icons/original/*
bundle exec rails r 'MakeKetcherailsSprites.perform_now'

bundle exec rake assets:precompile