FROM ubuntu:latest

RUN sed -i 's~http://archive.ubuntu.com/ubuntu/~http://de.archive.ubuntu.com/ubuntu/~g' /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y upgrade && apt-get -y autoremove

SHELL ["/bin/bash", "-c"]

ENV TIMEZONE=Europe/Berlin

ARG CHEMOTION_VERSION
ARG FLAVOR=production

RUN AREA=$(echo $TIMEZONE | cut -d'/' -f1) && \
    REGION=$(echo $TIMEZONE | cut -d'/' -f2) && \
    rm -f /etc/localtime && \
    ln -s /usr/share/zoneinfo/$AREA/$REGION /etc/localtime

RUN apt-get -y --autoremove --fix-missing install \
    apt-transport-https apt-utils autoconf automake \
    bison ca-certificates cmake cron curl dirmngr \
    fonts-crosextra-caladea fonts-crosextra-carlito \
    fonts-dejavu fonts-dejavu-core fonts-dejavu-extra \
    fonts-liberation fonts-liberation2 fonts-linuxlibertine \
    fonts-noto-core fonts-noto-extra fonts-noto-ui-core \
    fonts-opensymbol fonts-sil-gentium fonts-sil-gentium-basic \
    g++ gconf-service git gnupg gnupg2 htop imagemagick \
    inkscape libboost-all-dev libeigen3-dev libffi-dev \
    libgconf-2-4 libgdbm-dev libgmp-dev libmagic-dev \
    libmagickcore-dev libmagickwand-dev libncurses5-dev \
    libnspr4 libnss3 libpango1.0-0 libpq-dev libreadline-dev \
    libsass-dev libsqlite3-dev libssl-dev libxml2-dev \
    libxslt-dev libxss1 libyaml-dev openssh-server p7zip-full \
    pandoc postgresql postgresql-client postgresql-contrib \
    python-dev ranger s3fs sqlite3 sudo swig tzdata ufw \
    xfonts-100dpi xfonts-75dpi xfonts-base xfonts-cyrillic \
    xfonts-scalable rsyslog tini make gcc vim

## Pandoc
ENV PANDOC_VERSION=2.10.1
RUN pandoc_pkg="pandoc-${PANDOC_VERSION}-1-amd64.deb" && \
    curl -o /tmp/${pandoc_pkg} -L https://github.com/jgm/pandoc/releases/download/${PANDOC_VERSION}/${pandoc_pkg} && \
    dpkg -i /tmp/${pandoc_pkg}

## Passenger. Use own repos for standalone version
## As of version 21.04-1 this is deprecated and we use `gem install` to install passenger standalone.
#
# RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 561F9B9CAC40B2F7
# RUN . /etc/lsb-release && echo "deb https://oss-binaries.phusionpassenger.com/apt/passenger ${DISTRIB_CODENAME} main" >/etc/apt/sources.list.d/passenger.list
# RUN apt-get -y update && apt-get -y --autoremove --fix-missing install passenger=1:6.0.5-1~focal1

## MinIO client
RUN curl -L -o /bin/mcli https://dl.min.io/client/mc/release/linux-amd64/mc && chmod +x /bin/mcli

## User
ENV PROD=chemotion
ENV PROD_HOME=/home/chemotion

## RUBY
ENV RUBY_VERSION=2.6.6
ENV BUNDLER_VERSION=1.17.3
ENV PASSENGER_VERSION=6.0.5
ENV PASSENGER_PORT=4000

## NODEJS
ENV NVM_VERSION=v0.35.3
ENV NODE_VERSION=12.21.0
ENV NPM_VERSION=7.6.2

## POSTGRESQL DB
ENV DB_ROLE=chemotion
ENV DB_NAME=chemotion
ENV DB_PW=changeme
ENV DB_HOST=db
ENV DB_PORT=5432

## Basics
ARG NCPU=24
ARG SECRET_KEY_BASE=SecretKey
ENV APP_NAME=chemotion
ENV APP_DIR=/${APP_NAME}
ENV RAILS_ENV=${FLAVOR}

RUN useradd -mUs /bin/bash -d ${PROD_HOME} ${PROD}
RUN mkdir -p ${APP_DIR}/app
RUN mkdir -p /shared

COPY ./src/ ${APP_DIR}/app

RUN chown -R ${PROD}:${PROD} ${APP_DIR}
RUN chown -R ${PROD}:${PROD} /shared

## We do not have cron installed, we do not need logrotate...
RUN echo -e ""\
    "${APP_DIR}/app/log/*.log {\n"\
    "  weekly\n"\
    "  missingok\n"\
    "  rotate 8\n"\
    "  compress\n"\
    "  delaycompress\n"\
    "  notifempty\n"\
    "  copytruncate\n"\
    "}\n" | sed 's/^ //g' >/etc/logrotate.d/${APP_NAME}

RUN curl -fsSL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y nodejs
RUN npm install -g npm@${NPM_VERSION}

RUN npm --version
RUN node --version

USER ${PROD}
RUN gpg --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB || true
RUN gpg2 --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB || true
RUN curl -sSL https://get.rvm.io | bash -s stable --ruby=${RUBY_VERSION} --auto-dotfiles
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && gem install bundler -v ${BUNDLER_VERSION}

# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/${NVM_VERSION}/install.sh | bash
# RUN . ~/.nvm/nvm.sh; nvm install ${NODE_VERSION}
# RUN . ~/.nvm/nvm.sh && nvm use ${NODE_VERSION} && npm install -g npm@${NPM_VERSION}

RUN cd ${APP_DIR}/app && \
    npm install .

RUN echo ${RUBY_VERSION} > ${APP_DIR}/app/.ruby-version
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    gem install passenger -v ${PASSENGER_VERSION}
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle config build.nokogiri --use-system-libraries
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle install --jobs ${NCPU}

RUN cd ${APP_DIR}/app && \
    mkdir -p backup/deploy_backup backup/weekly_backup \
    config log public/simulations public/assets public/images \
    tmp/pids tmp/sockets tmp/cache uploads

RUN echo -e "SECRET_KEY_BASE='${SECRET_KEY_BASE}'\n"\
    "DB_NAME='${DB_NAME}'\n"\
    "DB_ROLE='${DB_ROLE}'\n"\
    "DB_PW='${DB_PW}'\n"\
    "DB_HOST='${DB_HOST}'\n"\
    "DB_PORT=${DB_PORT}\n"\
    "RAILS_ENV=${RAILS_ENV}\n" | sed 's/^ //g' >>${APP_DIR}/app/.env 

RUN echo -e "${RAILS_ENV}:\n"\
    "  adapter: postgresql\n"\
    "  encoding: unicode\n"\
    "  database: <%=ENV['DB_NAME']%>\n"\
    "  pool: 5\n"\
    "  username: <%=ENV['DB_ROLE']%>\n"\
    "  password: <%=ENV['DB_PW']%>\n"\
    "  host: <%=ENV['DB_HOST']%>\n"\
    "  port: <%=ENV['DB_PORT']%>\n" | sed 's/^ //g' >${APP_DIR}/app/config/database.yml

RUN echo -e "${RAILS_ENV}:\n"\
    "  :primary_store: 'local'\n"\
    "  :secundary_store: ''\n"\
    "  :stores:\n"\
    "    :tmp:\n"\
    "      :data_folder: 'tmp/uploads/'\n"\
    "      :thumbnail_folder: 'tmp/uploads/'\n"\
    "    :local:\n"\
    "      :data_folder: 'uploads'\n"\
    "      :thumbnail_folder: 'uploads'\n" | sed 's/^ //g' >${APP_DIR}/app/config/storage.yml

RUN echo -e "${RAILS_ENV}:\n"\
    "  :name_abbreviation:\n"\
    "    :reserved_list:\n"\
    "      - 'CRR'\n"\
    "      - 'CRS'\n"\
    "      - 'CRD'\n"\
    "    :length_group: [2, 5]\n"\
    "    :length_device: [2, 6]\n"\
    "    :length_default: [2, 3]\n" | sed 's/^ //g' >${APP_DIR}/app/config/user_props.yml

RUN echo -e "${RAILS_ENV}:\n"\
    "  secret_key_base: <%=ENV['SECRET_KEY_BASE']%>\n" | sed 's/^ //g' >${APP_DIR}/app/config/secrets.yml

RUN echo -e "${RAILS_ENV}:\n"\
    "  :services:\n"\
    "    - :name: 'mailcollector'\n"\
    "      :every: 5 # minutes\n"\
    "    - :name: 'folderwatchersftp'\n"\
    "      :every: 5 # minutes\n"\
    "    - :name: 'folderwatcherlocal'\n"\
    "      :every: 5 # minutes\n"\
    "    - :name: 'filewatchersftp'\n"\
    "      :every: 2 # minutes\n"\
    "    - :name: 'filewatcherlocal'\n"\
    "      :every: 2 # minutes\n"\
    "\n"\
    "  :mailcollector:\n"\
    "    :server: 'imap.server.de'\n"\
    "    :mail_address: 'service@mail'\n"\
    "    :password: 'password'\n"\
    "    # :port: 993 default\n"\
    "    # :ssl: true default\n"\
    "    :aliases:\n"\
    "      - 'alias_one@kit.edu'\n"\
    "      - 'alias_two@kit.edu'\n"\
    "\n"\
    "  :sftpusers:\n"\
    "    - :user: 'user1'\n"\
    "      :password: 'pass'\n"\
    "    - :user: 'user2'\n"\
    "      :password: 'pass'\n"\
    "\n"\
    "  :localcollectors:\n"\
    "    - :path: '/home/ftpuser'\n"\
    "    - :path: '/home/eln/public'\n" | sed 's/^ //g' >${APP_DIR}/app/config/datacollectors.yml

RUN cat ${APP_DIR}/app/config/datacollectors.yml

