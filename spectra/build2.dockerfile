
# GLOBALs that should be passed to this file
ARG SPECTRA_VERSION=1.4.0
ARG BUILDRUN

# GLOBALs that should not be changed.
ARG MSC_PORT=8088

FROM chambm/pwiz-skyline-i-agree-to-the-vendor-licenses as msconvert
RUN apt-get -y update && \
    apt-get -y upgrade && \
    apt-get --no-install-recommends -y install python3-pip curl && \
    pip install flask gevent

COPY --from=chemotion-build/base:1.4.0 /tini /tini
COPY ./scripts/mscrunner.py /app/mscrunner.py

ARG MSC_PORT

RUN mkdir -p /shared && rm -rf /data && \
    ln -s /shared /data

EXPOSE ${MSC_PORT}

ENV FLASK_ENV=production \
    FLASK_DEBUG=0 \
    MSC_PORT=${MSC_PORT}

ENTRYPOINT ["/tini", "--"]
WORKDIR "/app"
CMD ["python3", "mscrunner.py"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=30s --retries=3 \
    CMD curl --fail http://localhost:${MSC_PORT}/ping || exit 1


FROM msconvert
ARG SPECTRA_VERSION
ARG BUILDRUN
LABEL "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion MSConvert" \
    "org.opencontainers.image.description"="Image for the Chemotion Spectra-MSConvert sidecar container" \
    "org.opencontainers.image.version"="${SPECTRA_VERSION}" \
    "chemotion.internal.buildrun"="${BUILDRUN}" \
    "chemotion.internal.service.id"="msconvert"
