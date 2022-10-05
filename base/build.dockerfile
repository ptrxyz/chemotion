ARG BASE=ubuntu:21.04
ARG BASE_VERSION=1.4.0
ARG BUILDRUN

ARG TINI_VERSION="v0.19.0"

####################################################################################################
# BASE
####################################################################################################
FROM ${BASE} AS base

# set timezone
ARG TZ=Europe/Berlin
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime

# locales
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8
RUN echo -e "LANG=${LANG}\nLC_ALL=${LANG}" > /etc/locale.conf && \
    echo "${LANG} UTF-8" > /etc/locale.gen

RUN sed -i -re 's/([a-z]{2}\.)?archive.ubuntu.com|security.ubuntu.com/old-releases.ubuntu.com/g' /etc/apt/sources.list

# install system packages
RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install -y --no-install-recommends locales && \
    locale-gen en_US.UTF-8

# RUN apt-get install -y software-properties-common && apt-add-repository -y ppa:rael-gc/rvm && apt-get update
# RUN apt-get install -y --allow-downgrades libssl-dev=1.1.1l-1ubuntu1.4

ARG TINI_VERSION
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

FROM base
ARG BASE_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Converter" \
    "org.opencontainers.image.description"="Image for Chemotion Converter" \
    "org.opencontainers.image.version"="${BASE_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="base"
