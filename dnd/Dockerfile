FROM ubuntu:22.04

# set timezone
ARG TZ=Europe/Berlin
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime

RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install -y --no-install-recommends --autoremove --fix-missing locales git docker.io vim make ca-certificates && \
    apt-get clean && \
    locale-gen en_US.UTF-8

ADD . /pake