RUN chmod 600 ${APP_DIR}/app/config/*.yml
RUN chmod 600 ${APP_DIR}/app/.env

USER root

RUN psql --host="${DB_HOST}" -c "DROP DATABASE IF EXISTS $DB_NAME;"
RUN psql --host="${DB_HOST}" -c "DROP ROLE IF EXISTS $DB_ROLE;"
RUN psql --host="${DB_HOST}" -c "CREATE ROLE $DB_ROLE LOGIN CREATEDB NOSUPERUSER PASSWORD '$DB_PW';"

RUN psql --host="${DB_HOST}" -c "CREATE DATABASE $DB_NAME OWNER $DB_ROLE;"
RUN psql --host="${DB_HOST}" --username="${DB_ROLE}" -c "CREATE EXTENSION IF NOT EXISTS \"pg_trgm\"; CREATE EXTENSION IF NOT EXISTS \"hstore\"; CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
RUN psql --host="${DB_HOST}" --username="${DB_ROLE}" -c "ALTER USER $DB_ROLE PASSWORD '$DB_PW';"

USER ${PROD}

RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle exec rake db:create
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle exec rake db:migrate
# RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
#     bundle exec rake db:seed
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle exec rake assets:precompile
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle exec rake ketcherails:import:common_templates
RUN rm -rf ${APP_DIR}/app/public/images/ketcherails/icons/original/*
RUN . ~/.rvm/scripts/rvm && rvm use ${RUBY_VERSION} && cd ${APP_DIR}/app && \
    bundle exec rails r 'MakeKetcherailsSprites.perform_now'


## Seed files:
# we link some shared folders here, 
RUN cd ${APP_DIR}/app && \
    mkdir -p ${APP_DIR}/seed && \
    ln -s /shared/uploads/         uploads/ && \
    ln -s /shared/tmp_uploads/     tmp/uploads

# move files to seed folder
RUN cd ${APP_DIR}/app && \
    mkdir -p ${APP_DIR}/seed ${APP_DIR}/seed/public/ ${APP_DIR}/seed/config/ && \
    mv ${APP_DIR}/app/public/zip                     ${APP_DIR}/seed/public/zip &&\
    mv ${APP_DIR}/app/public/images                  ${APP_DIR}/seed/public/images && \
    mv ${APP_DIR}/app/.env                           ${APP_DIR}/seed/.env && \
    mv ${APP_DIR}/app/config/database.yml            ${APP_DIR}/seed/config/database.yml && \
    mv ${APP_DIR}/app/config/storage.yml             ${APP_DIR}/seed/config/storage.yml && \
    mv ${APP_DIR}/app/config/user_props.yml          ${APP_DIR}/seed/config/user_props.yml && \
    mv ${APP_DIR}/app/config/secrets.yml             ${APP_DIR}/seed/config/secrets.yml && \
    mv ${APP_DIR}/app/config/datacollectors.yml      ${APP_DIR}/seed/config/datacollectors.yml && \
    mv ${APP_DIR}/app/config/editors.yml.example     ${APP_DIR}/seed/config/editors.yml && \
    mv ${APP_DIR}/app/config/inference.yml.example   ${APP_DIR}/seed/config/inference.yml && \
    mv ${APP_DIR}/app/config/spectra.yml.example     ${APP_DIR}/seed/config/spectra.yml


RUN cd ${APP_DIR}/app && \
    ln -s /shared/public/zip                ${APP_DIR}/app/public/zip &&\
    ln -s /shared/public/images             ${APP_DIR}/app/public/images && \
    ln -s /shared/.env                      ${APP_DIR}/app/.env && \
    ln -s /shared/config/database.yml       ${APP_DIR}/app/config/database.yml && \
    ln -s /shared/config/storage.yml        ${APP_DIR}/app/config/storage.yml && \
    ln -s /shared/config/user_props.yml     ${APP_DIR}/app/config/user_props.yml && \
    ln -s /shared/config/secrets.yml        ${APP_DIR}/app/config/secrets.yml && \
    ln -s /shared/config/datacollectors.yml ${APP_DIR}/app/config/datacollectors.yml && \
    ln -s /shared/config/editors.yml        ${APP_DIR}/app/config/editors.yml && \
    ln -s /shared/config/inference.yml      ${APP_DIR}/app/config/inference.yml && \
    ln -s /shared/config/spectra.yml        ${APP_DIR}/app/config/spectra.yml

USER root

ENV INIT_BASE=/etc/init

RUN mkdir -p ${INIT_BASE} ${INIT_BASE}/init-scripts
COPY ./build_conf/init.sh         ${INIT_BASE}/init.sh
COPY ./build_conf/run-app.sh      ${INIT_BASE}/run-app.sh
COPY ./build_conf/run-worker.sh   ${INIT_BASE}/run-worker.sh
COPY ./build_conf/functions.sh    ${INIT_BASE}/functions.sh
COPY ./build_conf/init-scripts/*  ${INIT_BASE}/init-scripts/

RUN chmod +x ${INIT_BASE}/init.sh ${INIT_BASE}/run-app.sh ${INIT_BASE}/run-worker.sh 

# Lines needed to be compatible with docker 1.26+ versions of tini
RUN ln -s /usr/bin/tini /tini
RUN ln -s ${INIT_BASE}/init.sh /init.sh
ENTRYPOINT ["/tini", "--", "/init.sh"]