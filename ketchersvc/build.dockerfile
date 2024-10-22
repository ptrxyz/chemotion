# GLOBALs that should be passed to this file
ARG VERSION=UNSET
ARG BUILD_TAG_KETCHERSVC=UNSET

# Derived values
# -

# Private variables, not passed in from outside, but helpful for this file
# -


# Stage 1: prepare the base image
# hadolint ignore=DL3006
FROM oven/bun:latest as ketcher-base
WORKDIR /app


RUN apt-get update && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing \
    git ca-certificates

RUN git clone --depth 1 https://github.com/ptrxyz/chemotion-ketchersvc /app

RUN ELECTRON_SKIP_BINARY_DOWNLOAD=1 bun install && \
    bunx playwright install && \
    bun run build


# Stage 2: Finalize the image
FROM oven/bun:latest as app
ARG VERSION

EXPOSE 4000
WORKDIR /app

COPY --from=ketcher-base /app/dist /app
COPY --from=ketcher-base /root/.cache/ms-playwright/ /root/.cache/ms-playwright/
RUN bunx playwright install-deps chromium
# RUN bunx playwright install --with-deps chromium

HEALTHCHECK --interval=5s --timeout=1s --start-period=30s --retries=3 \
    CMD /usr/bin/pidof bun

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion KetcherSVC" \
    "org.opencontainers.image.description"="Image for Chemotion KetcherSVC" \
    "org.opencontainers.image.version"="${VERSION}" \    
    "chemotion.internal.service.id"="ketchersvc"

CMD ["bun", "index.js"]
