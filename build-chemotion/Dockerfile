################################################################################
# RUBY                                                                         #
################################################################################
FROM chemotion-build:latest-ruby AS ruby

################################################################################
# NODE                                                                         #
################################################################################
FROM chemotion-build:latest-node AS node

################################################################################
# MAIN                                                                         #
################################################################################
FROM ubuntu:21.04 AS main
# set the shell. we need `[[` and `echo -e`...
SHELL ["/bin/bash", "-c"]

# set timezone
ARG TZ=Europe/Berlin
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime

# locales
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
RUN echo -e "LANG=${LANG}\nLC_ALL=${LANG}" > /etc/locale.conf && \
    echo "${LANG} UTF-8" > /etc/locale.gen

# install system packages
RUN apt-get -y update && apt-get -y upgrade
RUN apt-get -y --autoremove --fix-missing install \
    libboost-serialization1.74.0 \
    libboost-iostreams1.74.0 \
    curl `# for pandoc` \
    postgresql-client \
    inkscape `# this installs python3` \
    xvfb \
    imagemagick \
    locales \
    tini

RUN PANDOC_VERSION=2.10.1 && \
    pandoc_pkg="pandoc-${PANDOC_VERSION}-1-amd64.deb" && \
    curl -o /tmp/${pandoc_pkg} -L https://github.com/jgm/pandoc/releases/download/${PANDOC_VERSION}/${pandoc_pkg} && \
    dpkg -i /tmp/${pandoc_pkg} && rm /tmp/${pandoc_pkg}

COPY fonts.conf /etc/fonts/conf.d/99-chemotion-fontfix.conf

# install ruby + node + modules
COPY --from=ruby /ruby         /ruby
COPY --from=node /node         /node
COPY --from=node /node_modules /chemotion/app/node_modules
COPY --from=node /yarn/cache   /yarn/cache

# symlink the binaries to /bin
RUN find /ruby/bin /node/bin -not -type d | xargs -i ln -s {} /bin/

# sanity checking
RUN ruby -v && bundle -v && gem -v && rails -v && node -v && npm -v && npx -v && yarn -v

# copy the app
COPY ./src /chemotion/app/
ENV RAILS_ENV=production
ENV NODE_ENV=production
# this line is needed due to moving config folder around (where webpacks config .js files life)
# it avoids MODULE_NOT_FOUND errors during webpack compilation
ENV NODE_PATH=/chemotion/app/node_modules/
ENV NODE_OPTIONS=--max_old_space_size=4096
RUN yarn cache dir && \
    yarn config set cache-folder /yarn/cache && \
    cd /chemotion/app && \
    yarn install --link-duplicates

# check dependencies 
COPY ./lddcheck /bin/lddcheck
RUN find /ruby/lib/ruby/gems/ /node/lib/node_modules -iname '*.so' -type f -exec lddcheck \{\} \; | tee /lddlog.txt | grep not | grep -v libRD | sort | uniq

# copy config template to image
RUN mkdir -p /shared /template 

# clean up some unnecessary files
RUN find /template /chemotion/app -iname '*.gitlab' -print -delete -or -iname '*.travis' -print -delete

COPY ./defaultLandscape /template/defaultLandscape

RUN apt-get -y update && apt-get -y upgrade
RUN apt-get -y --autoremove --fix-missing install \
    python3-click \
    python3-yaml \
    python3-psutil

COPY ./scripts /etc/scripts

ENV THOR_SILENCE_DEPRECATION=1

ENV CONFIG_USERID=1000
ENV CONFIG_FIX_PERMISSIONS_ON_STARTUP=true
ENV CONFIG_PASSENGER_PORT=4000
ENV CONFIG_ROLE=eln

# Lines needed to be compatible with docker 1.26+ versions of tini
# as of Apr. 28, we use cmd.sh as init system.
RUN ln -s /usr/bin/tini /tini
RUN ln -s /etc/scripts/cmd.sh /init
ENTRYPOINT ["/usr/bin/tini", "--", "/init"]
