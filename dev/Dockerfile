ARG BASEIMG="chemotion-build/eln:latest"

# hadolint ignore=DL3006
FROM ${BASEIMG}

ENV CONFIG_ROLE=dev \
    SECRET_KEY_BASE=dev \
    DROP_DATABASE=1

COPY ./shared/pullin/config /chemotion/app/config
COPY Makefile /embed/Makefile
COPY dev.sh /embed/scripts/dev.sh
COPY run.sh /embed/run.sh
COPY health.sh /embed/health.sh

RUN chmod +x /embed/run.sh /embed/scripts/dev.sh /embed/health.sh

COPY build_dev.sh /embed
RUN chmod +x /embed/build_dev.sh && /embed/build_dev.sh && rm -r /embed/build_dev.sh
