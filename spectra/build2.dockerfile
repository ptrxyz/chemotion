# GLOBALs that should be passed to this file
ARG VERSION=UNSET

# Derived values
ARG BASE=chemotion-build/base:${VERSION}

# Private variables, not passed in from outside, but helpful for this file
# -


# Stage 0: import this for later
FROM chemotion-build/base:${VERSION} as chemotion-build-base

# Stage 1: prepare the base image
# hadolint ignore=DL3006
#FROM chambm/pwiz-skyline-i-agree-to-the-vendor-licenses as prebuild
FROM proteowizard/pwiz-skyline-i-agree-to-the-vendor-licenses as prebuild
RUN apt-get -y update && \
    apt-get -y upgrade || true

# 23-08-19: Workaround Ubuntu bug
# https://github.com/termux/proot-distro/issues/90
RUN rm -f /var/lib/dpkg/info/fprintd.postinst; \
    rm -f /var/lib/dpkg/info/libfprint-2-2*.postinst; \
    rm -f /var/lib/dpkg/info/libpam-fprintd*.postinst; \
    dpkg --configure -a

RUN apt-get install -y --no-install-recommends --autoremove --fix-missing python3-flask python3-gevent curl

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

FROM scratch as msconvert
COPY --from=prebuild / /
COPY --from=chemotion-build-base /tini /tini
COPY ./additives/mscrunner.py /app/mscrunner.py

RUN mkdir -p /shared && rm -rf /data && \
    ln -s /shared /data

# FROM scratch as final
# COPY --from=msconvert / /
# SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c", "--"]

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    CONTAINER_GITHUB=https://github.com/ProteoWizard/container \
    WINEDISTRO=devel \
    WINEVERSION=7.8~focal-1 \
    WINEDEBUG=-all \
    WINEPREFIX=/wineprefix64 \
    WINEPATH=C:\\pwiz;C:\\pwiz\\skyline

# Stage 2: finalize the image
FROM msconvert as app
ARG VERSION
ARG BUILDRUN
ARG MSC_PORT

ENV FLASK_ENV=production \
    FLASK_DEBUG=0 \
    MSC_PORT=4000

EXPOSE 4000

WORKDIR "/app"
ENTRYPOINT ["/tini", "--"]
CMD ["/bin/bash", "-c", "wine msconvert &>/dev/null; exec python3 -u mscrunner.py"]

HEALTHCHECK --interval=15s --timeout=3s --start-period=30s --retries=3 \
    CMD curl --fail http://localhost:4000/ping || exit 1

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion MSConvert" \
    "org.opencontainers.image.description"="Image for the Chemotion Spectra-MSConvert sidecar container" \
    "org.opencontainers.image.version"="${VERSION}" \
    "chemotion.internal.service.id"="msconvert"
