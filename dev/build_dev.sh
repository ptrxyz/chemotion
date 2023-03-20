#!/bin/bash
cd /chemotion/app || exit 1

# mirror the production sections of all config/*.yml to a development section
mkdir -p /tmp/prod/ /tmp/dev/ && cp config/*.yml /tmp/prod
for i in /tmp/prod/*.yml; do
	fname=$(basename "${i}")
	j=/tmp/dev/"${fname}"
	cat "${i}" > "${j}"
	cat "${i}" | sed 's/^production:\s*$/development:/g' >> "${j}"
done
cp /tmp/dev/*.yml /chemotion/app/config/
rm -rf /tmp/dev /tmp/prod/

# install dev packages
PROC_ONLINE=$(getconf _NPROCESSORS_ONLN)
PROC_ONLINE=${PROC_ONLINE-4}

# make sure we really do have all the packages we need
yarn install --ignore-engines
sed -i -e "/gem 'faker'/d" /chemotion/app/Gemfile
echo "gem 'faker'" >> /chemotion/app/Gemfile
bundle check || bundle install --jobs="${PROC_ONLINE}" --path "${BUNDLE_PATH}"


# expose logs
mkdir -p log/ && ln -s ../log/server.log public/
mkdir -p log/ && ln -s ../log/worker.log public/

# make dev seeds the prod seeds to have some users
# cp db/seeds/development.rb db/seeds/production.rb
echo "" >> db/seeds/production.rb
echo 'require_relative "development/users.seed.rb"' >> db/seeds/production.rb
