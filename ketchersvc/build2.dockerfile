# GLOBALs that should be passed to this file
ARG VERSION=UNSET

# Derived values
# -

# Private variables, not passed in from outside, but helpful for this file
# -

# Stage 1: the builder image
FROM alpine:3.16.2 as builder

# add necessary files to the builder image: caddy config, curl, and ketcher
COPY ./Caddyfile  /builder/etc/caddy/Caddyfile
ADD https://github.com/moparisthebest/static-curl/releases/download/v7.85.0/curl-amd64 /builder/bin/curl

# we want to use the latest version of git
# hadolint ignore=DL3018
RUN apk add --no-cache git && git clone --branch chemotionified https://github.com/ptrxyz/ketcher.git /builder/www/

RUN chmod +x /builder/bin/curl && \
    mv /builder/www/ketcher.html /builder/www/index.html && \
    grep "Ketcher" /builder/www/index.html `# this is just some sanity checking...`

# Stage 2: merge with caddy + finalize the image
FROM caddy:2 as app
ARG VERSION

COPY --from=builder /builder /

HEALTHCHECK --interval=10s --timeout=3s --start-period=3s --retries=3 \
    CMD (curl --fail http://localhost:80/ | grep Ketcher) || exit 1

LABEL \
    "org.opencontainers.image.authors"="Chemotion Team" \
    "org.opencontainers.image.title"="Chemotion KetcherSVC-SC" \
    "org.opencontainers.image.description"="Image for the Chemotion KetcherSVC-Caddy sidecar container" \
    "org.opencontainers.image.version"="${VERSION}" \
    "chemotion.internal.service.id"="ketchersvc-sc"

