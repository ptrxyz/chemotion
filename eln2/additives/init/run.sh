#!/bin/bash
cd /chemotion/app

# make sure we really do have all the packages we need
yarn install --ignore-engines
bundle check || bundle install --jobs=$(getconf _NPROCESSORS_ONLN) --path "${BUNDLE_PATH}"

# if the DATABASE env variable is set, rewrite the config file to use DATABASE
if [[ -n "${DATABASE}" ]]; then
  yq -o json '.development.database="'"${DATABASE-chemotion}"'"' ./config/database.yml | \
  yq -P -o yaml > /tmp/database.yml
  mv /tmp/database.yml ./config/database.yml
fi

# source database environment and wait for it to spin up
. /init/dbenv.sh
while ! pg_isready; do echo "Waiting for database."; sleep 3; done

# set DO_RESET env variable to make everything go away.

# heuristics to determine if there is user data present on the file system
if [[ "$(find ./public/images/molecules/ ./public/images/reactions/ ./public/images/samples/ | grep -v .keep | wc -l)" -gt 3 ]]; then
  USERDATA="present"
fi

if [[ -n "${DO_RESET}" ]]; then
  echo "DO_RESET triggered!"
  echo "Your data will be gone in..."
  for i in $(seq 10 -1 1); do
    echo "${i}..."
    sleep 1
  done

  bundle exec rake db:reset || bundle exec rake db:schema:load
  rm -rf ./uploads || true

  if [[ -f ./reset.tar.gz ]]; then
    rm -rf ./public/images || true
    mkdir -p ./public/images
    tar xvz -f ./reset.tar.gz .
    (cd ./public/images && find .)
  fi
fi

# house keeping...
bundle exec rake db:create
bundle exec rake db:migrate

# heuristically determine if database needs seeding.
seedCheck=$(psql -tA -c "select 'non-empty' from molecule_names limit 1;")
[[ "${seedCheck}" =~ "non-empty" ]] || bundle exec rake db:seed

# add ketcherails templates if needed
seedCheck=$(psql -tA -c "select 'non-empty' from ketcherails_common_templates limit 1;")
if [[ ! "${seedCheck}" =~ "non-empty" ]]; then
  bundle exec rake ketcherails:import:common_templates
  # rm -rf ./public/images/ketcherails/icons/original/*
  bundle exec rails r 'MakeKetcherailsSprites.perform_now'
fi

bundle exec rake assets:precompile

# Experimental: enable web console inside container
sed -i '$d' ./config/environments/development.rb
cat >> ./config/environments/development.rb <<EOF
  config.web_console.whitelisted_ips = (ENV['WEBCONSOLE_WHITELISTED_IPS'] || "").split(' ')
end
EOF

# trying to find our IP without `iputils`
export APPLICATION_URL=${HOST_URL}
export WEBCONSOLE_WHITELISTED_IPS="$(awk '/0.0\/[123456789]/ { print $2 }' /proc/net/fib_trie | sort | uniq)"

# fire it up ...
rm -f ./tmp/pids/server.pid
(sleep 5 ; bundle exec rake jobs:work) &
exec bundle exec rails s -b 0.0.0.0 -p${RAILS_PORT-4000}
