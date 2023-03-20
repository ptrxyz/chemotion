# GLOBALs that should be passed to this file
ARG VERSION=UNSET

# Derived values
#-

# Private variables, not passed in from outside, but helpful for this file
ARG BASE=ubuntu:22.04
ARG TINI_VERSION="v0.19.0"

# Stage 1: the base image
# hadolint ignore=DL3006
FROM ${BASE} AS base

# set timezone
ARG TZ=Europe/Berlin
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime

# locales
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8
RUN echo "LANG=${LANG}" >/etc/locale.conf && \
    echo "LC_ALL=${LANG}" >>/etc/locale.conf && \
    echo "${LANG} UTF-8" >/etc/locale.gen

# RUN sed -i -re 's/([a-z]{2}\.)?archive.ubuntu.com|security.ubuntu.com/old-releases.ubuntu.com/g' /etc/apt/sources.list

# install system packages
# hadolint ignore=DL3009
RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing locales && \
    apt-get clean && \
    locale-gen en_US.UTF-8

# RUN apt-get install -y software-properties-common && apt-add-repository -y ppa:rael-gc/rvm && apt-get update
# RUN apt-get install -y --allow-downgrades libssl-dev=1.1.1l-1ubuntu1.4

ARG TINI_VERSION
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
ADD https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 /bin/yq
RUN chmod +x /bin/yq && chmod +x /tini

FROM scratch as final
COPY --from=base / /

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

ENTRYPOINT ["/tini", "--"]

# Stage 2: finalize the image
FROM final as app
ARG VERSION

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion Converter" \
    "org.opencontainers.image.description"="Image for Chemotion Converter" \
    "org.opencontainers.image.version"="${VERSION}" \
    "chemotion.internal.service.id"="base"
