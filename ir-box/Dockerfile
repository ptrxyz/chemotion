####################################################################################################
# BASE
####################################################################################################
FROM ubuntu:21.04 AS base
ARG TINI_VERSION="v0.19.0"

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
RUN apt-get -y update && apt-get -y upgrade && \
    apt-get install locales && \
    locale-gen en_US.UTF-8

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

####################################################################################################
# SPECTRA
####################################################################################################
FROM base AS chemir

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

RUN apt-get install -y --no-install-recommends \
    gcc libxrender1 libxext-dev pkg-config \
    git ca-certificates

# COPY Anaconda3-2020.11-Linux-x86_64.sh /root/conda.sh
# COPY Anaconda2-2019.10-Linux-x86_64.sh /root/conda.sh
# ADD Miniconda3-py39_4.9.2-Linux-x86_64.sh /root/conda.sh
ADD ./miniconda.sh /root/conda.sh
ADD ./environment.yml /root/environment.yml
ADD ./src /app

RUN bash /root/conda.sh -p /anaconda3 -b && \
    echo "PATH=/anaconda3/condabin/:$PATH" >> ~/.profile && \
    rm /root/conda.sh && \
    /anaconda3/condabin/conda update -y -n base -c defaults conda && \
    /anaconda3/condabin/conda env create -f /root/environment.yml

# Make RUN commands use the new environment:
SHELL ["/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-ir", "/bin/bash", "-c"]

RUN cd /app && \
    pip install -r requirements.txt --ignore-installed

EXPOSE 3008

ENV FLASK_ENV=production

WORKDIR /app

ENTRYPOINT ["/tini", "--", "/anaconda3/condabin/conda", "run", "--no-capture-output", "-n", "chem-ir"]
CMD gunicorn -t 600 -w 1 -b 0.0.0.0:3008 server:app
